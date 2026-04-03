package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/cache"
	"github.com/XingfenD/yoresee_doc/pkg/key"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *DocumentService) Update(ctx context.Context, req *dto.UpdateDocumentRequest) (bool, error) {
	if err := validateUpdateDocumentReq(req); err != nil {
		logrus.Errorf("[Service layer: DocumentService] validateUpdateDocumentReq failed, err=%+v", err)
		return false, status.GenErrWithCustomMsg(err, "invalid update document request")
	}

	var (
		moved   bool
		oldPath string
		newPath string
	)

	cacheKey := key.KeyModelByExternalID(key.KeyObjectTypeEnum_Doc, req.ExternalID)
	err := cache.DoubleDelete(
		context.Background(),
		func() error {
			return utils.WithTransaction(func(tx *gorm.DB) error {
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
					logrus.Errorf("[Service layer: DocumentService] create document version failed, document_id=%d, err=%+v", oldDoc.ID, err)
					return status.GenErrWithCustomMsg(status.StatusWriteDBError, "create document version failed")
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
					logrus.Errorf("[Service layer: DocumentService] update document failed, external_id=%s, err=%+v", req.ExternalID, err)
					return status.GenErrWithCustomMsg(status.StatusWriteDBError, "update document failed")
				}
				if moved {
					if err := s.documentRepo.MoveSubtree(oldDoc.ID, newParentID).WithTx(tx).Exec(); err != nil {
						logrus.Errorf("[Service layer: DocumentService] move subtree failed, document_id=%d, new_parent_id=%d, err=%+v", oldDoc.ID, newParentID, err)
						return status.GenErrWithCustomMsg(status.StatusWriteDBError, "move document subtree failed")
					}
					path, err := s.documentRepo.GetPathByIDWithTx(tx, oldDoc.ID)
					if err != nil {
						logrus.Errorf("[Service layer: DocumentService] query new path failed, document_id=%d, err=%+v", oldDoc.ID, err)
						return status.GenErrWithCustomMsg(status.StatusReadDBError, "query document path failed")
					}
					newPath = path
				}

				return nil
			})
		},
		cacheKey,
	)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] Update failed, external_id=%s, err=%+v", req.ExternalID, err)
		return false, status.GenErrWithCustomMsg(err, "update document failed")
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

	if err := domain_event.PublishDocumentUpsertEvent(ctx, req.ExternalID); err != nil {
		logrus.Warnf("[Service layer: DocumentService] publish search sync event failed, external_id=%s, err=%+v", req.ExternalID, err)
	}

	return true, nil
}
