package comment_service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository/comment_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/constant"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/sirupsen/logrus"
)

type CommentService struct {
	commentRepo  *comment_repo.CommentRepository
	documentRepo *document_repo.DocumentRepository
	userRepo     *user_repo.UserRepository
}

func NewCommentService() *CommentService {
	return &CommentService{
		commentRepo:  comment_repo.CommentRepo,
		documentRepo: &document_repo.DocumentRepo,
		userRepo:     user_repo.UserRepo,
	}
}

func (s *CommentService) CreateComment(req *dto.CreateCommentRequest) (*model.DocumentComment, error) {
	if req == nil || strings.TrimSpace(req.DocumentExternalID) == "" || strings.TrimSpace(req.Content) == "" {
		return nil, status.StatusParamError
	}
	if strings.TrimSpace(req.CreatorExternalID) == "" {
		return nil, status.StatusTokenInvalid
	}
	doc, err := s.documentRepo.GetByExternalID(req.DocumentExternalID).Exec(context.Background())
	if err != nil || doc == nil {
		return nil, status.StatusDocumentNotFound
	}

	creatorID, err := s.userRepo.GetIDByExternalID(req.CreatorExternalID).Exec()
	if err != nil {
		return nil, status.StatusUserNotFound
	}

	parentID := int64(0)
	var parentComment *model.DocumentComment
	if req.ParentExternalID != nil && strings.TrimSpace(*req.ParentExternalID) != "" {
		parent, err := s.commentRepo.GetByExternalID(*req.ParentExternalID).Exec()
		if err != nil || parent == nil {
			return nil, status.StatusReadDBError
		}
		parentID = parent.ID
		parentComment = parent
	}

	item := &model.DocumentComment{
		ExternalID: utils.GenerateExternalID(utils.ExternalIDContextComment),
		DocumentID: doc.ID,
		ParentID:   parentID,
		CreatorID:  creatorID,
		Content:    req.Content,
	}
	if req.AnchorID != nil && strings.TrimSpace(*req.AnchorID) != "" {
		item.AnchorID = strings.TrimSpace(*req.AnchorID)
	}
	if err := s.commentRepo.Create(item).Exec(); err != nil {
		return nil, status.StatusWriteDBError
	}

	notifyQuote := ""
	if req.Quote != nil {
		notifyQuote = strings.TrimSpace(*req.Quote)
	}
	s.notifyCommentTargets(doc, item, parentComment, notifyQuote)
	return item, nil
}

func (s *CommentService) ListComments(req *dto.ListCommentsRequest) ([]model.DocumentComment, int64, error) {
	if req == nil || strings.TrimSpace(req.DocumentExternalID) == "" {
		return nil, 0, status.StatusParamError
	}
	doc, err := s.documentRepo.GetByExternalID(req.DocumentExternalID).Exec(context.Background())
	if err != nil || doc == nil {
		return nil, 0, status.StatusDocumentNotFound
	}
	op := s.commentRepo.ListByDocument(doc.ID).
		WithPagination(req.Page, req.PageSize)
	switch req.Scope {
	case dto.CommentScopeNormal:
		op = op.WithNormalOnly()
	case dto.CommentScopeInline:
		op = op.WithInlineOnly()
	}
	return op.ExecWithTotal()
}

func (s *CommentService) DeleteComment(req *dto.DeleteCommentRequest, isAdmin bool) error {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" {
		return status.StatusParamError
	}
	comment, err := s.commentRepo.GetByExternalID(req.ExternalID).Exec()
	if err != nil || comment == nil {
		return status.StatusReadDBError
	}
	if !isAdmin {
		operatorID, err := s.userRepo.GetIDByExternalID(req.OperatorExternalID).Exec()
		if err != nil {
			return status.StatusUserNotFound
		}
		if operatorID != comment.CreatorID {
			return status.StatusPermissionDenied
		}
	}
	if err := s.commentRepo.Delete(comment.ID).Exec(); err != nil {
		return status.StatusWriteDBError
	}
	return nil
}

func (s *CommentService) UpdateComment(req *dto.UpdateCommentRequest, isAdmin bool) (*model.DocumentComment, error) {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" || strings.TrimSpace(req.Content) == "" {
		return nil, status.StatusParamError
	}
	comment, err := s.commentRepo.GetByExternalID(req.ExternalID).Exec()
	if err != nil || comment == nil {
		return nil, status.StatusReadDBError
	}
	if !isAdmin {
		operatorID, err := s.userRepo.GetIDByExternalID(req.OperatorExternalID).Exec()
		if err != nil {
			return nil, status.StatusUserNotFound
		}
		if operatorID != comment.CreatorID {
			return nil, status.StatusPermissionDenied
		}
	}
	if err := s.commentRepo.UpdateContentByID(comment.ID, req.Content).Exec(); err != nil {
		return nil, status.StatusWriteDBError
	}
	comment.Content = req.Content
	return comment, nil
}

var CommentSvc = NewCommentService()

func (s *CommentService) notifyCommentTargets(doc *model.Document, comment *model.DocumentComment, parent *model.DocumentComment, quote string) {
	if doc == nil || comment == nil {
		return
	}
	receiverSet := map[string]struct{}{}

	if doc.UserID != comment.CreatorID {
		owner, err := s.userRepo.GetByID(doc.UserID).Exec()
		if err == nil && owner != nil && strings.TrimSpace(owner.ExternalID) != "" {
			receiverSet[owner.ExternalID] = struct{}{}
		}
	}

	parentExternalID := ""
	if parent != nil && parent.ID != 0 {
		parentExternalID = parent.ExternalID
		parentCreator, err := s.userRepo.GetByID(parent.CreatorID).Exec()
		if err == nil && parentCreator != nil && strings.TrimSpace(parentCreator.ExternalID) != "" {
			if parentCreator.ID != comment.CreatorID {
				receiverSet[parentCreator.ExternalID] = struct{}{}
			}
		}
	}

	if len(receiverSet) == 0 {
		return
	}

	receivers := make([]string, 0, len(receiverSet))
	for id := range receiverSet {
		receivers = append(receivers, id)
	}

	payload := map[string]any{
		"document_external_id":       doc.ExternalID,
		"comment_external_id":        comment.ExternalID,
		"document_title":             doc.Title,
		"parent_comment_external_id": parentExternalID,
	}
	payloadJSON, _ := json.Marshal(payload)

	evt := map[string]any{
		"receiver_external_ids": receivers,
		"type":                  "comment",
		"title":                 "新评论",
		"content":               comment.Content,
		"anchor_id":             comment.AnchorID,
		"quote":                 quote,
		"payload_json":          string(payloadJSON),
	}
	data, err := json.Marshal(evt)
	if err != nil {
		return
	}

	topic := constant.NotificationTopicDefault
	if err := mq.PublishTo(context.Background(), mq.BackendRabbitMQ, topic, data); err != nil {
		logrus.Errorf("publish comment notification failed: %v", err)
	}
}
