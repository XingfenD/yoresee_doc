package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/search"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("Init config failed: %v", err)
	}
	if config.GlobalConfig == nil || !config.GlobalConfig.Elasticsearch.Enabled {
		logrus.Println("Elasticsearch disabled, skip search-sync-worker")
		return
	}

	if err := storage.InitPostgres(&config.GlobalConfig.Database); err != nil {
		logrus.Fatalf("Init Postgres failed: %v", err)
	}
	if err := storage.InitRedis(&config.GlobalConfig.Redis); err != nil {
		logrus.Fatalf("Init Redis failed: %v", err)
	}
	if err := storage.InitElasticsearch(&config.GlobalConfig.Elasticsearch); err != nil {
		logrus.Fatalf("Init Elasticsearch failed: %v", err)
	}
	if err := mq.Init(&config.GlobalConfig.MQConfig); err != nil {
		logrus.Fatalf("Init MQ failed: %v", err)
	}
	repository.MustInit()

	backend := resolveMQBackend()
	topic := domain_event.DocumentSyncTopic()
	logrus.Infof("Search sync worker started: backend=%s topic=%s", backend, topic)

	go func() {
		if err := mq.SubscribeTo(backend, topic, handleDocumentEvent); err != nil {
			logrus.Fatalf("Subscribe search sync topic failed: %v", err)
		}
	}()

	utils.WaitForShutdownSignal()
	time.Sleep(500 * time.Millisecond)
	if err := mq.Close(); err != nil {
		logrus.Errorf("Close MQ failed: %v", err)
	}
	if err := storage.CloseElasticsearch(); err != nil {
		logrus.Errorf("Close Elasticsearch failed: %v", err)
	}
	if err := storage.CloseRedis(); err != nil {
		logrus.Errorf("Close Redis failed: %v", err)
	}
	if err := storage.ClosePostgres(); err != nil {
		logrus.Errorf("Close Postgres failed: %v", err)
	}
}

func resolveMQBackend() mq.Backend {
	raw := strings.TrimSpace(os.Getenv("SEARCH_SYNC_MQ"))
	switch strings.ToLower(raw) {
	case string(mq.BackendRedis):
		return mq.BackendRedis
	case string(mq.BackendRabbitMQ), "":
		return mq.BackendRabbitMQ
	default:
		return mq.BackendRabbitMQ
	}
}

func handleDocumentEvent(ctx context.Context, data []byte) error {
	evt, err := domain_event.DecodeDocumentSyncEvent(data)
	if err != nil {
		logrus.Warnf("Parse search sync event failed: %v", err)
		return nil
	}

	switch evt.Action {
	case domain_event.DocumentActionUpsert:
		doc, err := document_repo.DocumentRepo.GetByExternalID(evt.ExternalID).Exec(ctx)
		if err != nil {
			logrus.Warnf("Search sync load document failed, external_id=%s, err=%v", evt.ExternalID, err)
			return nil
		}
		if err := search.UpsertDocument(ctx, doc); err != nil {
			logrus.Warnf("Search sync upsert failed, external_id=%s, err=%v", evt.ExternalID, err)
			return nil
		}
		logrus.Infof("Search sync upsert success, external_id=%s", evt.ExternalID)
	case domain_event.DocumentActionDelete:
		// reserved for future hard-delete sync.
	default:
		logrus.Warnf("Search sync skip unknown action, action=%s, external_id=%s", evt.Action, evt.ExternalID)
	}
	return nil
}
