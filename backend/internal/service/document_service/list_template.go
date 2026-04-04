package document_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/mapper/doc_type_mapper"
	"github.com/XingfenD/yoresee_doc/internal/model"
	internal_dto "github.com/XingfenD/yoresee_doc/internal/service/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/bytedance/gg/gslice"
	"github.com/sirupsen/logrus"
)

type templateListOperation struct {
	req  *internal_dto.TemplateListReq
	srvc *DocumentService
}

func (s *DocumentService) listTemplates(req *internal_dto.TemplateListReq) *templateListOperation {
	return &templateListOperation{
		req:  req,
		srvc: s,
	}
}

func (op *templateListOperation) ExecWithTotal() ([]*dto.TemplateResponse, int64, error) {
	listOp, err := op.srvc.buildListTemplateOperation(op.req)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] buildListTemplateOperation failed, err=%+v", err)
		return nil, 0, status.GenErrWithCustomMsg(err, "build template list operation failed")
	}
	templates, total, err := listOp.ExecWithTotal()
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] list templates failed, err=%+v", err)
		return nil, 0, status.GenErrWithCustomMsg(err, "list templates failed")
	}

	kbExternalIDMap := make(map[int64]string)
	kbIDs := make([]int64, 0)
	for _, tmpl := range templates {
		if tmpl.KnowledgeBaseID == nil || *tmpl.KnowledgeBaseID == 0 {
			continue
		}
		kbIDs = append(kbIDs, *tmpl.KnowledgeBaseID)
	}
	kbIDs = gslice.Uniq(kbIDs)
	if len(kbIDs) > 0 {
		kbs, err := op.srvc.kbRepo.MGetKnowledgeBaseByIDs(kbIDs).Exec()
		if err != nil {
			logrus.Errorf("[Service layer: DocumentService] list templates get knowledge bases failed: %+v", err)
			return nil, 0, status.StatusReadDBError
		}
		for _, kb := range kbs {
			kbExternalIDMap[kb.ID] = kb.ExternalID
		}
	}

	resp := make([]*dto.TemplateResponse, 0, len(templates))
	for _, tmpl := range templates {
		kbExternalID := ""
		if tmpl.KnowledgeBaseID != nil {
			kbExternalID = kbExternalIDMap[*tmpl.KnowledgeBaseID]
		}
		resp = append(resp, &dto.TemplateResponse{
			ID:                      tmpl.ID,
			Name:                    tmpl.Name,
			Description:             tmpl.Description,
			Content:                 tmpl.Content,
			Type:                    doc_type_mapper.FromModelType(tmpl.DocumentType),
			Scope:                   tmpl.Scope,
			KnowledgeBaseExternalID: kbExternalID,
			Tags:                    tmpl.Tags,
			CreatedAt:               tmpl.CreatedAt,
			UpdatedAt:               tmpl.UpdatedAt,
		})
	}

	return resp, total, nil
}

func (s *DocumentService) ListTemplatesByExternal(req *dto.TemplateListByExternalReq) ([]*dto.TemplateResponse, int64, error) {
	var creatorID *int64
	if req.CreatorExternalID != "" {
		id, err := s.userRepo.GetIDByExternalID(req.CreatorExternalID).Exec()
		if err != nil {
			return nil, 0, status.StatusUserNotFound
		}
		creatorID = &id
	}

	internalReq, err := buildTemplateListReqFromExternal(req, creatorID)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] buildTemplateListReqFromExternal failed, err=%+v", err)
		return nil, 0, status.GenErrWithCustomMsg(err, "build template list request failed")
	}
	list, total, err := s.listTemplates(internalReq).ExecWithTotal()
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] ListTemplatesByExternal failed, err=%+v", err)
		return nil, 0, status.GenErrWithCustomMsg(err, "list templates failed")
	}
	return list, total, nil
}

func buildTemplateListReqFromExternal(req *dto.TemplateListByExternalReq, creatorID *int64) (*internal_dto.TemplateListReq, error) {
	filter := &dto.TemplateListFilterArgs{}
	if req.FilterArgs != nil {
		*filter = *req.FilterArgs
	}

	return &internal_dto.TemplateListReq{
		CreatorID:  creatorID,
		FilterArgs: filter,
		SortArgs:   req.SortArgs,
		Pagination: req.Pagination,
	}, nil
}

func (s *DocumentService) buildListTemplateOperation(req *internal_dto.TemplateListReq) (*templateListOperationBuilder, error) {
	builder := &templateListOperationBuilder{
		req:  req,
		srvc: s,
	}
	return builder, nil
}

type templateListOperationBuilder struct {
	req  *internal_dto.TemplateListReq
	srvc *DocumentService
}

func (b *templateListOperationBuilder) ExecWithTotal() ([]*model.Template, int64, error) {
	op := b.srvc.templateRepo.List(&model.Template{})

	if b.req.CreatorID != nil {
		op = op.WithUserID(b.req.CreatorID)
	}

	if b.req.FilterArgs != nil {
		if b.req.FilterArgs.TargetContainer != nil {
			if *b.req.FilterArgs.TargetContainer == dto.TemplateContainerKnowledgeBase &&
				(b.req.FilterArgs.KnowledgeBaseID == nil || *b.req.FilterArgs.KnowledgeBaseID == "") {
				return nil, 0, status.StatusParamError
			}
			scope, _ := scopeFromContainer(*b.req.FilterArgs.TargetContainer)
			op = op.WithScope(&scope)
		}
		if b.req.FilterArgs.KnowledgeBaseID != nil && *b.req.FilterArgs.KnowledgeBaseID != "" {
			kbID, err := b.srvc.kbRepo.GetIDByExternalID(*b.req.FilterArgs.KnowledgeBaseID).Exec()
			if err != nil {
				logrus.Errorf("[Service layer: DocumentService] Get knowledge base id for template list failed, external_id=%s, err=%+v", *b.req.FilterArgs.KnowledgeBaseID, err)
				return nil, 0, status.StatusKnowledgeBaseNotFound
			}
			op = op.WithKnowledgeBaseID(&kbID)
		}
		if b.req.FilterArgs.NameKeyword != nil {
			op = op.WithNameKeyword(b.req.FilterArgs.NameKeyword)
		}
		if b.req.FilterArgs.Type != nil {
			modelType := model.DocumentType(doc_type_mapper.ToModelType(*b.req.FilterArgs.Type))
			op = op.WithDocumentType(&modelType)
		}
	}

	sortField := b.req.SortArgs.Field
	sortDesc := b.req.SortArgs.Desc
	if sortField == "" {
		sortField = "created_at"
		sortDesc = true
	}
	op = op.WithSort(sortField, sortDesc)

	if b.req.Pagination.Page > 0 && b.req.Pagination.PageSize > 0 {
		op = op.WithPagination(b.req.Pagination.Page, b.req.Pagination.PageSize)
	}

	templates, total, err := op.ExecWithTotal()
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] template repository list failed, err=%+v", err)
		return nil, 0, status.GenErrWithCustomMsg(status.StatusReadDBError, "list templates failed")
	}
	return templates, total, nil
}
