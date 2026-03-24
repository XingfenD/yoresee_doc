package grpcserver

import (
	"context"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"github.com/XingfenD/yoresee_doc/internal/service/auth_service"
	"github.com/XingfenD/yoresee_doc/internal/service/invitation_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type InvitationServiceServer struct {
	pb.UnimplementedInvitationServiceServer
	userRepo *user_repo.UserRepository
}

func NewInvitationServiceServer() *InvitationServiceServer {
	return &InvitationServiceServer{
		userRepo: user_repo.UserRepo,
	}
}

func (s *InvitationServiceServer) requireAdmin(ctx context.Context) error {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return status.StatusTokenInvalid
	}
	isAdmin, err := auth_service.AuthSvc.IsAdmin(userExternalID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return status.StatusPermissionDenied
	}
	return nil
}

func parseOptionalTime(input *string) (*time.Time, error) {
	if input == nil || strings.TrimSpace(*input) == "" {
		return nil, nil
	}
	tm, err := time.Parse(time.RFC3339, *input)
	if err != nil {
		return nil, err
	}
	return &tm, nil
}

func (s *InvitationServiceServer) CreateInvitation(ctx context.Context, req *pb.CreateInvitationRequest) (*pb.CreateInvitationResponse, error) {
	if req == nil {
		return &pb.CreateInvitationResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.CreateInvitationResponse{Base: baseResponseFromErr(err)}, nil
	}

	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.CreateInvitationResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}

	expiresAt, err := parseOptionalTime(req.ExpiresAt)
	if err != nil {
		return &pb.CreateInvitationResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	inv, err := invitation_service.InvitationSvc.CreateInvitation(&dto.CreateInvitationRequest{
		CreatorExternalID: userExternalID,
		MaxUsedCnt:        req.MaxUsedCnt,
		ExpiresAt:         expiresAt,
		Note:             req.Note,
	})
	if err != nil {
		return &pb.CreateInvitationResponse{Base: baseResponseFromErr(err)}, nil
	}

	user, _ := s.userRepo.GetByExternalID(userExternalID).Exec()
	resp := &dto.InvitationResponse{
		ID:                  inv.ID,
		Code:                inv.Code,
		CreatedByExternalID: userExternalID,
		CreatedByName:       "",
		UsedCnt:             inv.UsedCnt,
		MaxUsedCnt:          inv.MaxUsedCnt,
		ExpiresAt:           inv.ExpiresAt,
		CreatedAt:           inv.CreatedAt,
		Disabled:            inv.Disabled,
		Note:               inv.Note,
	}
	if user != nil {
		if user.Nickname != "" {
			resp.CreatedByName = user.Nickname
		} else {
			resp.CreatedByName = user.Username
		}
	}

	return &pb.CreateInvitationResponse{
		Base:       baseResponseFromErr(nil),
		Invitation: toInvitationResponse(resp),
	}, nil
}

func (s *InvitationServiceServer) ListInvitations(ctx context.Context, req *pb.ListInvitationsRequest) (*pb.ListInvitationsResponse, error) {
	if req == nil {
		return &pb.ListInvitationsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.ListInvitationsResponse{Base: baseResponseFromErr(err)}, nil
	}

	var creatorID *int64
	if req.CreatorExternalId != nil && strings.TrimSpace(req.GetCreatorExternalId()) != "" {
		id, err := s.userRepo.GetIDByExternalID(req.GetCreatorExternalId()).Exec()
		if err != nil {
			return &pb.ListInvitationsResponse{Base: baseResponseFromErr(status.StatusUserNotFound)}, nil
		}
		creatorID = utils.Of(id)
	}

	sortArgs := dto.SortArgs{Field: "created_at", Desc: true}
	if req.OrderBy != nil {
		sortArgs.Field = req.GetOrderBy()
	}
	if req.OrderDesc != nil {
		sortArgs.Desc = req.GetOrderDesc()
	}

	listReq := &dto.ListInvitationsReq{
		CreatorID:      creatorID,
		MaxUsedCnt:     req.MaxUsedCnt,
		ExpiresAtStart: req.ExpiresAtStart,
		ExpiresAtEnd:   req.ExpiresAtEnd,
		CreatedAtStart: req.CreatedAtStart,
		CreatedAtEnd:   req.CreatedAtEnd,
		Disabled:       req.Disabled,
		SortArgs:       sortArgs,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}

	invitations, total, err := invitation_service.InvitationSvc.ListInvitations(listReq)
	if err != nil {
		return &pb.ListInvitationsResponse{Base: baseResponseFromErr(err)}, nil
	}

	creatorIDs := make([]int64, 0, len(invitations))
	for _, inv := range invitations {
		creatorIDs = append(creatorIDs, inv.CreatedBy)
	}
	creatorMap, _ := s.userRepo.MGetUserByID(creatorIDs).Exec()

	respInvites := make([]*pb.InvitationResponse, 0, len(invitations))
	for _, inv := range invitations {
		createdByExternalID := ""
		createdByName := ""
		if creator := creatorMap[inv.CreatedBy]; creator != nil {
			createdByExternalID = creator.ExternalID
			if creator.Nickname != "" {
				createdByName = creator.Nickname
			} else {
				createdByName = creator.Username
			}
		}
		resp := &dto.InvitationResponse{
			ID:                  inv.ID,
			Code:                inv.Code,
			CreatedByExternalID: createdByExternalID,
			CreatedByName:       createdByName,
			UsedCnt:             inv.UsedCnt,
			MaxUsedCnt:          inv.MaxUsedCnt,
			ExpiresAt:           inv.ExpiresAt,
			CreatedAt:           inv.CreatedAt,
			Disabled:            inv.Disabled,
			Note:               inv.Note,
		}
		respInvites = append(respInvites, toInvitationResponse(resp))
	}

	return &pb.ListInvitationsResponse{
		Base:        baseResponseFromErr(nil),
		Invitations: respInvites,
		Total:       total,
	}, nil
}

func (s *InvitationServiceServer) UpdateInvitation(ctx context.Context, req *pb.UpdateInvitationRequest) (*pb.UpdateInvitationResponse, error) {
	if req == nil || strings.TrimSpace(req.Code) == "" {
		return &pb.UpdateInvitationResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.UpdateInvitationResponse{Base: baseResponseFromErr(err)}, nil
	}

	expiresAt, err := parseOptionalTime(req.ExpiresAt)
	if err != nil {
		return &pb.UpdateInvitationResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	if err := invitation_service.InvitationSvc.UpdateInvitation(&dto.UpdateInvitationRequest{
		Code:       req.Code,
		MaxUsedCnt: req.MaxUsedCnt,
		ExpiresAt:  expiresAt,
		Disabled:   req.Disabled,
		Note:       req.Note,
	}); err != nil {
		return &pb.UpdateInvitationResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.UpdateInvitationResponse{Base: baseResponseFromErr(nil)}, nil
}

func (s *InvitationServiceServer) DeleteInvitation(ctx context.Context, req *pb.DeleteInvitationRequest) (*pb.DeleteInvitationResponse, error) {
	if req == nil || strings.TrimSpace(req.Code) == "" {
		return &pb.DeleteInvitationResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.DeleteInvitationResponse{Base: baseResponseFromErr(err)}, nil
	}

	if err := invitation_service.InvitationSvc.DeleteInvitation(&dto.DeleteInvitationRequest{Code: req.Code}); err != nil {
		return &pb.DeleteInvitationResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.DeleteInvitationResponse{Base: baseResponseFromErr(nil)}, nil
}

func (s *InvitationServiceServer) ListInvitationRecords(ctx context.Context, req *pb.ListInvitationRecordsRequest) (*pb.ListInvitationRecordsResponse, error) {
	if req == nil {
		return &pb.ListInvitationRecordsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.ListInvitationRecordsResponse{Base: baseResponseFromErr(err)}, nil
	}

	listReq := &dto.ListInvitationRecordsRequest{
		Code:        req.Code,
		Status:      req.Status,
		UsedAtStart: req.UsedAtStart,
		UsedAtEnd:   req.UsedAtEnd,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}

	records, total, err := invitation_service.InvitationSvc.ListInvitationRecords(listReq)
	if err != nil {
		return &pb.ListInvitationRecordsResponse{Base: baseResponseFromErr(err)}, nil
	}

	userIDs := make([]int64, 0, len(records))
	for _, r := range records {
		if r.UsedByUserID != nil {
			userIDs = append(userIDs, *r.UsedByUserID)
		}
	}
	userMap, _ := s.userRepo.MGetUserByID(userIDs).Exec()

	respRecords := make([]*pb.InvitationRecordResponse, 0, len(records))
	for _, r := range records {
		usedBy := r.UsedBy
		usedByExternalID := ""
		if r.UsedByUserID != nil {
			if user := userMap[*r.UsedByUserID]; user != nil {
				usedByExternalID = user.ExternalID
				if user.Nickname != "" {
					usedBy = user.Nickname
				} else if user.Username != "" {
					usedBy = user.Username
				}
			}
		}
		resp := &dto.InvitationRecordResponse{
			ID:               r.ID,
			Code:             r.Code,
			UsedBy:           usedBy,
			UsedByExternalID: usedByExternalID,
			UsedAt:           r.UsedAt,
			Status:           r.Status,
		}
		respRecords = append(respRecords, toInvitationRecordResponse(resp))
	}

	return &pb.ListInvitationRecordsResponse{
		Base:    baseResponseFromErr(nil),
		Records: respRecords,
		Total:   total,
	}, nil
}
