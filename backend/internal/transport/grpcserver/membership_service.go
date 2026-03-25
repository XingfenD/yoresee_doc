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

func (s *MembershipServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	if req == nil {
		return &pb.ListUsersResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.ListUsersResponse{Base: baseResponseFromErr(err)}, nil
	}

	serviceReq := &dto.ListUsersRequest{
		Keyword: req.Keyword,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}
	users, total, err := membership_service.MembershipSvc.ListUsers(serviceReq)
	if err != nil {
		return &pb.ListUsersResponse{Base: baseResponseFromErr(err)}, nil
	}

	respUsers := make([]*pb.UserResponse, 0, len(users))
	for _, user := range users {
		respUsers = append(respUsers, toUserResponse(user))
	}

	return &pb.ListUsersResponse{
		Base:  baseResponseFromErr(nil),
		Users: respUsers,
		Total: total,
	}, nil
}

func (s *MembershipServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalId) == "" {
		return &pb.UpdateUserResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.UpdateUserResponse{Base: baseResponseFromErr(err)}, nil
	}

	serviceReq := &dto.UpdateUserRequest{
		ExternalID: req.ExternalId,
		Username:   req.Username,
		Email:      req.Email,
		Nickname:   req.Nickname,
		Status:     req.Status,
	}
	if err := membership_service.MembershipSvc.UpdateUser(serviceReq); err != nil {
		return &pb.UpdateUserResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.UpdateUserResponse{
		Base: baseResponseFromErr(nil),
	}, nil
}

func (s *MembershipServiceServer) ListUserGroupMembers(ctx context.Context, req *pb.ListUserGroupMembersRequest) (*pb.ListUserGroupMembersResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalId) == "" {
		return &pb.ListUserGroupMembersResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.ListUserGroupMembersResponse{Base: baseResponseFromErr(err)}, nil
	}

	serviceReq := &dto.ListUserGroupMembersRequest{
		ExternalID: req.ExternalId,
		Keyword:    req.Keyword,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}
	users, total, err := membership_service.MembershipSvc.ListUserGroupMembers(serviceReq)
	if err != nil {
		return &pb.ListUserGroupMembersResponse{Base: baseResponseFromErr(err)}, nil
	}

	respUsers := make([]*pb.UserResponse, 0, len(users))
	for _, user := range users {
		respUsers = append(respUsers, toUserResponse(user))
	}

	return &pb.ListUserGroupMembersResponse{
		Base:  baseResponseFromErr(nil),
		Users: respUsers,
		Total: total,
	}, nil
}

func (s *MembershipServiceServer) ListOrgNodes(ctx context.Context, req *pb.ListOrgNodesRequest) (*pb.ListOrgNodesResponse, error) {
	if req == nil {
		return &pb.ListOrgNodesResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.ListOrgNodesResponse{Base: baseResponseFromErr(err)}, nil
	}

	serviceReq := &dto.ListOrgNodesRequest{
		ParentExternalID: req.ParentExternalId,
		Keyword:          req.Keyword,
		IncludeChildren:  req.IncludeChildren,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}
	nodes, total, err := membership_service.MembershipSvc.ListOrgNodes(serviceReq)
	if err != nil {
		return &pb.ListOrgNodesResponse{Base: baseResponseFromErr(err)}, nil
	}

	respNodes := make([]*pb.OrgNodeResponse, 0, len(nodes))
	for _, node := range nodes {
		respNodes = append(respNodes, toOrgNodeResponse(node))
	}

	return &pb.ListOrgNodesResponse{
		Base:     baseResponseFromErr(nil),
		OrgNodes: respNodes,
		Total:    total,
	}, nil
}

func (s *MembershipServiceServer) GetOrgNode(ctx context.Context, req *pb.GetOrgNodeRequest) (*pb.GetOrgNodeResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalId) == "" {
		return &pb.GetOrgNodeResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.GetOrgNodeResponse{Base: baseResponseFromErr(err)}, nil
	}

	node, err := membership_service.MembershipSvc.GetOrgNode(&dto.GetOrgNodeRequest{
		ExternalID:      req.ExternalId,
		IncludeChildren: req.IncludeChildren,
	})
	if err != nil {
		return &pb.GetOrgNodeResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.GetOrgNodeResponse{
		Base:    baseResponseFromErr(nil),
		OrgNode: toOrgNodeResponse(node),
	}, nil
}

func (s *MembershipServiceServer) CreateOrgNode(ctx context.Context, req *pb.CreateOrgNodeRequest) (*pb.CreateOrgNodeResponse, error) {
	if req == nil || strings.TrimSpace(req.Name) == "" {
		return &pb.CreateOrgNodeResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.CreateOrgNodeResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.CreateOrgNodeResponse{Base: baseResponseFromErr(err)}, nil
	}

	externalID, err := membership_service.MembershipSvc.CreateOrgNode(&dto.CreateOrgNodeRequest{
		CreatorUserExternalID: userExternalID,
		ParentExternalID:      req.ParentExternalId,
		Name:                  req.Name,
		Description:           req.Description,
		MemberUserExternalIDs: req.MemberUserExternalIds,
	})
	if err != nil {
		return &pb.CreateOrgNodeResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.CreateOrgNodeResponse{
		Base:       baseResponseFromErr(nil),
		ExternalId: externalID,
	}, nil
}

func (s *MembershipServiceServer) UpdateOrgNode(ctx context.Context, req *pb.UpdateOrgNodeRequest) (*pb.UpdateOrgNodeResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalId) == "" {
		return &pb.UpdateOrgNodeResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.UpdateOrgNodeResponse{Base: baseResponseFromErr(err)}, nil
	}

	err := membership_service.MembershipSvc.UpdateOrgNode(&dto.UpdateOrgNodeRequest{
		ExternalID:            req.ExternalId,
		Name:                  req.Name,
		Description:           req.Description,
		SyncMembers:           req.SyncMembers,
		MemberUserExternalIDs: req.MemberUserExternalIds,
	})
	if err != nil {
		return &pb.UpdateOrgNodeResponse{Base: baseResponseFromErr(err)}, nil
	}
	return &pb.UpdateOrgNodeResponse{
		Base: baseResponseFromErr(nil),
	}, nil
}

func (s *MembershipServiceServer) DeleteOrgNode(ctx context.Context, req *pb.DeleteOrgNodeRequest) (*pb.DeleteOrgNodeResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalId) == "" {
		return &pb.DeleteOrgNodeResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.DeleteOrgNodeResponse{Base: baseResponseFromErr(err)}, nil
	}

	if err := membership_service.MembershipSvc.DeleteOrgNode(&dto.DeleteOrgNodeRequest{
		ExternalID: req.ExternalId,
	}); err != nil {
		return &pb.DeleteOrgNodeResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.DeleteOrgNodeResponse{
		Base: baseResponseFromErr(nil),
	}, nil
}

func (s *MembershipServiceServer) MoveOrgNode(ctx context.Context, req *pb.MoveOrgNodeRequest) (*pb.MoveOrgNodeResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalId) == "" {
		return &pb.MoveOrgNodeResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.MoveOrgNodeResponse{Base: baseResponseFromErr(err)}, nil
	}

	if err := membership_service.MembershipSvc.MoveOrgNode(&dto.MoveOrgNodeRequest{
		ExternalID:          req.ExternalId,
		NewParentExternalID: req.NewParentExternalId,
	}); err != nil {
		return &pb.MoveOrgNodeResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.MoveOrgNodeResponse{
		Base: baseResponseFromErr(nil),
	}, nil
}

func (s *MembershipServiceServer) ListOrgNodeMembers(ctx context.Context, req *pb.ListOrgNodeMembersRequest) (*pb.ListOrgNodeMembersResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalId) == "" {
		return &pb.ListOrgNodeMembersResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.ListOrgNodeMembersResponse{Base: baseResponseFromErr(err)}, nil
	}

	serviceReq := &dto.ListOrgNodeMembersRequest{
		ExternalID: req.ExternalId,
		Keyword:    req.Keyword,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}
	users, total, err := membership_service.MembershipSvc.ListOrgNodeMembers(serviceReq)
	if err != nil {
		return &pb.ListOrgNodeMembersResponse{Base: baseResponseFromErr(err)}, nil
	}

	respUsers := make([]*pb.UserResponse, 0, len(users))
	for _, user := range users {
		respUsers = append(respUsers, toUserResponse(user))
	}

	return &pb.ListOrgNodeMembersResponse{
		Base:  baseResponseFromErr(nil),
		Users: respUsers,
		Total: total,
	}, nil
}

func toOrgNodeResponse(node *dto.OrgNodeResponse) *pb.OrgNodeResponse {
	resp := &pb.OrgNodeResponse{
		ExternalId:            node.ExternalID,
		ParentExternalId:      node.ParentExternalID,
		Name:                  node.Name,
		Path:                  node.Path,
		Description:           node.Description,
		CreatorUserExternalId: node.CreatorUserExternalID,
		MemberCount:           int32(node.MemberCount),
	}
	if node.Children != nil {
		resp.Children = make([]*pb.OrgNodeResponse, 0, len(node.Children))
		for _, child := range node.Children {
			resp.Children = append(resp.Children, toOrgNodeResponse(child))
		}
	}
	return resp
}
