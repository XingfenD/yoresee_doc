package grpcserver

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthInterceptor(allowUnauth map[string]struct{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, ok := allowUnauth[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		authHeader := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get("authorization"); len(values) > 0 {
				authHeader = values[0]
			} else if values := md.Get("Authorization"); len(values) > 0 {
				authHeader = values[0]
			}
		}

		claims, err := middleware.JWTAuth.ValidateAuthorizationHeader(authHeader)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		ctx = context.WithValue(ctx, "user_external_id", claims.ExternalID)
		ctx = context.WithValue(ctx, "username", claims.Username)

		return handler(ctx, req)
	}
}
