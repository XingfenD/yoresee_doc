package main

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/search"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

var errSearchSyncWorkerDisabled = errors.New("search-sync-worker disabled")

func main() {
	if err := initSearchSyncWorker(); err != nil {
		if errors.Is(err, errSearchSyncWorkerDisabled) {
			logrus.Println("Elasticsearch disabled, skip search-sync-worker")
			return
		}
		logrus.Fatalf("Init search-sync-worker failed: %v", err)
	}

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

func initSearchSyncWorker() error {
	return bootstrap.NewInitializer().
		InitConfig().
		Check("check elasticsearch enabled", func() error {
			if config.GlobalConfig == nil || !config.GlobalConfig.Elasticsearch.Enabled {
				return errSearchSyncWorkerDisabled
			}
			return nil
		}).
		InitPostgres().
		InitRedis().
		InitElasticsearch().
		InitMQ().
		InitRepository().
		Err()
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
