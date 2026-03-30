package document_service

import (
	"context"
	"strconv"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/model"
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
	if err := storage.ES.UpsertDocument(ctx, s.documentSearchIndexName(), strconv.FormatInt(docModel.ID, 10), buildSearchIndexDocument(docModel)); err != nil {
		logrus.Warnf("[Service layer: DocumentService] sync search index failed, external_id=%s, err=%+v", externalID, err)
	}
}

func buildSearchIndexDocument(doc *model.Document) map[string]interface{} {
	payload := map[string]interface{}{
		"id":          doc.ID,
		"external_id": doc.ExternalID,
		"title":       doc.Title,
		"summary":     doc.Summary,
		"content":     doc.Content,
		"type":        doc.Type,
		"user_id":     doc.UserID,
		"status":      doc.Status,
		"tags":        doc.Tags,
		"created_at":  doc.CreatedAt.Format(time.RFC3339),
		"updated_at":  doc.UpdatedAt.Format(time.RFC3339),
	}
	if doc.KnowledgeID != nil {
		payload["knowledge_id"] = *doc.KnowledgeID
	}
	return payload
}
