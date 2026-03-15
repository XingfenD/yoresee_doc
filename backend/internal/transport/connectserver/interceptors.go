package connectserver

import (
	"context"

	"connectrpc.com/connect"
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/middleware"
)

func UnaryAuthInterceptor(allowUnauth map[string]struct{}) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if _, ok := allowUnauth[req.Spec().Procedure]; ok {
				return next(ctx, req)
			}

			internalKey := req.Header().Get("x-internal-key")
			if internalKey != "" && config.GlobalConfig.Backend.InternalRPCKey != "" && internalKey == config.GlobalConfig.Backend.InternalRPCKey {
				return next(ctx, req)
			}

			authHeader := req.Header().Get("Authorization")
			if authHeader == "" {
				authHeader = req.Header().Get("authorization")
			}

			claims, err := middleware.JWTAuth.ValidateAuthorizationHeader(authHeader)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			ctx = context.WithValue(ctx, "user_external_id", claims.ExternalID)
			ctx = context.WithValue(ctx, "username", claims.Username)

			return next(ctx, req)
		}
	})
}
