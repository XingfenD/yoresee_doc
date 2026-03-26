package grpcserver

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/status"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(grpcPort, grpcWebPort int) (*grpc.Server, error) {
	allowUnauth := map[string]struct{}{
		"/yoresee_doc.v1.AuthService/Login":        {},
		"/yoresee_doc.v1.AuthService/Register":     {},
		"/yoresee_doc.v1.SystemService/Health":     {},
		"/yoresee_doc.v1.SystemService/SystemInfo": {},
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryAuthInterceptor(allowUnauth)),
	)

	pb.RegisterAuthServiceServer(grpcServer, NewAuthServiceServer())
	pb.RegisterDocumentServiceServer(grpcServer, NewDocumentServiceServer())
	pb.RegisterKnowledgeBaseServiceServer(grpcServer, NewKnowledgeBaseServiceServer())
	pb.RegisterSystemServiceServer(grpcServer, NewSystemServiceServer())
	pb.RegisterMembershipServiceServer(grpcServer, NewMembershipServiceServer())
	pb.RegisterInvitationServiceServer(grpcServer, NewInvitationServiceServer())

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return nil, status.GenErrWithCustomMsg(status.StatusServiceInternalError, "listen grpc port failed")
	}

	go func() {
		if serveErr := grpcServer.Serve(grpcListener); serveErr != nil {
			logrus.Errorf("grpc server exited with error: %v", serveErr)
		}
	}()

	grpcWebServer := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(isOriginAllowed),
	)

	grpcWebHandler := grpcWebHTTPHandler(grpcWebServer)
	grpcWebAddr := fmt.Sprintf(":%d", grpcWebPort)
	go func() {
		if serveErr := http.ListenAndServe(grpcWebAddr, grpcWebHandler); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			logrus.Errorf("grpc-web server exited with error: %v", serveErr)
		}
	}()

	return grpcServer, nil
}

func grpcWebHTTPHandler(wrapped *grpcweb.WrappedGrpcServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setGrpcWebCorsHeaders(w, r)
		if wrapped.IsGrpcWebRequest(r) || wrapped.IsGrpcWebSocketRequest(r) {
			wrapped.ServeHTTP(w, r)
			return
		}
		if r.Method == http.MethodOptions && wrapped.IsAcceptableGrpcCorsRequest(r) {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	})
}

func setGrpcWebCorsHeaders(w http.ResponseWriter, r *http.Request) {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if allowedOrigin := resolveAllowedOrigin(origin); allowedOrigin != "" {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	}
	w.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsAllowCredentials()))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsAllowedHeaders(), ","))
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsAllowedMethods(), ","))
	w.Header().Set("Access-Control-Expose-Headers", "Grpc-Status,Grpc-Message,Grpc-Status-Details-Bin")
	if maxAge := corsMaxAge(); maxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(maxAge))
	}
}

func isOriginAllowed(origin string) bool {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return false
	}
	allowedOrigins := corsAllowedOrigins()
	if len(allowedOrigins) == 0 {
		return true
	}
	for _, allowed := range allowedOrigins {
		allowed = strings.TrimSpace(allowed)
		if allowed == "*" || strings.EqualFold(allowed, origin) {
			return true
		}
	}
	return false
}

func resolveAllowedOrigin(origin string) string {
	if origin == "" {
		if len(corsAllowedOrigins()) == 0 {
			return "*"
		}
		return ""
	}
	if isOriginAllowed(origin) {
		return origin
	}
	return ""
}

func corsAllowedOrigins() []string {
	if config.GlobalConfig == nil {
		return nil
	}
	return config.GlobalConfig.Backend.Security.CORS.AllowedOrigins
}

func corsAllowedMethods() []string {
	if config.GlobalConfig != nil && len(config.GlobalConfig.Backend.Security.CORS.AllowedMethods) > 0 {
		return config.GlobalConfig.Backend.Security.CORS.AllowedMethods
	}
	return []string{"POST", "GET", "OPTIONS"}
}

func corsAllowedHeaders() []string {
	if config.GlobalConfig != nil && len(config.GlobalConfig.Backend.Security.CORS.AllowedHeaders) > 0 {
		return config.GlobalConfig.Backend.Security.CORS.AllowedHeaders
	}
	return []string{"Accept", "Accept-Language", "Content-Type", "Content-Length", "X-Grpc-Web", "X-User-Agent", "Grpc-Timeout", "Authorization"}
}

func corsAllowCredentials() bool {
	if config.GlobalConfig != nil {
		return config.GlobalConfig.Backend.Security.CORS.AllowCredentials
	}
	return true
}

func corsMaxAge() int {
	if config.GlobalConfig != nil {
		return config.GlobalConfig.Backend.Security.CORS.MaxAge
	}
	return 0
}
