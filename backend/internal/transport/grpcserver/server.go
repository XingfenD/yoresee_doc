package grpcserver

import (
	"fmt"
	"net"
	"net/http"

	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
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

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return nil, fmt.Errorf("listen grpc port failed: %w", err)
	}

	go func() {
		_ = grpcServer.Serve(grpcListener)
	}()

	grpcWebServer := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	grpcWebHandler := grpcWebHTTPHandler(grpcWebServer)
	grpcWebAddr := fmt.Sprintf(":%d", grpcWebPort)
	go func() {
		_ = http.ListenAndServe(grpcWebAddr, grpcWebHandler)
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
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept,Accept-Language,Content-Type,Content-Length,X-Grpc-Web,X-User-Agent,Grpc-Timeout,Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
	w.Header().Set("Access-Control-Expose-Headers", "Grpc-Status,Grpc-Message,Grpc-Status-Details-Bin")
}
