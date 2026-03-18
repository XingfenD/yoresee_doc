package grpcserver

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service/document_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type DocumentServiceServer struct {
	pb.UnimplementedDocumentServiceServer
}

func NewDocumentServiceServer() *DocumentServiceServer {
	return &DocumentServiceServer{}
}

func (s *DocumentServiceServer) ListDocuments(ctx context.Context, req *pb.ListDocumentsRequest) (*pb.ListDocumentsResponse, error) {
	if req == nil {
		return &pb.ListDocumentsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	filterArgs := &dto.DocumentsListFilterArgs{
		TitleKeyword: req.TitleKeyword,
		DocType:      req.Type,
		Tags:         req.Tags,
	}
	if req.Status != nil {
		filterArgs.Status = utils.Of(int(req.GetStatus()))
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

	serviceReq := &dto.ListDocumentsByExternalReq{
		ExternalArgs: &dto.DocumentsListExternalArgs{
			UserExternalID:         req.UserExternalId,
			RootDocumentExternalID: req.RootDocumentExternalId,
		},
		FilterArgs: filterArgs,
		SortArgs:   sortArgs,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}
	if req.Options != nil {
		serviceReq.Options = &dto.RecursiveOptions{
			IncludeChildren: req.Options.IncludeChildren,
			Recursive:       req.Options.Recursive,
			Depth:           utils.Of(int(req.Options.Depth)),
		}
	}

	docs, _, err := document_service.DocumentSvc.ListDocumentsByExternal(ctx, serviceReq)
	if err != nil {
		return &pb.ListDocumentsResponse{Base: baseResponseFromErr(err)}, nil
	}

	respDocs := make([]*pb.DocumentResponse, 0, len(docs))
	for _, doc := range docs {
		respDocs = append(respDocs, toDocumentResponse(doc))
	}

	return &pb.ListDocumentsResponse{
		Base:      baseResponseFromErr(nil),
		Documents: respDocs,
	}, nil
}

func (s *DocumentServiceServer) GetDocumentContent(ctx context.Context, req *pb.GetDocumentContentRequest) (*pb.GetDocumentContentResponse, error) {
	if req == nil || req.DocumentExternalId == "" {
		return &pb.GetDocumentContentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	document, err := document_service.DocumentSvc.GetDocumentByExternalID(ctx, req.DocumentExternalId)
	if err != nil {
		return &pb.GetDocumentContentResponse{Base: baseResponseFromErr(status.StatusDocumentNotFound)}, nil
	}

	return &pb.GetDocumentContentResponse{
		Base:     baseResponseFromErr(nil),
		Document: toDocumentResponse(&document.DocumentMetaResponse),
		Content:  document.Content,
	}, nil
}

func (s *DocumentServiceServer) GetDocumentYjsSnapshot(ctx context.Context, req *pb.GetDocumentYjsSnapshotRequest) (*pb.GetDocumentYjsSnapshotResponse, error) {
	if req == nil || req.DocumentExternalId == "" {
		return &pb.GetDocumentYjsSnapshotResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	state, err := document_service.DocumentSvc.GetDocumentYjsSnapshot(ctx, req.DocumentExternalId)
	if err != nil {
		return &pb.GetDocumentYjsSnapshotResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.GetDocumentYjsSnapshotResponse{
		Base:  baseResponseFromErr(nil),
		State: state,
	}, nil
}

func (s *DocumentServiceServer) SaveDocumentYjsSnapshot(ctx context.Context, req *pb.SaveDocumentYjsSnapshotRequest) (*pb.SaveDocumentYjsSnapshotResponse, error) {
	if req == nil || req.DocumentExternalId == "" || len(req.State) == 0 {
		return &pb.SaveDocumentYjsSnapshotResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	if err := document_service.DocumentSvc.SaveDocumentYjsSnapshot(ctx, req.DocumentExternalId, req.State); err != nil {
		return &pb.SaveDocumentYjsSnapshotResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.SaveDocumentYjsSnapshotResponse{
		Base: baseResponseFromErr(nil),
	}, nil
}

func (s *DocumentServiceServer) GetOwnDocuments(ctx context.Context, req *pb.GetOwnDocumentsRequest) (*pb.GetOwnDocumentsResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.GetOwnDocumentsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	queryReq := &dto.ListDocumentsByExternalReq{
		ExternalArgs: &dto.DocumentsListExternalArgs{
			UserExternalID: &userExternalID,
		},
		ListDocumentsBaseArgs: dto.ListDocumentsBaseArgs{
			ListOwnDoc:    true,
			DirectoryOnly: req.DirectoryOnly,
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

	docs, count, err := document_service.DocumentSvc.ListDocumentsByExternal(ctx, queryReq)
	if err != nil {
		return &pb.GetOwnDocumentsResponse{Base: baseResponseFromErr(err)}, nil
	}

	respDocs := make([]*pb.DocumentResponse, 0, len(docs))
	for _, doc := range docs {
		respDocs = append(respDocs, toDocumentResponse(doc))
	}

	return &pb.GetOwnDocumentsResponse{
		Base:       baseResponseFromErr(nil),
		Documents:  respDocs,
		TotalCount: count,
	}, nil
}

func (s *DocumentServiceServer) CreateDocument(ctx context.Context, req *pb.CreateDocumentRequest) (*pb.CreateDocumentResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.CreateDocumentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if req == nil || req.Title == "" {
		return &pb.CreateDocumentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	dtoReq := &dto.CreateDocumentReq{
		Title:             req.Title,
		Type:              fromDocumentType(req.Type),
		CreatorExternalID: utils.Of(userExternalID),
		ParentExternalID:  req.ParentExternalId,
	}

	switch req.ContainerType {
	case pb.CreateDocumentContainerType_CREATE_DOCUMENT_CONTAINER_TYPE_KNOWLEDGE_BASE:
		dtoReq.CreateAsOwnDoc = false
		dtoReq.KnowledgeExternalID = req.KnowledgeBaseExternalId
	case pb.CreateDocumentContainerType_CREATE_DOCUMENT_CONTAINER_TYPE_OWN:
		dtoReq.CreateAsOwnDoc = true
	default:
		return &pb.CreateDocumentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	resp, err := document_service.DocumentSvc.Create(ctx, dtoReq)
	if err != nil {
		return &pb.CreateDocumentResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.CreateDocumentResponse{
		Base:       baseResponseFromErr(nil),
		ExternalId: resp.ExternalID,
	}, nil
}

func validateUpdateDocumentRequest(req *pb.UpdateDocumentRequest) error {
	if req == nil {
		return status.StatusParamError
	}
	if req.Content == nil && req.KnowledgeBaseExternalId == nil && req.ParentExternalId == nil && req.Title == nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "no update fields")
	}
	if req.MoveToContainer == nil || *req.MoveToContainer == pb.CreateDocumentContainerType_CREATE_DOCUMENT_CONTAINER_TYPE_OWN {
		if req.KnowledgeBaseExternalId != nil {
			return status.GenErrWithCustomMsg(status.StatusParamError, "KnowledgeBaseExternalId not nil")
		}
	}
	if *req.MoveToContainer == pb.CreateDocumentContainerType_CREATE_DOCUMENT_CONTAINER_TYPE_KNOWLEDGE_BASE &&
		req.KnowledgeBaseExternalId == nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "KnowledgeBaseExternalId is nil")
	}
	return nil
}

func (s *DocumentServiceServer) UpdateDocument(ctx context.Context, req *pb.UpdateDocumentRequest) (*pb.UpdateDocumentResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.UpdateDocumentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := validateUpdateDocumentRequest(req); err != nil {
		return &pb.UpdateDocumentResponse{Base: baseResponseFromErr(err)}, nil
	}

	serviceReq := &dto.UpdateDocumentRequest{
		ExternalID:              req.ExternalId,
		Title:                   req.Title,
		ParentExternalID:        req.ParentExternalId,
		KnowledgeBaseExternalID: req.KnowledgeBaseExternalId,
	}

	if req.GetMoveToContainer() == pb.CreateDocumentContainerType_CREATE_DOCUMENT_CONTAINER_TYPE_OWN {
		serviceReq.MoveAsOwn = true
	}

	return &pb.UpdateDocumentResponse{
		Base: baseResponseFromErr(nil),
	}, nil
}

func (s *DocumentServiceServer) UpdateDocumentMeta(ctx context.Context, req *pb.UpdateDocumentMetaRequest) (*pb.UpdateDocumentMetaResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.UpdateDocumentMetaResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if req == nil || req.ExternalId == "" {
		return &pb.UpdateDocumentMetaResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if req.Title == nil && req.Summary == nil && req.Status == nil && req.Tags == nil {
		return &pb.UpdateDocumentMetaResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	serviceReq := &dto.UpdateDocumentMetaRequest{
		ExternalID: req.ExternalId,
		Title:      req.Title,
		Summary:    req.Summary,
	}
	if req.Status != nil {
		statusVal := int(req.GetStatus())
		serviceReq.Status = &statusVal
	}
	if req.Tags != nil {
		tags := req.Tags
		serviceReq.Tags = &tags
	}

	if _, err := document_service.DocumentSvc.UpdateDocumentMeta(ctx, serviceReq); err != nil {
		return &pb.UpdateDocumentMetaResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.UpdateDocumentMetaResponse{
		Base: baseResponseFromErr(nil),
	}, nil
}
