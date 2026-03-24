package grpcserver

import (
	"context"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service/auth_service"
	"github.com/XingfenD/yoresee_doc/internal/service/membership_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type MembershipServiceServer struct {
	pb.UnimplementedMembershipServiceServer
}

func NewMembershipServiceServer() *MembershipServiceServer {
	return &MembershipServiceServer{}
}

func (s *MembershipServiceServer) requireAdmin(ctx context.Context) error {
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

func (s *MembershipServiceServer) ListUserGroups(ctx context.Context, req *pb.ListUserGroupsRequest) (*pb.ListUserGroupsResponse, error) {
	if req == nil {
		return &pb.ListUserGroupsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.ListUserGroupsResponse{Base: baseResponseFromErr(err)}, nil
	}

	serviceReq := &dto.ListUserGroupsRequest{
		Keyword: req.Keyword,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}
	groups, total, err := membership_service.MembershipSvc.ListUserGroups(serviceReq)
	if err != nil {
		return &pb.ListUserGroupsResponse{Base: baseResponseFromErr(err)}, nil
	}

	respGroups := make([]*pb.UserGroupResponse, 0, len(groups))
	for _, group := range groups {
		respGroups = append(respGroups, toUserGroupResponse(group))
	}

	return &pb.ListUserGroupsResponse{
		Base:       baseResponseFromErr(nil),
		UserGroups: respGroups,
		Total:      total,
	}, nil
}

func (s *MembershipServiceServer) GetUserGroup(ctx context.Context, req *pb.GetUserGroupRequest) (*pb.GetUserGroupResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalId) == "" {
		return &pb.GetUserGroupResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.GetUserGroupResponse{Base: baseResponseFromErr(err)}, nil
	}

	group, err := membership_service.MembershipSvc.GetUserGroup(&dto.GetUserGroupRequest{
		ExternalID: req.ExternalId,
	})
	if err != nil {
		return &pb.GetUserGroupResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.GetUserGroupResponse{
		Base:      baseResponseFromErr(nil),
		UserGroup: toUserGroupResponse(group),
	}, nil
}

func (s *MembershipServiceServer) CreateUserGroup(ctx context.Context, req *pb.CreateUserGroupRequest) (*pb.CreateUserGroupResponse, error) {
	if req == nil || strings.TrimSpace(req.Name) == "" {
		return &pb.CreateUserGroupResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.CreateUserGroupResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.CreateUserGroupResponse{Base: baseResponseFromErr(err)}, nil
	}

	externalID, err := membership_service.MembershipSvc.CreateUserGroup(&dto.CreateUserGroupRequest{
		CreatorUserExternalID: userExternalID,
		Name:                  req.Name,
		Description:           req.Description,
		MemberUserExternalIDs: req.MemberUserExternalIds,
	})
	if err != nil {
		return &pb.CreateUserGroupResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.CreateUserGroupResponse{
		Base:       baseResponseFromErr(nil),
		ExternalId: externalID,
	}, nil
}

func (s *MembershipServiceServer) UpdateUserGroup(ctx context.Context, req *pb.UpdateUserGroupRequest) (*pb.UpdateUserGroupResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalId) == "" {
		return &pb.UpdateUserGroupResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.UpdateUserGroupResponse{Base: baseResponseFromErr(err)}, nil
	}

	err := membership_service.MembershipSvc.UpdateUserGroup(&dto.UpdateUserGroupRequest{
		ExternalID:            req.ExternalId,
		Name:                  req.Name,
		Description:           req.Description,
		SyncMembers:           req.SyncMembers,
		MemberUserExternalIDs: req.MemberUserExternalIds,
	})
	if err != nil {
		return &pb.UpdateUserGroupResponse{Base: baseResponseFromErr(err)}, nil
	}
	return &pb.UpdateUserGroupResponse{
		Base: baseResponseFromErr(nil),
	}, nil
}

func (s *MembershipServiceServer) DeleteUserGroup(ctx context.Context, req *pb.DeleteUserGroupRequest) (*pb.DeleteUserGroupResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalId) == "" {
		return &pb.DeleteUserGroupResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.DeleteUserGroupResponse{Base: baseResponseFromErr(err)}, nil
	}

	if err := membership_service.MembershipSvc.DeleteUserGroup(&dto.DeleteUserGroupRequest{
		ExternalID: req.ExternalId,
	}); err != nil {
		return &pb.DeleteUserGroupResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.DeleteUserGroupResponse{
		Base: baseResponseFromErr(nil),
	}, nil
}
