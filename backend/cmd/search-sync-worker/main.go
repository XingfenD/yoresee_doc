package main

import (
	"context"
	"errors"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/search"
	"github.com/XingfenD/yoresee_doc/internal/service/mq_service"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/sirupsen/logrus"
)

var errSearchSyncWorkerDisabled = errors.New("search-sync-worker disabled")

func main() {
	initializer, err := initSearchSyncWorker()
	if err != nil {
		if errors.Is(err, errSearchSyncWorkerDisabled) {
			logrus.Println("Elasticsearch disabled, skip search-sync-worker")
			return
		}
		logrus.Fatalf("Init search-sync-worker failed: %v", err)
	}

	backend := mq.BackendRabbitMQ
	topic := domain_event.DocumentSyncTopic()
	group := utils.GetEnvVar("SEARCH_SYNC_MQ_GROUP", "search-sync-worker")
	logrus.Infof("Search sync worker started: backend=%s topic=%s group=%s", backend, topic, group)

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
				return handleDocumentEvent(ctx, message.Body)
			},
		); err != nil {
			logrus.Fatalf("Subscribe search sync topic failed: %v", err)
		}
	}()

	initializer.ShutdownOnSignal(500 * time.Millisecond)
}

func initSearchSyncWorker() (*bootstrap.Initializer, error) {
	initializer := bootstrap.NewInitializer().
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
		InitRepository()
	return initializer, initializer.Err()
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
			return err
		}
		if err := search.UpsertDocument(ctx, doc); err != nil {
			logrus.Warnf("Search sync upsert failed, external_id=%s, err=%v", evt.ExternalID, err)
			return err
		}
		logrus.Infof("Search sync upsert success, external_id=%s", evt.ExternalID)
	case domain_event.DocumentActionDelete:
		// reserved for future hard-delete sync.
	default:
		logrus.Warnf("Search sync skip unknown action, action=%s, external_id=%s", evt.Action, evt.ExternalID)
	}
	return nil
}
