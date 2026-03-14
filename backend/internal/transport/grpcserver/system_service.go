package grpcserver

import (
	"context"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/service"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type SystemServiceServer struct {
	pb.UnimplementedSystemServiceServer
}

func NewSystemServiceServer() *SystemServiceServer {
	return &SystemServiceServer{}
}

func (s *SystemServiceServer) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Base:      baseResponseFromErr(nil),
		Timestamp: time.Now().Format(time.RFC3339),
		Status:    "healthy",
		Version:   "1.0.0",
	}, nil
}

func (s *SystemServiceServer) SystemInfo(ctx context.Context, req *pb.SystemInfoRequest) (*pb.SystemInfoResponse, error) {
	return &pb.SystemInfoResponse{
		Base:               baseResponseFromErr(nil),
		SystemName:         config.GlobalConfig.Backend.SystemName,
		SystemRegisterMode: service.ConfigSvc.GetSystemRegisterMode(ctx),
	}, nil
}
