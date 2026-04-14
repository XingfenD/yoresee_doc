package comment_service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository/comment_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/knowledge_base_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
)

type CommentService struct {
	commentRepo  *comment_repo.CommentRepository
	documentRepo *document_repo.DocumentRepository
	userRepo     *user_repo.UserRepository
	kbRepo       *knowledge_base_repo.KnowledgeBaseRepository
}

func NewCommentService() *CommentService {
	return &CommentService{
		commentRepo:  comment_repo.CommentRepo,
		documentRepo: &document_repo.DocumentRepo,
		userRepo:     user_repo.UserRepo,
		kbRepo:       knowledge_base_repo.KnowledgeBaseRepo,
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
	} else if parentComment != nil && strings.TrimSpace(parentComment.AnchorID) != "" {
		item.AnchorID = parentComment.AnchorID
	}
	if strings.TrimSpace(item.AnchorID) == "" {
		return nil, status.StatusParamError
	}
	if err := s.commentRepo.Create(item).Exec(); err != nil {
		return nil, status.StatusWriteDBError
	}

	notifyQuote := ""
	if req.Quote != nil {
		notifyQuote = strings.TrimSpace(*req.Quote)
	}
	s.notifyCommentTargets(doc, item, parentComment, notifyQuote, req.MentionedUserExternalIDs)
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
	list, total, err := op.ExecWithTotal()
	if err != nil {
		logrus.Errorf("[Service layer: CommentService] list comments failed, document_external_id=%s, err=%+v", req.DocumentExternalID, err)
		return nil, 0, status.GenErrWithCustomMsg(status.StatusReadDBError, "list comments failed")
	}
	return list, total, nil
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

func (s *CommentService) notifyCommentTargets(doc *model.Document, comment *model.DocumentComment, parent *model.DocumentComment, quote string, mentionedUserExternalIDs []string) {
	if doc == nil || comment == nil {
		return
	}

	// Resolve knowledge base external ID for navigation payload
	kbExternalID := ""
	if doc.KnowledgeID != nil && *doc.KnowledgeID != 0 {
		kb, err := s.kbRepo.GetByID(*doc.KnowledgeID).Exec()
		if err == nil && kb != nil {
			kbExternalID = kb.ExternalID
		}
	}

	parentExternalID := ""
	if parent != nil && parent.ID != 0 {
		parentExternalID = parent.ExternalID
	}

	payloadBase := map[string]any{
		"document_external_id":        doc.ExternalID,
		"comment_external_id":         comment.ExternalID,
		"document_title":              doc.Title,
		"parent_comment_external_id":  parentExternalID,
		"knowledge_base_external_id":  kbExternalID,
		"mentioned_user_external_ids": mentionedUserExternalIDs,
	}
	payloadJSON, _ := json.Marshal(payloadBase)
	payloadStr := string(payloadJSON)

	publish := func(receivers []string, notifType, title string) {
		if len(receivers) == 0 {
			return
		}
		evt := domain_event.NotificationCreateEvent{
			ReceiverExternalIDs: receivers,
			Type:                notifType,
			Title:               title,
			Content:             comment.Content,
			PayloadJSON:         payloadStr,
		}
		if err := domain_event.PublishNotificationCreateEvent(context.Background(), evt); err != nil {
			logrus.Errorf("publish %s notification failed: %v", notifType, err)
		}
	}

	// Build mention receiver set (highest priority — skip for other types)
	mentionSet := map[string]struct{}{}
	for _, id := range mentionedUserExternalIDs {
		id = strings.TrimSpace(id)
		if id != "" && id != comment.ExternalID {
			mentionSet[id] = struct{}{}
		}
	}

	// Build reply receivers (parent comment creator, excluding commenter and mention receivers)
	var replyReceivers []string
	if parent != nil && parent.ID != 0 {
		parentCreator, err := s.userRepo.GetByID(parent.CreatorID).Exec()
		if err == nil && parentCreator != nil {
			id := strings.TrimSpace(parentCreator.ExternalID)
			if id != "" && parentCreator.ID != comment.CreatorID {
				if _, isMentioned := mentionSet[id]; !isMentioned {
					replyReceivers = append(replyReceivers, id)
				}
			}
		}
	}

	// Build comment receivers (document owner, excluding commenter, reply receivers, and mention receivers)
	replySet := map[string]struct{}{}
	for _, id := range replyReceivers {
		replySet[id] = struct{}{}
	}
	var commentReceivers []string
	if doc.UserID != comment.CreatorID {
		owner, err := s.userRepo.GetByID(doc.UserID).Exec()
		if err == nil && owner != nil {
			id := strings.TrimSpace(owner.ExternalID)
			if id != "" {
				if _, isMentioned := mentionSet[id]; !isMentioned {
					if _, isReply := replySet[id]; !isReply {
						commentReceivers = append(commentReceivers, id)
					}
				}
			}
		}
	}

	mentionReceivers := make([]string, 0, len(mentionSet))
	for id := range mentionSet {
		mentionReceivers = append(mentionReceivers, id)
	}

	publish(mentionReceivers, constant.NotificationType_Mention, constant.GetNotificationTitle(constant.NotificationType_Mention))
	publish(replyReceivers, constant.NotificationType_Reply, constant.GetNotificationTitle(constant.NotificationType_Reply))
	publish(commentReceivers, constant.NotificationType_Comment, constant.GetNotificationTitle(constant.NotificationType_Comment))
}
