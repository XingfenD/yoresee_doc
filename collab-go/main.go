package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/XingfenD/yoresee_doc/collab-go/auth"
	"github.com/XingfenD/yoresee_doc/collab-go/health"
	yoreseedocpb "github.com/XingfenD/yoresee_doc/collab-go/pkg/gen/yoresee_doc/v1"
	"github.com/XingfenD/yoresee_doc/collab-go/proxy"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":1234"
	}
	secret := os.Getenv("JWT_SECRET")
	coreURL := os.Getenv("COLLAB_CORE_URL")
	if coreURL == "" {
		coreURL = "ws://collab-core:1234"
	}
	backendAddr := os.Getenv("BACKEND_GRPC_ADDR")
	if backendAddr == "" {
		backendAddr = "backend:9090"
	}

	// var grpcConn *grpc.ClientConn
	var systemService yoreseedocpb.SystemServiceClient
	var err error

	// 使用推荐的 NewClient 方法替代 deprecated 的 Dial
	grpcClient, err := grpc.NewClient(backendAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Warning: failed to initialize gRPC client: %v. Health check will show backend as unknown.", err)
	} else {
		systemService = yoreseedocpb.NewSystemServiceClient(grpcClient)
		defer grpcClient.Close()
	}

	authenticator := auth.NewAuthenticator(secret)
	proxyHandler := proxy.NewProxy(coreURL)
	healthChecker := health.NewHealthChecker(systemService)

	http.HandleFunc("/health", healthChecker.Check)
	http.HandleFunc("/ws/doc/", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, authenticator, proxyHandler)
	})

	log.Printf("collab-gateway listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, authenticator *auth.Authenticator, proxyHandler *proxy.Proxy) {
	docID := strings.TrimPrefix(r.URL.Path, "/ws/doc/")
	if docID == "" || strings.Contains(docID, "/") {
		http.Error(w, "invalid doc id", http.StatusBadRequest)
		return
	}

	token := r.URL.Query().Get("token")
	if err := authenticator.ValidateToken(token); err != nil {
		log.Printf("collab-gateway unauthorized path=%s remote=%s err=%v", r.URL.Path, r.RemoteAddr, err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("collab-gateway upgrade failed path=%s remote=%s err=%v", r.URL.Path, r.RemoteAddr, err)
		return
	}
	defer conn.Close()

	coreConn, err := proxyHandler.DialCore(docID)
	if err != nil {
		log.Printf("collab-gateway dial core failed docID=%s err=%v", docID, err)
		return
	}
	defer coreConn.Close()

	errCh := make(chan error, 2)
	go proxyHandler.ProxyWS(conn, coreConn, errCh)
	go proxyHandler.ProxyWS(coreConn, conn, errCh)

	<-errCh
}
