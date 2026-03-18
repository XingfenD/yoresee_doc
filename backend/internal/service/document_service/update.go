package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *DocumentService) Update(ctx context.Context, req *dto.UpdateDocumentRequest) (bool, error) {
	if err := validateUpdateDocumentReq(req); err != nil {
		return false, err
	}

	var (
		moved   bool
		oldPath string
		newPath string
	)

	err := utils.WithTransaction(func(tx *gorm.DB) error {
		oldDoc, err := s.documentRepo.GetByExternalID(req.ExternalID).WithTx(tx).Exec(ctx)
		if err != nil {
			return status.StatusDocumentNotFound
		}
		oldPath = oldDoc.Path

		docModel := &model.Document{
			ID: oldDoc.ID,
		}
		op := s.documentRepo.Update(docModel).WithTx(tx)
		if req.Content != nil {
			docModel.Content = *req.Content
			op.UpdateContent()
		}
		if req.Title != nil {
			docModel.Title = *req.Title
		}

		var newParentID int64
		if req.ParentExternalID != nil {
			parentID, err := s.documentRepo.GetIDByExternalID(*req.ParentExternalID).Exec(ctx)
			if err != nil {
				return status.StatusDocumentNotFound
			}
			docModel.ParentID = parentID
			newParentID = parentID
			moved = true
			op = op.UpdateParentID()
		}

		// create version
		versionModel := &model.DocumentVersion{
			DocumentID:    oldDoc.ID,
			Title:         oldDoc.Title,
			Content:       oldDoc.Content,
			UserID:        oldDoc.UserID,
			ChangeSummary: "",
		}

		if err := s.docVersionRepo.Create(versionModel).WithTx(tx).Exec(); err != nil {
			return err
		}

		// update kb relation
		if req.KnowledgeBaseExternalID != nil && !req.MoveAsOwn {
			kbID, err := s.kbRepo.GetIDByExternalID(*req.KnowledgeBaseExternalID).Exec()
			if err != nil {
				return status.StatusKnowledgeBaseNotFound
			}
			docModel.KnowledgeID = &kbID
			op = op.UpdateKnowledgeID()
		}

		if req.MoveAsOwn {
			docModel.KnowledgeID = nil
			op = op.UpdateKnowledgeID()
		}

		if err := op.Exec(); err != nil {
			return err
		}
		if moved {
			if err := s.documentRepo.MoveSubtree(oldDoc.ID, newParentID).WithTx(tx).Exec(); err != nil {
				return err
			}
			path, err := s.documentRepo.GetPathByIDWithTx(tx, oldDoc.ID)
			if err != nil {
				return err
			}
			newPath = path
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	if moved {
		if err := s.documentRepo.BumpSubtreeVersionsByPath(ctx, oldPath); err != nil {
			logrus.Warnf("bump subtree version failed: %v", err)
			return true, nil
		}
		if newPath != "" && newPath != oldPath {
			if err := s.documentRepo.BumpSubtreeVersionsByPath(ctx, newPath); err != nil {
				logrus.Warnf("bump subtree version failed: %v", err)
				return true, nil
			}
		}
	}

	return true, nil
}
