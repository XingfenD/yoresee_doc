package grpcserver

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service/auth_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthServiceServer() *AuthServiceServer {
	return &AuthServiceServer{}
}

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.AuthLoginRequest) (*pb.AuthLoginResponse, error) {
	if req == nil || req.Email == "" || req.Password == "" {
		return &pb.AuthLoginResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	token, user, err := auth_service.AuthSvc.Login(req.Email, req.Password)
	if err != nil {
		return &pb.AuthLoginResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.AuthLoginResponse{
		Base:  baseResponseFromErr(nil),
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.AuthRegisterRequest) (*pb.AuthRegisterResponse, error) {
	if req == nil || req.Username == "" || req.Password == "" || req.Email == "" {
		return &pb.AuthRegisterResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	userCreate := &dto.UserCreate{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	if req.InvitationCode != nil {
		userCreate.InvitationCode = req.InvitationCode
	}

	err := auth_service.AuthSvc.Register(ctx, userCreate)
	if err != nil {
		return &pb.AuthRegisterResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.AuthRegisterResponse{Base: baseResponseFromErr(nil)}, nil
}

func (s *AuthServiceServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	if req == nil {
		return &pb.UpdateProfileResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.UpdateProfileResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}

	serviceReq := &dto.UpdateProfileRequest{
		Username:          req.Username,
		Email:             req.Email,
		Nickname:          req.Nickname,
		Password:          req.Password,
		Avatar:            req.Avatar,
		AvatarFile:        req.AvatarFile,
		AvatarFilename:    req.AvatarFilename,
		AvatarContentType: req.AvatarContentType,
	}

	user, err := auth_service.AuthSvc.UpdateProfile(userExternalID, serviceReq)
	if err != nil {
		return &pb.UpdateProfileResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.UpdateProfileResponse{
		Base: baseResponseFromErr(nil),
		User: toUserResponse(user),
	}, nil
}

func (s *AuthServiceServer) QuerySideBarDisplay(ctx context.Context, req *pb.QuerySideBarDisplayRequest) (*pb.QuerySideBarDisplayResponse, error) {
	if req == nil || req.Scene == "" {
		return &pb.QuerySideBarDisplayResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.QuerySideBarDisplayResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}

	isAdmin, err := auth_service.AuthSvc.IsAdmin(userExternalID)
	if err != nil {
		return &pb.QuerySideBarDisplayResponse{Base: baseResponseFromErr(err)}, nil
	}

	displayTabs, err := auth_service.AuthSvc.QuerySideBarDisplay(req.Scene, isAdmin)
	if err != nil {
		return &pb.QuerySideBarDisplayResponse{Base: baseResponseFromErr(err)}, nil
	}
	return &pb.QuerySideBarDisplayResponse{
		Base:        baseResponseFromErr(nil),
		DisplayTabs: displayTabs,
	}, nil
}

func (s *AuthServiceServer) QueryTopNavDisplay(ctx context.Context, req *pb.QueryTopNavDisplayRequest) (*pb.QueryTopNavDisplayResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return &pb.QueryTopNavDisplayResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}

	isAdmin, err := auth_service.AuthSvc.IsAdmin(userExternalID)
	if err != nil {
		return &pb.QueryTopNavDisplayResponse{Base: baseResponseFromErr(err)}, nil
	}

	menus, err := auth_service.AuthSvc.QueryTopNavDisplay(isAdmin)
	if err != nil {
		return &pb.QueryTopNavDisplayResponse{Base: baseResponseFromErr(err)}, nil
	}

	return &pb.QueryTopNavDisplayResponse{
		Base:         baseResponseFromErr(nil),
		DisplayMenus: menus,
	}, nil
}
