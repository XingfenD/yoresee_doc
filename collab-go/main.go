package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/XingfenD/yoresee_doc/collab-go/auth"
	"github.com/XingfenD/yoresee_doc/collab-go/config"
	"github.com/XingfenD/yoresee_doc/collab-go/handler"
	"github.com/XingfenD/yoresee_doc/collab-go/health"
	yoreseedocpb "github.com/XingfenD/yoresee_doc/collab-go/pkg/gen/yoresee_doc/v1"
	"github.com/XingfenD/yoresee_doc/collab-go/proxy"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var systemService yoreseedocpb.SystemServiceClient
	var documentService yoreseedocpb.DocumentServiceClient

	grpcClient, err := grpc.NewClient(cfg.BackendGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Warning: failed to initialize gRPC client: %v. Health check will show backend as unknown.", err)
	} else {
		systemService = yoreseedocpb.NewSystemServiceClient(grpcClient)
		documentService = yoreseedocpb.NewDocumentServiceClient(grpcClient)
		defer grpcClient.Close()
	}

	authenticator := auth.NewAuthenticator(cfg.JWTSecret)
	proxyHandler := proxy.NewProxy(cfg.CollabCoreURL)
	healthChecker := health.NewHealthChecker(systemService)
	wsHandler := handler.NewWSHandler(authenticator, proxyHandler, documentService, cfg.InternalRPCKey)

	mux := http.NewServeMux()
	healthChecker.RegisterProbeRoutes(mux)
	mux.Handle("/ws/doc/", wsHandler)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("collab-gateway listening on %s", cfg.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("collab-gateway exited unexpectedly: %v", err)
		}
	}()

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-signalCtx.Done()

	log.Printf("shutdown signal received, start draining")
	healthChecker.SetDraining(true)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("collab-gateway graceful shutdown failed: %v", err)
	}
}
