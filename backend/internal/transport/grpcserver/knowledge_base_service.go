package grpcserver

import (
	"context"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service/document_service"
	"github.com/XingfenD/yoresee_doc/internal/service/knowledge_base_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type KnowledgeBaseServiceServer struct {
	pb.UnimplementedKnowledgeBaseServiceServer
}

func NewKnowledgeBaseServiceServer() *KnowledgeBaseServiceServer {
	return &KnowledgeBaseServiceServer{}
}

func (s *KnowledgeBaseServiceServer) ListKnowledgeBases(ctx context.Context, req *pb.ListKnowledgeBasesRequest) (*pb.ListKnowledgeBasesResponse, error) {
	if req == nil {
		return &pb.ListKnowledgeBasesResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.ListKnowledgeBasesResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	filterArgs := &dto.KnowledgeBaseListFilterArgs{
		IsPublic:    req.IsPublic,
		NameKeyword: req.NameKeyword,
	}
	if req.CreateTimeRange != nil {
		filterArgs.CreateTimeRangeStart = req.CreateTimeRange.Start
		filterArgs.CreateTimeRangeEnd = req.CreateTimeRange.End
	}
	if req.UpdateTimeRange != nil {
		filterArgs.UpdateTimeRangeStart = req.UpdateTimeRange.Start
		filterArgs.UpdateTimeRangeEnd = req.UpdateTimeRange.End
	}

	sortArgs := dto.SortArgs{Field: "created_at", Desc: true}
	if req.OrderBy != nil {
		sortArgs.Field = req.GetOrderBy()
	}
	if req.OrderDesc != nil {
		sortArgs.Desc = req.GetOrderDesc()
	}

	serviceReq := &dto.KnowledgeBaseListByExternalReq{
		CreatorExternalID: "",
		FilterArgs:        filterArgs,
		SortArgs:          sortArgs,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}
	if req.OnlyMine {
		serviceReq.CreatorExternalID = userExternalID
	}

	kbs, total, err := knowledge_base_service.KnowledgeBaseSvc.ListByExternal(serviceReq)
	if err != nil {
		return &pb.ListKnowledgeBasesResponse{Base: baseResponseFromErr(err)}, nil
	}

	respKBs := make([]*pb.KnowledgeBaseResponse, 0, len(kbs))
	for _, kb := range kbs {
		respKBs = append(respKBs, toKnowledgeBaseResponse(kb))
	}

	return &pb.ListKnowledgeBasesResponse{
		Base:           baseResponseFromErr(nil),
		KnowledgeBases: respKBs,
		Total:          total,
	}, nil
}

func (s *KnowledgeBaseServiceServer) GetKnowledgeBase(ctx context.Context, req *pb.GetKnowledgeBaseRequest) (*pb.GetKnowledgeBaseResponse, error) {
	if req == nil || req.KnowledgeBaseExternalId == "" {
		return &pb.GetKnowledgeBaseResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.GetKnowledgeBaseResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	kbDTO, err := knowledge_base_service.KnowledgeBaseSvc.GetByExternalID(&dto.KnowledgeBaseGetByExternalIDReq{
		KnowledgeBaseExternalID: req.KnowledgeBaseExternalId,
	}).WithExtend().Exec()
	if err != nil {
		return &pb.GetKnowledgeBaseResponse{Base: baseResponseFromErr(status.StatusKnowledgeBaseNotFound)}, nil
	}

	if req.RecordRecentLog {
		knowledge_base_service.KnowledgeBaseSvc.CreateRecentKnowledgeBase(&dto.CreateRecentKnowledgeBaseRequest{
			UserExternalID:          userExternalID,
			KnowledgeBaseExternalID: req.KnowledgeBaseExternalId,
			AssessTime:              time.Now(),
		})
	}

	svcReq := &dto.ListDocumentsByExternalReq{
		ExternalArgs: &dto.DocumentsListExternalArgs{
			KnowledgeExternalID: utils.Of(kbDTO.ExternalID),
		},
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
		Options: &dto.RecursiveOptions{
			IncludeChildren: true,
			Recursive:       true,
		},
	}
	documents, totalCount, err := document_service.DocumentSvc.ListDocumentsByExternal(
		ctx,
		svcReq,
	)
	if err != nil {
		return &pb.GetKnowledgeBaseResponse{Base: baseResponseFromErr(status.StatusDocumentNotFound)}, nil
	}

	respDocs := make([]*pb.DocumentResponse, 0, len(documents))
	for _, doc := range documents {
		respDocs = append(respDocs, toDocumentResponse(doc))
	}

	return &pb.GetKnowledgeBaseResponse{
		Base:          baseResponseFromErr(nil),
		KnowledgeBase: toKnowledgeBaseResponse(kbDTO),
		Documents:     respDocs,
		TotalCount:    totalCount,
	}, nil
}
