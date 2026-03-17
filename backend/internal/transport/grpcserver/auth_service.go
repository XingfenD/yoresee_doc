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
