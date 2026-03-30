package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/search"
	"github.com/XingfenD/yoresee_doc/internal/service/document_service"
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
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("Init config failed: %v", err)
	}

	if err := storage.InitPostgres(&config.GlobalConfig.Database); err != nil {
		logrus.Fatalf("Init Postgres failed: %v", err)
	}

	if err := storage.InitRedis(&config.GlobalConfig.Redis); err != nil {
		logrus.Fatalf("Init Redis failed: %v", err)
	}
	if err := storage.InitElasticsearch(&config.GlobalConfig.Elasticsearch); err != nil {
		logrus.Warnf("Init Elasticsearch failed, snapshot index sync disabled: %v", err)
	}

	if err := mq.Init(&config.GlobalConfig.MQConfig); err != nil {
		logrus.Fatalf("Init MQ failed: %v", err)
	}

	repository.MustInit()

	topic := constant.DirtyDocTopicDefault

	backend := mq.BackendRabbitMQ

	collabCoreHTTP := os.Getenv("COLLAB_CORE_HTTP")
	if collabCoreHTTP == "" {
		collabCoreHTTP = "http://collab-core:1234"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	dirtySetKey := constant.DirtyDocSetDefault
	inFlight := &sync.Map{}

	logrus.Infof("Snapshot worker started: topic=%s collabCore=%s", topic, collabCoreHTTP)

	go func() {
		if err := mq.SubscribeTo(backend, topic, func(ctx context.Context, data []byte) error {
			docID := parseDocID(data)
			if docID == "" {
				return nil
			}
			return snapshotDoc(ctx, inFlight, client, collabCoreHTTP, dirtySetKey, docID, true)
		}); err != nil {
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

	waitForShutdown()
	if err := mq.Close(); err != nil {
		logrus.Errorf("Close MQ failed: %v", err)
	}
	if err := storage.CloseRedis(); err != nil {
		logrus.Errorf("Close Redis failed: %v", err)
	}
	if err := storage.CloseElasticsearch(); err != nil {
		logrus.Errorf("Close Elasticsearch failed: %v", err)
	}
	if err := storage.ClosePostgres(); err != nil {
		logrus.Errorf("Close Postgres failed: %v", err)
	}
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
	state, err := decodeBase64(payload.State)
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
		last, parseErr := parseInt64(lastStr)
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
	syncDocumentSearchIndex(ctx, docID)

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

func syncDocumentSearchIndex(ctx context.Context, externalID string) {
	if config.GlobalConfig == nil || !config.GlobalConfig.Elasticsearch.Enabled || storage.ES == nil {
		return
	}

	doc, err := document_repo.DocumentRepo.GetByExternalID(externalID).Exec(ctx)
	if err != nil {
		logrus.Warnf("Snapshot sync search index query doc failed docId=%s err=%v", externalID, err)
		return
	}
	if err := search.UpsertDocument(ctx, doc); err != nil {
		logrus.Warnf("Snapshot sync search index failed docId=%s err=%v", externalID, err)
	}
}

func parseInt64(value string) (int64, error) {
	var result int64
	_, err := fmt.Sscanf(value, "%d", &result)
	return result, err
}

func decodeBase64(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}

func waitForShutdown() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
}
