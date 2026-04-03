package document_service

import (
	"context"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/mapper/doc_container_mapper"
	"github.com/XingfenD/yoresee_doc/internal/mapper/doc_type_mapper"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *DocumentService) Create(ctx context.Context, req *dto.CreateDocumentReq) (*dto.CreateDocumentResponse, error) {
	if err := validateCreateDocumentReq(req); err != nil {
		logrus.Errorf("[Service layer: DocumentService] validateCreateDocumentReq failed, err=%+v", err)
		return nil, status.GenErrWithCustomMsg(err, "invalid create document request")
	}

	// TODO: redis support
	docExternalID := utils.GenerateExternalID(utils.ExternalIDContextDocument)
	err := utils.WithTransaction(func(tx *gorm.DB) error {
		docModel := &model.Document{
			ExternalID:    docExternalID,
			Title:         req.Title,
			Type:          doc_type_mapper.ToModelType(req.Type),
			ContainerType: doc_container_mapper.ToModelType(req.ContainerType),
			Summary:       "",
			IsPublic:      req.IsPublic,
			Content:       "",
			Path:          "0",
			Depth:         0,
		}

		// query user_id
		userID, err := s.userRepo.GetIDByExternalID(*req.CreatorExternalID).WithTx(tx).Exec()
		if err != nil {
			return status.StatusUserNotFound
		}

		docModel.UserID = userID

		// query knowledge_id
		if req.ContainerType == dto.ContainerType_KnowledgeBase {
			kbID, err := s.kbRepo.GetIDByExternalID(*req.KnowledgeExternalID).WithTx(tx).Exec()
			if err != nil {
				return status.StatusKnowledgeBaseNotFound
			}
			docModel.KnowledgeID = &kbID
		}

		// TODO: permission check for knowledgebase

		// query parent_id
		if req.ParentExternalID != nil {
			parentDocID, err := s.documentRepo.GetIDByExternalID(*req.ParentExternalID).WithTx(tx).Exec(ctx)
			if err != nil {
				return status.StatusDocumentNotFound
			}
			docModel.ParentID = parentDocID
		}

		// apply template content if provided
		if req.TemplateID != nil && *req.TemplateID > 0 {
			tpl, err := s.templateRepo.GetByID(*req.TemplateID).WithTx(tx).Exec()
			if err != nil || tpl == nil {
				logrus.Errorf("template not found: err:%+v", err)
				return status.GenErrWithCustomMsg(status.StatusParamError, "template not found")
			}
			docModel.Content = tpl.Content
		}

		// create document meta
		err = s.documentRepo.Create(docModel).WithTx(tx).Exec()
		if err != nil {
			status.GenErrWithCustomMsg(status.StatusWriteDBError, "create document meta failed")
			return status.StatusWriteDBError
		}
		if err := s.documentRepo.UpdatePathDepth(docModel.ID, docModel.ParentID).WithTx(tx).Exec(); err != nil {
			return status.GenErrWithCustomMsg(status.StatusWriteDBError, "update document path depth failed")
		}

		// create doc version
		ver := &model.DocumentVersion{
			DocumentID:    docModel.ID,
			Content:       docModel.Content,
			UserID:        userID,
			Title:         docModel.Title,
			ChangeSummary: "Create the document",
		}
		err = s.docVersionRepo.Create(ver).WithTx(tx).Exec()
		if err != nil {
			return status.GenErrWithCustomMsg(status.StatusWriteDBError, "create version failed")
		}

		return nil
	})

	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] Create err: %+v", err)
		return nil, status.StatusWriteDBError
	}

	if createdDoc, err := s.documentRepo.GetByExternalID(docExternalID).Exec(ctx); err == nil {
		if err := s.documentRepo.BumpSubtreeVersionsByPath(ctx, createdDoc.Path); err != nil {
			logrus.Warnf("bump subtree version failed: %v", err)
		}
	}
	if req.TemplateID != nil && req.CreatorExternalID != nil && *req.CreatorExternalID != "" {
		_ = s.CreateRecentTemplate(&dto.CreateRecentTemplateRequest{
			UserExternalID: *req.CreatorExternalID,
			TemplateID:     *req.TemplateID,
			AccessTime:     time.Now(),
		})
	}
	if err := domain_event.PublishDocumentUpsertEvent(ctx, docExternalID); err != nil {
		logrus.Warnf("[Service layer: DocumentService] publish search sync event failed, external_id=%s, err=%+v", docExternalID, err)
	}

	return &dto.CreateDocumentResponse{
		ExternalID: docExternalID,
	}, nil
}
