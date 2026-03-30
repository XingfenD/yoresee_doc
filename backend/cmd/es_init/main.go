package main

import (
	"context"
	"errors"

	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/search"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const reindexBatchSize = 200

var errESInitDisabled = errors.New("es_init disabled")

func main() {
	if err := bootstrap.NewInitializer().
		InitConfig().
		Check("check elasticsearch enabled", func() error {
			if config.GlobalConfig == nil || !config.GlobalConfig.Elasticsearch.Enabled {
				return errESInitDisabled
			}
			return nil
		}).
		InitPostgres().
		InitConsul().
		RequireConsulEnabled().
		InitElasticsearch().
		Err(); err != nil {
		if errors.Is(err, errESInitDisabled) {
			logrus.Println("Elasticsearch disabled, skip es_init")
			return
		}
		logrus.Fatalf("Init es_init failed: %v", err)
	}

	defer storage.ClosePostgres()
	defer storage.CloseElasticsearch()

	ctx := context.Background()
	indexName := search.DocumentIndexName()
	initialized := isESInitialized(ctx)

	exists, err := storage.ES.IndexExists(ctx, indexName)
	if err != nil {
		logrus.Fatalf("Check Elasticsearch index failed, index=%s, err=%v", indexName, err)
	}

	if initialized && exists {
		logrus.Printf("Elasticsearch already initialized, index=%s", indexName)
		return
	}

	if !exists {
		if err := storage.ES.CreateIndex(ctx, indexName, search.BuildDocumentIndexMapping()); err != nil {
			logrus.Fatalf("Create Elasticsearch index failed, index=%s, err=%v", indexName, err)
		}
		logrus.Printf("Elasticsearch index created, index=%s", indexName)
	}

	if err := reindexDocuments(ctx); err != nil {
		logrus.Fatalf("Reindex documents failed: %v", err)
	}

	if err := markESInitializedInConsul(ctx); err != nil {
		logrus.Fatalf("Mark Elasticsearch initialized failed: %v", err)
	}

	logrus.Println("Elasticsearch initialization completed successfully")
}

func isESInitialized(ctx context.Context) bool {
	value, ok, err := storage.Consul.Get(ctx, utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_ES,
		constant.ConfigKey_Third_Initialized,
	))
	if err != nil || !ok {
		return false
	}
	return value == constant.ES_Initialized_True
}

func markESInitializedInConsul(ctx context.Context) error {
	return storage.Consul.Set(
		ctx,
		utils.GenConfigKey(
			constant.ConfigKey_First_System,
			constant.ConfigKey_Second_ES,
			constant.ConfigKey_Third_Initialized,
		),
		constant.ES_Initialized_True,
	)
}

func reindexDocuments(ctx context.Context) error {
	logrus.Println("Start reindexing documents to Elasticsearch...")

	var docs []model.Document
	return storage.DB.
		Model(&model.Document{}).
		Where("deleted_at IS NULL").
		FindInBatches(&docs, reindexBatchSize, func(tx *gorm.DB, batch int) error {
			for i := range docs {
				doc := docs[i]
				if err := search.UpsertDocument(ctx, &doc); err != nil {
					return err
				}
			}
			logrus.Printf("Reindex batch done, batch=%d, size=%d", batch, len(docs))
			return nil
		}).Error
}
