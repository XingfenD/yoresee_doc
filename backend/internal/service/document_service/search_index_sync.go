package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/sirupsen/logrus"
)

func (s *DocumentService) publishDocumentSearchSyncUpsertEvent(ctx context.Context, externalID string) {
	if s == nil || externalID == "" {
		return
	}
	if config.GlobalConfig == nil || !config.GlobalConfig.Elasticsearch.Enabled {
		return
	}
	if err := domain_event.PublishDocumentUpsertEvent(ctx, externalID); err != nil {
		logrus.Warnf("[Service layer: DocumentService] publish search sync event failed, external_id=%s, err=%+v", externalID, err)
	}
}
