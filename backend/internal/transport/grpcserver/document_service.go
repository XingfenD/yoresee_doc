package grpcserver

import (
	"context"
	"time"

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
			KnowledgeExternalID:    req.KnowledgeBaseExternalId,
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

	docs, total, err := document_service.DocumentSvc.ListDocumentsByExternal(ctx, serviceReq)
	if err != nil {
		return &pb.ListDocumentsResponse{Base: baseResponseFromErr(err)}, nil
	}

	respDocs := make([]*pb.DocumentResponse, 0, len(docs))
	for _, doc := range docs {
		respDocs = append(respDocs, toDocumentResponse(doc))
	}

	return &pb.ListDocumentsResponse{
		Base:       baseResponseFromErr(nil),
		Documents:  respDocs,
		TotalCount: total,
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

func (s *DocumentServiceServer) GetDocumentSettings(ctx context.Context, req *pb.GetDocumentSettingsRequest) (*pb.GetDocumentSettingsResponse, error) {
	if req == nil || req.DocumentExternalId == "" {
		return &pb.GetDocumentSettingsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	settings, err := document_service.DocumentSvc.GetDocumentSettings(ctx, &dto.GetDocumentSettingsRequest{
		ExternalID: req.DocumentExternalId,
	})
	if err != nil {
		return &pb.GetDocumentSettingsResponse{Base: baseResponseFromErr(status.StatusDocumentNotFound)}, nil
	}

	return &pb.GetDocumentSettingsResponse{
		Base:     baseResponseFromErr(nil),
		IsPublic: utils.Of(settings.IsPublic),
	}, nil
}

func (s *DocumentServiceServer) RecordRecentDocument(ctx context.Context, req *pb.RecordRecentDocumentRequest) (*pb.RecordRecentDocumentResponse, error) {
	if req == nil || req.DocumentExternalId == "" {
		return &pb.RecordRecentDocumentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.RecordRecentDocumentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := document_service.DocumentSvc.RecordRecentDocument(&dto.RecordRecentDocumentRequest{
		UserExternalID:     userExternalID,
		DocumentExternalID: req.DocumentExternalId,
	}); err != nil {
		return &pb.RecordRecentDocumentResponse{Base: baseResponseFromErr(err)}, nil
	}
	return &pb.RecordRecentDocumentResponse{Base: baseResponseFromErr(nil)}, nil
}

func (s *DocumentServiceServer) ListRecentDocuments(ctx context.Context, req *pb.ListRecentDocumentsRequest) (*pb.ListRecentDocumentsResponse, error) {
	if req == nil {
		return &pb.ListRecentDocumentsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.ListRecentDocumentsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	var startTime *time.Time
	if req.StartTime != nil && *req.StartTime != "" {
		t, err := time.Parse(time.RFC3339, *req.StartTime)
		if err != nil {
			return &pb.ListRecentDocumentsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
		}
		startTime = &t
	}
	var endTime *time.Time
	if req.EndTime != nil && *req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, *req.EndTime)
		if err != nil {
			return &pb.ListRecentDocumentsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
		}
		endTime = &t
	}

	serviceReq := &dto.ListRecentDocumentsRequest{
		UserExternalID: userExternalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}
	documents, total, err := document_service.DocumentSvc.ListRecentDocuments(serviceReq)
	if err != nil {
		return &pb.ListRecentDocumentsResponse{Base: baseResponseFromErr(err)}, nil
	}

	respDocs := make([]*pb.DocumentResponse, 0, len(documents))
	for _, doc := range documents {
		respDocs = append(respDocs, toDocumentResponse(&doc.DocumentMetaResponse))
	}

	return &pb.ListRecentDocumentsResponse{
		Base:      baseResponseFromErr(nil),
		Documents: respDocs,
		Total:     total,
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
		IsPublic:          req.GetIsPublic(),
		CreatorExternalID: utils.Of(userExternalID),
		ParentExternalID:  req.ParentExternalId,
	}
	if req.TemplateId != nil && *req.TemplateId > 0 {
		dtoReq.TemplateID = req.TemplateId
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
	if req.Title == nil && req.Summary == nil && req.Tags == nil {
		return &pb.UpdateDocumentMetaResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	serviceReq := &dto.UpdateDocumentMetaRequest{
		ExternalID: req.ExternalId,
		Title:      req.Title,
		Summary:    req.Summary,
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

func (s *DocumentServiceServer) UpdateDocumentSettings(ctx context.Context, req *pb.UpdateDocumentSettingsRequest) (*pb.UpdateDocumentSettingsResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.UpdateDocumentSettingsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if req == nil || req.ExternalId == "" || req.IsPublic == nil {
		return &pb.UpdateDocumentSettingsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	serviceReq := &dto.UpdateDocumentSettingsRequest{
		ExternalID: req.ExternalId,
		IsPublic:   req.GetIsPublic(),
	}
	updated, err := document_service.DocumentSvc.UpdateDocumentSettings(ctx, serviceReq)
	if err != nil {
		return &pb.UpdateDocumentSettingsResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.UpdateDocumentSettingsResponse{
		Base:     baseResponseFromErr(nil),
		IsPublic: utils.Of(updated.IsPublic),
	}, nil
}

func (s *DocumentServiceServer) UploadDocumentAttachment(ctx context.Context, req *pb.UploadDocumentAttachmentRequest) (*pb.UploadDocumentAttachmentResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.UploadDocumentAttachmentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if req == nil || req.DocumentExternalId == "" || len(req.FileContent) == 0 {
		return &pb.UploadDocumentAttachmentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	attachment, err := document_service.DocumentSvc.UploadAttachment(
		ctx,
		userExternalID,
		req.DocumentExternalId,
		req.FileContent,
		req.FileName,
		req.GetContentType(),
	)
	if err != nil {
		return &pb.UploadDocumentAttachmentResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.UploadDocumentAttachmentResponse{
		Base:       baseResponseFromErr(nil),
		Attachment: toAttachmentResponse(attachment),
	}, nil
}

func (s *DocumentServiceServer) ListDocumentAttachments(ctx context.Context, req *pb.ListDocumentAttachmentsRequest) (*pb.ListDocumentAttachmentsResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.ListDocumentAttachmentsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if req == nil || req.DocumentExternalId == "" {
		return &pb.ListDocumentAttachmentsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	attachments, err := document_service.DocumentSvc.ListAttachments(ctx, req.DocumentExternalId)
	if err != nil {
		return &pb.ListDocumentAttachmentsResponse{Base: baseResponseFromErr(err)}, nil
	}
	respAttachments := make([]*pb.AttachmentResponse, 0, len(attachments))
	for _, attachment := range attachments {
		respAttachments = append(respAttachments, toAttachmentResponse(attachment))
	}

	return &pb.ListDocumentAttachmentsResponse{
		Base:        baseResponseFromErr(nil),
		Attachments: respAttachments,
	}, nil
}

func (s *DocumentServiceServer) DeleteDocumentAttachment(ctx context.Context, req *pb.DeleteDocumentAttachmentRequest) (*pb.DeleteDocumentAttachmentResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.DeleteDocumentAttachmentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if req == nil || req.DocumentExternalId == "" || req.AttachmentExternalId == "" {
		return &pb.DeleteDocumentAttachmentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	if err := document_service.DocumentSvc.DeleteAttachment(ctx, req.DocumentExternalId, req.AttachmentExternalId); err != nil {
		return &pb.DeleteDocumentAttachmentResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.DeleteDocumentAttachmentResponse{
		Base: baseResponseFromErr(nil),
	}, nil
}

func (s *DocumentServiceServer) CreateTemplate(ctx context.Context, req *pb.CreateTemplateRequest) (*pb.CreateTemplateResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.CreateTemplateResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}
	if req == nil {
		return &pb.CreateTemplateResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	serviceReq := &dto.CreateTemplateRequest{
		UserExternalID:  userExternalID,
		TargetContainer: fromCreateTemplateContainer(req.TargetContainer),
		TemplateContent: req.TemplateContent,
	}
	if req.KnowledgeBaseId != "" {
		serviceReq.KnowledgeBaseExternalID = utils.Of(req.KnowledgeBaseId)
	}

	if err := document_service.DocumentSvc.CreateTemplate(ctx, serviceReq); err != nil {
		return &pb.CreateTemplateResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.CreateTemplateResponse{
		Base: baseResponseFromErr(nil),
	}, nil
}

func (s *DocumentServiceServer) GetTemplate(ctx context.Context, req *pb.GetTemplateRequest) (*pb.GetTemplateResponse, error) {
	if req == nil || req.TemplateId <= 0 {
		return &pb.GetTemplateResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.GetTemplateResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	resp, err := document_service.DocumentSvc.GetTemplateByID(req.TemplateId)
	if err != nil {
		return &pb.GetTemplateResponse{Base: baseResponseFromErr(err)}, nil
	}
	if req.RecordRecentLog {
		_ = document_service.DocumentSvc.CreateRecentTemplate(&dto.CreateRecentTemplateRequest{
			UserExternalID: userExternalID,
			TemplateID:     req.TemplateId,
			AccessTime:     time.Now(),
		})
	}
	return &pb.GetTemplateResponse{
		Base:     baseResponseFromErr(nil),
		Template: toTemplateResponse(resp),
	}, nil
}

func (s *DocumentServiceServer) ListRecentTemplates(ctx context.Context, req *pb.ListRecentTemplatesRequest) (*pb.ListRecentTemplatesResponse, error) {
	if req == nil {
		return &pb.ListRecentTemplatesResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.ListRecentTemplatesResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	var startTime *time.Time
	if req.StartTime != nil && *req.StartTime != "" {
		t, err := time.Parse(time.RFC3339, *req.StartTime)
		if err != nil {
			return &pb.ListRecentTemplatesResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
		}
		startTime = &t
	}
	var endTime *time.Time
	if req.EndTime != nil && *req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, *req.EndTime)
		if err != nil {
			return &pb.ListRecentTemplatesResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
		}
		endTime = &t
	}

	serviceReq := &dto.ListRecentTemplatesRequest{
		UserExternalID: userExternalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}
	templates, total, err := document_service.DocumentSvc.ListRecentTemplates(serviceReq)
	if err != nil {
		return &pb.ListRecentTemplatesResponse{Base: baseResponseFromErr(err)}, nil
	}
	respTemplates := make([]*pb.TemplateResponse, 0, len(templates))
	for _, tpl := range templates {
		respTemplates = append(respTemplates, toTemplateResponse(tpl))
	}
	return &pb.ListRecentTemplatesResponse{
		Base:      baseResponseFromErr(nil),
		Templates: respTemplates,
		Total:     total,
	}, nil
}

func (s *DocumentServiceServer) ListTemplates(ctx context.Context, req *pb.ListTemplatesRequest) (*pb.ListTemplatesResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.ListTemplatesResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if req == nil {
		return &pb.ListTemplatesResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	filterArgs := &dto.TemplateListFilterArgs{
		NameKeyword:     req.NameKeyword,
		KnowledgeBaseID: req.KnowledgeBaseId,
	}
	if req.TargetContainer != nil {
		container := fromCreateTemplateContainer(req.GetTargetContainer())
		filterArgs.TargetContainer = utils.Of(container)
	}

	sortArgs := dto.SortArgs{Field: "created_at", Desc: true}
	if req.OrderBy != nil {
		sortArgs.Field = req.GetOrderBy()
	}
	if req.OrderDesc != nil {
		sortArgs.Desc = req.GetOrderDesc()
	}

	serviceReq := &dto.TemplateListByExternalReq{
		FilterArgs: filterArgs,
		SortArgs:   sortArgs,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}
	if req.OnlyMine {
		serviceReq.CreatorExternalID = userExternalID
	}
	if req.TargetContainer != nil && req.GetTargetContainer() == pb.CreateTemplateContainer_OWN_TEMPLATE {
		serviceReq.CreatorExternalID = userExternalID
	}

	templates, total, err := document_service.DocumentSvc.ListTemplatesByExternal(serviceReq)
	if err != nil {
		return &pb.ListTemplatesResponse{Base: baseResponseFromErr(err)}, nil
	}

	respTemplates := make([]*pb.TemplateResponse, 0, len(templates))
	for _, tpl := range templates {
		respTemplates = append(respTemplates, toTemplateResponse(tpl))
	}

	return &pb.ListTemplatesResponse{
		Base:      baseResponseFromErr(nil),
		Templates: respTemplates,
		Total:     total,
	}, nil
}
