package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/search"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func (s *DocumentService) syncDocumentSearchIndexByExternalID(ctx context.Context, externalID string) {
	if s == nil || externalID == "" {
		return
	}
	if config.GlobalConfig == nil || !config.GlobalConfig.Elasticsearch.Enabled || storage.ES == nil {
		return
	}

	docModel, err := s.documentRepo.GetByExternalID(externalID).Exec(ctx)
	if err != nil {
		logrus.Warnf("[Service layer: DocumentService] sync search index query document failed, external_id=%s, err=%+v", externalID, err)
		return
	}
	if err := search.UpsertDocument(ctx, docModel); err != nil {
		logrus.Warnf("[Service layer: DocumentService] sync search index failed, external_id=%s, err=%+v", externalID, err)
	}
}
