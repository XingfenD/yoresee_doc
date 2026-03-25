package grpcserver

import (
	"context"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"github.com/XingfenD/yoresee_doc/internal/service/auth_service"
	"github.com/XingfenD/yoresee_doc/internal/service/comment_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type CommentServiceServer struct {
	pb.UnimplementedCommentServiceServer
	userRepo *user_repo.UserRepository
}

func NewCommentServiceServer() *CommentServiceServer {
	return &CommentServiceServer{
		userRepo: user_repo.UserRepo,
	}
}

func (s *CommentServiceServer) CreateDocumentComment(ctx context.Context, req *pb.CreateDocumentCommentRequest) (*pb.CreateDocumentCommentResponse, error) {
	if req == nil {
		return &pb.CreateDocumentCommentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || strings.TrimSpace(userExternalID) == "" {
		return &pb.CreateDocumentCommentResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}

	var parentExternalID *string
	if req.ParentExternalId != nil && strings.TrimSpace(req.GetParentExternalId()) != "" {
		parentExternalID = req.ParentExternalId
	}

	comment, err := comment_service.CommentSvc.CreateComment(&dto.CreateCommentRequest{
		DocumentExternalID: req.DocumentExternalId,
		Content:            req.Content,
		ParentExternalID:   parentExternalID,
		CreatorExternalID:  userExternalID,
	})
	if err != nil {
		return &pb.CreateDocumentCommentResponse{Base: baseResponseFromErr(err)}, nil
	}

	creator, _ := s.userRepo.GetByExternalID(userExternalID).Exec()
	resp := toCommentResponse(comment, req.DocumentExternalId, req.GetParentExternalId(), creator)
	return &pb.CreateDocumentCommentResponse{Base: baseResponseFromErr(nil), Comment: resp}, nil
}

func (s *CommentServiceServer) ListDocumentComments(ctx context.Context, req *pb.ListDocumentCommentsRequest) (*pb.ListDocumentCommentsResponse, error) {
	if req == nil {
		return &pb.ListDocumentCommentsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	comments, total, err := comment_service.CommentSvc.ListComments(&dto.ListCommentsRequest{
		DocumentExternalID: req.DocumentExternalId,
		Page:               int(req.Page),
		PageSize:           int(req.PageSize),
	})
	if err != nil {
		return &pb.ListDocumentCommentsResponse{Base: baseResponseFromErr(err)}, nil
	}

	creatorIDs := make([]int64, 0, len(comments))
	for _, item := range comments {
		creatorIDs = append(creatorIDs, item.CreatorID)
	}
	creatorMap, _ := s.userRepo.MGetUserByID(creatorIDs).Exec()

	parentMap := make(map[int64]string, len(comments))
	for _, item := range comments {
		parentMap[item.ID] = item.ExternalID
	}

	respItems := make([]*pb.DocumentComment, 0, len(comments))
	for _, item := range comments {
		parentExternalID := ""
		if item.ParentID != 0 {
			parentExternalID = parentMap[item.ParentID]
		}
		respItems = append(respItems, toCommentResponse(&item, req.DocumentExternalId, parentExternalID, creatorMap[item.CreatorID]))
	}

	return &pb.ListDocumentCommentsResponse{
		Base:     baseResponseFromErr(nil),
		Comments: respItems,
		Total:    total,
	}, nil
}

func (s *CommentServiceServer) DeleteDocumentComment(ctx context.Context, req *pb.DeleteDocumentCommentRequest) (*pb.DeleteDocumentCommentResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalId) == "" {
		return &pb.DeleteDocumentCommentResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || strings.TrimSpace(userExternalID) == "" {
		return &pb.DeleteDocumentCommentResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}
	isAdmin, _ := auth_service.AuthSvc.IsAdmin(userExternalID)
	if err := comment_service.CommentSvc.DeleteComment(&dto.DeleteCommentRequest{
		ExternalID:         req.ExternalId,
		OperatorExternalID: userExternalID,
	}, isAdmin); err != nil {
		return &pb.DeleteDocumentCommentResponse{Base: baseResponseFromErr(err)}, nil
	}
	return &pb.DeleteDocumentCommentResponse{Base: baseResponseFromErr(nil)}, nil
}

func toCommentResponse(comment *model.DocumentComment, docExternalID, parentExternalID string, creator *model.User) *pb.DocumentComment {
	if comment == nil {
		return nil
	}
	resp := &pb.DocumentComment{
		ExternalId:            comment.ExternalID,
		DocumentExternalId:    docExternalID,
		ParentExternalId:      parentExternalID,
		Content:               comment.Content,
		CreatorUserExternalId: "",
		CreatorName:           "",
		CreatorAvatar:         "",
		CreatedAt:             timeToString(comment.CreatedAt),
	}
	if creator != nil {
		resp.CreatorUserExternalId = creator.ExternalID
		if creator.Nickname != "" {
			resp.CreatorName = creator.Nickname
		} else {
			resp.CreatorName = creator.Username
		}
		resp.CreatorAvatar = creator.Avatar
	}
	return resp
}
