package connectserver

import (
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Start(grpcPort, grpcWebPort int) (*http.Server, *http.Server, error) {
	allowUnauth := map[string]struct{}{
		pb.AuthService_Login_FullMethodName:        {},
		pb.AuthService_Register_FullMethodName:     {},
		pb.SystemService_Health_FullMethodName:     {},
		pb.SystemService_SystemInfo_FullMethodName: {},
	}

	interceptor := UnaryAuthInterceptor(allowUnauth)
	opts := []connect.HandlerOption{
		connect.WithInterceptors(interceptor),
	}

	mux := http.NewServeMux()
	registerHandlers(mux, opts)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := withCORS(mux)
	h2cHandler := h2c.NewHandler(handler, &http2.Server{})

	grpcAddr := fmt.Sprintf(":%d", grpcPort)
	grpcServer := &http.Server{Addr: grpcAddr, Handler: h2cHandler}
	go func() {
		_ = grpcServer.ListenAndServe()
	}()

	grpcWebAddr := fmt.Sprintf(":%d", grpcWebPort)
	grpcWebServer := &http.Server{Addr: grpcWebAddr, Handler: h2cHandler}
	go func() {
		_ = grpcWebServer.ListenAndServe()
	}()

	return grpcServer, grpcWebServer, nil
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setConnectCorsHeaders(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func setConnectCorsHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept,Accept-Language,Content-Type,Content-Length,X-Grpc-Web,X-User-Agent,Grpc-Timeout,Authorization,Connect-Protocol-Version")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
	w.Header().Set("Access-Control-Expose-Headers", "Grpc-Status,Grpc-Message,Grpc-Status-Details-Bin,Connect-Protocol-Version")
}
