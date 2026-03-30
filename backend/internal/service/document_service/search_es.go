package document_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/config"
	internal_dto "github.com/XingfenD/yoresee_doc/internal/service/dto"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

const defaultSearchIndexPrefix = "yoresee_doc"

func (s *DocumentService) applyElasticsearchKeywordFilter(ctx context.Context, req *internal_dto.DocumentsListReq) {
	if req == nil || req.FilterArgs == nil || req.FilterArgs.TitleKeyword == nil {
		return
	}
	keyword := strings.TrimSpace(*req.FilterArgs.TitleKeyword)
	if keyword == "" {
		return
	}
	if config.GlobalConfig == nil || !config.GlobalConfig.Elasticsearch.Enabled || storage.ES == nil {
		return
	}

	searchReq := storage.SearchDocumentsRequest{
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

	ids, err := storage.ES.SearchDocumentIDs(ctx, s.documentSearchIndexName(), searchReq)
	if err != nil {
		logrus.Warnf("[Service layer: DocumentService] elasticsearch keyword search failed, keyword=%s, err=%+v", keyword, err)
		return
	}
	req.SearchDocIDs = ids
}

func (s *DocumentService) documentSearchIndexName() string {
	prefix := defaultSearchIndexPrefix
	if config.GlobalConfig != nil {
		configPrefix := strings.TrimSpace(config.GlobalConfig.Elasticsearch.IndexPrefix)
		if configPrefix != "" {
			prefix = configPrefix
		}
	}
	return fmt.Sprintf("%s_documents", prefix)
}
