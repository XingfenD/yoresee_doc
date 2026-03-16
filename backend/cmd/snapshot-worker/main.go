package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/bytedance/sonic"
	"github.com/sirupsen/logrus"
)

const (
	defaultDirtyDocTopic = "collab.dirty_docs"
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

	if err := mq.Init(&config.GlobalConfig.MQConfig); err != nil {
		logrus.Fatalf("Init MQ failed: %v", err)
	}

	repository.MustInit()

	topic := os.Getenv("DIRTY_DOC_TOPIC")
	if topic == "" {
		topic = defaultDirtyDocTopic
	}

	mqBackend := strings.ToLower(os.Getenv("DIRTY_DOC_MQ"))
	var backend mq.Backend
	switch mqBackend {
	case "rabbit", "rabbitmq":
		backend = mq.BackendRabbitMQ
	default:
		backend = mq.BackendRedis
	}

	collabCoreHTTP := os.Getenv("COLLAB_CORE_HTTP")
	if collabCoreHTTP == "" {
		collabCoreHTTP = "http://collab-core:1234"
	}

	client := &http.Client{Timeout: 10 * time.Second}

	logrus.Infof("Snapshot worker started: topic=%s collabCore=%s", topic, collabCoreHTTP)

	go func() {
		if err := mq.SubscribeTo(backend, topic, func(ctx context.Context, data []byte) error {
			docID := parseDocID(data)
			if docID == "" {
				return nil
			}
			state, err := fetchDocState(client, collabCoreHTTP, docID)
			if err != nil {
				return err
			}
			if len(state) == 0 {
				return nil
			}

			if err := service.DocumentSvc.SaveDocumentYjsSnapshot(ctx, docID, state); err != nil {
				return err
			}

			if storage.GetRedis() != nil {
				key := fmt.Sprintf("yjs:doc:%s", docID)
				if err := storage.GetRedis().Set(ctx, key, state, 0).Err(); err != nil {
					return err
				}
			}

			logrus.Infof("Snapshot saved for %s", docID)
			return nil
		}); err != nil {
			logrus.Fatalf("Subscribe failed: %v", err)
		}
	}()

	waitForShutdown()
	if err := mq.Close(); err != nil {
		logrus.Errorf("Close MQ failed: %v", err)
	}
	if err := storage.CloseRedis(); err != nil {
		logrus.Errorf("Close Redis failed: %v", err)
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

func fetchDocState(client *http.Client, baseURL, docID string) ([]byte, error) {
	url := fmt.Sprintf("%s/internal/yjs/doc/%s", strings.TrimRight(baseURL, "/"), docID)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func waitForShutdown() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
}
