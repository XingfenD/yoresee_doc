package document_service

import (
	"context"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/sirupsen/logrus"
)

func (s *DocumentService) ListDocumentVersions(ctx context.Context, documentExternalID string, page, pageSize int32) ([]*dto.DocumentVersionResponse, int64, error) {
	if s == nil || strings.TrimSpace(documentExternalID) == "" {
		return nil, 0, status.GenErrWithCustomMsg(status.StatusParamError, "invalid list document versions request")
	}

	p := int(page)
	if p <= 0 {
		p = 1
	}
	ps := int(pageSize)
	if ps <= 0 {
		ps = 20
	}
	if ps > 100 {
		ps = 100
	}
	offset := (p - 1) * ps

	documentID, err := s.documentRepo.GetIDByExternalID(documentExternalID).Exec(ctx)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] query document id for list versions failed, document_external_id=%s, err=%+v", documentExternalID, err)
		return nil, 0, status.GenErrWithCustomMsg(status.StatusReadDBError, "query document failed")
	}
	if documentID <= 0 {
		return nil, 0, status.GenErrWithCustomMsg(status.StatusDocumentNotFound, "document not found")
	}

	items, total, err := s.docVersionRepo.ListByDocumentID(documentID, offset, ps)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] list document versions failed, document_external_id=%s, document_id=%d, err=%+v", documentExternalID, documentID, err)
		return nil, 0, status.GenErrWithCustomMsg(status.StatusReadDBError, "list document versions failed")
	}

	resp := make([]*dto.DocumentVersionResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, &dto.DocumentVersionResponse{
			Version:       item.Version,
			Title:         item.Title,
			ChangeSummary: item.ChangeSummary,
			CreatedAt:     item.CreatedAt,
		})
	}
	return resp, total, nil
}

func (s *DocumentService) GetDocumentVersionContent(ctx context.Context, documentExternalID string, version int32) (*dto.DocumentVersionResponse, error) {
	if s == nil || strings.TrimSpace(documentExternalID) == "" || version <= 0 {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "invalid get document version request")
	}

	documentID, err := s.documentRepo.GetIDByExternalID(documentExternalID).Exec(ctx)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] query document id for get version failed, document_external_id=%s, err=%+v", documentExternalID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusReadDBError, "query document failed")
	}
	if documentID <= 0 {
		return nil, status.GenErrWithCustomMsg(status.StatusDocumentNotFound, "document not found")
	}

	item, err := s.docVersionRepo.GetByDocumentIDAndVersion(documentID, int(version))
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] get document version failed, document_external_id=%s, document_id=%d, version=%d, err=%+v", documentExternalID, documentID, version, err)
		return nil, status.GenErrWithCustomMsg(status.StatusDocumentNotFound, "document version not found")
	}

	return &dto.DocumentVersionResponse{
		Version:       item.Version,
		Title:         item.Title,
		Content:       item.Content,
		ChangeSummary: item.ChangeSummary,
		CreatedAt:     item.CreatedAt,
	}, nil
}
