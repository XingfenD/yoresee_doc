package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/service/document_service"
	"github.com/XingfenD/yoresee_doc/internal/service/mq_service"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/constant"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/bytedance/sonic"
	"github.com/sirupsen/logrus"
)

type dirtyDocMessage struct {
	DocID string `json:"doc_id"`
	DocId string `json:"docId"`
}

func main() {
	initializer := bootstrap.NewInitializer().
		InitConfig().
		InitPostgres().
		InitRedis().
		InitElasticsearchAllowFail().
		InitMQ().
		InitRepository()
	if err := initializer.Err(); err != nil {
		logrus.Fatalf("Init snapshot-worker failed: %v", err)
	}

	topic := constant.DirtyDocTopicDefault

	backend := mq.BackendRabbitMQ
	group := resolveMQConsumerGroup()

	collabCoreHTTP := os.Getenv("COLLAB_CORE_HTTP")
	if collabCoreHTTP == "" {
		collabCoreHTTP = "http://collab-core:1234"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	dirtySetKey := constant.DirtyDocSetDefault
	inFlight := &sync.Map{}

	logrus.Infof("Snapshot worker started: topic=%s group=%s collabCore=%s", topic, group, collabCoreHTTP)

	go func() {
		if err := mq_service.MQSvc.Consume(
			context.Background(),
			backend,
			mq.ConsumeOptions{
				Topic:   topic,
				Mode:    mq.ConsumeModeGroup,
				Group:   group,
				AutoAck: false,
				OnError: mq.ErrorActionRequeue,
			},
			func(ctx context.Context, message mq.Message) error {
				docID := parseDocID(message.Body)
				if docID == "" {
					return nil
				}
				return snapshotDoc(ctx, inFlight, client, collabCoreHTTP, dirtySetKey, docID, true)
			},
		); err != nil {
			logrus.Fatalf("Subscribe failed: %v", err)
		}
	}()

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			docIDs, err := scanDirtyDocs(client, collabCoreHTTP, dirtySetKey)
			if err != nil {
				logrus.Errorf("Scan dirty docs failed: %v", err)
				continue
			}
			logrus.Infof("Scan dirty docs: candidates=%d", len(docIDs))
			if len(docIDs) > 0 {
				logrus.Infof("Dirty doc candidates: %v", docIDs)
			}
			for _, docID := range docIDs {
				ctx := context.Background()
				_ = snapshotDoc(ctx, inFlight, client, collabCoreHTTP, dirtySetKey, docID, false)
			}
		}
	}()

	initializer.ShutdownOnSignal(0)
}

func resolveMQConsumerGroup() string {
	group := strings.TrimSpace(os.Getenv("SNAPSHOT_MQ_GROUP"))
	if group == "" {
		return "snapshot-worker"
	}
	return group
}

func parseDocID(data []byte) string {
	payload := strings.TrimSpace(string(data))
	if payload == "" {
		return ""
	}
	var msg dirtyDocMessage
	if err := sonic.Unmarshal(data, &msg); err == nil {
		if msg.DocID != "" {
			return msg.DocID
		}
		if msg.DocId != "" {
			return msg.DocId
		}
	}
	return payload
}

func fetchDocSnapshot(client *http.Client, baseURL, docID string) ([]byte, string, error) {
	url := fmt.Sprintf("%s/internal/yjs/doc-snapshot/%s", strings.TrimRight(baseURL, "/"), docID)
	resp, err := client.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	var payload struct {
		State   string `json:"state"`
		Content string `json:"content"`
	}
	if err := sonic.Unmarshal(body, &payload); err != nil {
		return nil, "", err
	}
	if payload.State == "" {
		return nil, payload.Content, nil
	}
	state, err := utils.DecodeBase64(payload.State)
	if err != nil {
		return nil, "", err
	}
	return state, payload.Content, nil
}

func scanDirtyDocs(client *http.Client, baseURL, dirtySetKey string) ([]string, error) {
	ctx := context.Background()
	if storage.GetRedis() == nil {
		return nil, nil
	}
	docIDs, err := storage.GetRedis().SMembers(ctx, dirtySetKey).Result()
	if err != nil {
		return nil, err
	}
	if len(docIDs) == 0 {
		return nil, nil
	}
	now := time.Now().UnixMilli()
	candidates := make([]string, 0, len(docIDs))
	for _, docID := range docIDs {
		if docID == "" {
			continue
		}
		roomKey := fmt.Sprintf("collab:room:doc-%s", docID)
		lastStr, err := storage.GetRedis().Get(ctx, roomKey).Result()
		if err != nil {
			logrus.Infof("Dirty doc %s skipped: room key missing", docID)
			continue
		}
		last, parseErr := utils.ParseInt64(lastStr)
		if parseErr != nil {
			logrus.Infof("Dirty doc %s skipped: invalid room timestamp %s", docID, lastStr)
			continue
		}
		if now-last < 10_000 {
			logrus.Infof("Dirty doc %s skipped: lastEditAgo=%dms", docID, now-last)
			continue
		}
		candidates = append(candidates, docID)
	}
	return candidates, nil
}

func snapshotDoc(ctx context.Context, inFlight *sync.Map, client *http.Client, baseURL, dirtySetKey, docID string, force bool) error {
	if _, loaded := inFlight.LoadOrStore(docID, struct{}{}); loaded {
		return nil
	}
	defer inFlight.Delete(docID)

	logrus.Infof("Snapshot start docId=%s force=%v", docID, force)
	state, content, err := fetchDocSnapshot(client, baseURL, docID)
	if err != nil {
		logrus.Errorf("Snapshot fetch failed docId=%s err=%v", docID, err)
		return err
	}
	if len(state) == 0 {
		logrus.Infof("Snapshot empty docId=%s", docID)
		return nil
	}
	logrus.Infof("Snapshot fetched docId=%s bytes=%d", docID, len(state))

	if err := document_service.DocumentSvc.SaveDocumentYjsSnapshot(ctx, docID, state); err != nil {
		logrus.Errorf("Snapshot save failed docId=%s err=%v", docID, err)
		return err
	}

	if err := document_repo.DocumentRepo.UpdateContentByExternalID(docID, content).Exec(ctx); err != nil {
		logrus.Errorf("Snapshot content update failed docId=%s err=%v", docID, err)
		return err
	}
	publishDocumentSearchSyncUpsertEvent(ctx, docID)

	if storage.GetRedis() != nil {
		key := fmt.Sprintf("collab:yjs:doc:updates:%s", docID)
		pipe := storage.GetRedis().TxPipeline()
		pipe.Del(ctx, key)
		pipe.RPush(ctx, key, state)
		if _, err := pipe.Exec(ctx); err != nil {
			logrus.Errorf("Snapshot redis list replace failed docId=%s err=%v", docID, err)
			return err
		}
		if err := storage.GetRedis().SRem(ctx, dirtySetKey, docID).Err(); err != nil {
			logrus.Errorf("Snapshot redis srem failed docId=%s err=%v", docID, err)
			return err
		}
	}

	if force {
		logrus.Infof("Snapshot saved (mq) for %s", docID)
	} else {
		logrus.Infof("Snapshot saved (scan) for %s", docID)
	}
	return nil
}

func publishDocumentSearchSyncUpsertEvent(ctx context.Context, externalID string) {
	if err := domain_event.PublishDocumentUpsertEvent(ctx, externalID); err != nil {
		logrus.Warnf("Snapshot publish search sync event failed docId=%s err=%v", externalID, err)
	}
}
