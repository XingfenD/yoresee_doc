package document_service

import (
	"context"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/search"
	internal_dto "github.com/XingfenD/yoresee_doc/internal/service/dto"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func (s *DocumentService) applyElasticsearchKeywordFilter(ctx context.Context, req *internal_dto.DocumentsListReq) {
	if req == nil || req.FilterArgs == nil || req.FilterArgs.TitleKeyword == nil {
		return
	}
	keyword := strings.TrimSpace(*req.FilterArgs.TitleKeyword)
	if keyword == "" {
		return
	}
	if storage.ES == nil {
		return
	}

	searchReq := search.DocumentSearchRequest{
		Keyword:              keyword,
		DocType:              req.FilterArgs.DocType,
		Status:               req.FilterArgs.Status,
		Tags:                 req.FilterArgs.Tags,
		CreateTimeRangeStart: req.FilterArgs.CreateTimeRangeStart,
		CreateTimeRangeEnd:   req.FilterArgs.CreateTimeRangeEnd,
		UpdateTimeRangeStart: req.FilterArgs.UpdateTimeRangeStart,
		UpdateTimeRangeEnd:   req.FilterArgs.UpdateTimeRangeEnd,
		Size:                 5000,
		ListOwnDoc:           req.ListOwnDoc,
	}
	if req.MetaArgs != nil {
		searchReq.UserID = req.MetaArgs.UserID
		searchReq.KnowledgeID = req.MetaArgs.KnowledgeID
	}

	ids, err := search.SearchDocumentIDs(ctx, searchReq)
	if err != nil {
		logrus.Warnf("[Service layer: DocumentService] elasticsearch keyword search failed, keyword=%s, err=%+v", keyword, err)
		return
	}
	req.SearchDocIDs = ids
}
