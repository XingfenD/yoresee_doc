package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc/metadata"

	"github.com/XingfenD/yoresee_doc/collab-go/auth"
	yoreseedocpb "github.com/XingfenD/yoresee_doc/collab-go/pkg/gen/yoresee_doc/v1"
	"github.com/XingfenD/yoresee_doc/collab-go/proxy"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
	authenticator   *auth.Authenticator
	proxy           *proxy.Proxy
	documentService yoreseedocpb.DocumentServiceClient
	internalRPCKey  string
}

func NewWSHandler(authenticator *auth.Authenticator, proxy *proxy.Proxy, documentService yoreseedocpb.DocumentServiceClient, internalRPCKey string) *WSHandler {
	return &WSHandler{
		authenticator:   authenticator,
		proxy:           proxy,
		documentService: documentService,
		internalRPCKey:  internalRPCKey,
	}
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	docID := strings.TrimPrefix(r.URL.Path, "/ws/doc/")
	if docID == "" || strings.Contains(docID, "/") {
		http.Error(w, "invalid doc id", http.StatusBadRequest)
		return
	}

	if err := authenticateRequest(r, h.authenticator); err != nil {
		log.Printf("collab-gateway unauthorized path=%s remote=%s err=%v", r.URL.Path, r.RemoteAddr, err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if code, err := checkDocumentExists(r.Context(), h.documentService, h.internalRPCKey, docID); err != nil {
		log.Printf("collab-gateway doc check failed docID=%s err=%v", docID, err)
		http.Error(w, err.Error(), code)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("collab-gateway upgrade failed path=%s remote=%s err=%v", r.URL.Path, r.RemoteAddr, err)
		return
	}
	defer conn.Close()

	coreConn, err := h.proxy.DialCore(docID)
	if err != nil {
		log.Printf("collab-gateway dial core failed docID=%s err=%v", docID, err)
		return
	}
	defer coreConn.Close()

	errCh := make(chan error, 2)
	go h.proxy.ProxyWS(conn, coreConn, errCh)
	go h.proxy.ProxyWS(coreConn, conn, errCh)

	<-errCh
}

func authenticateRequest(r *http.Request, authenticator *auth.Authenticator) error {
	token := r.URL.Query().Get("token")
	return authenticator.ValidateToken(token)
}

func checkDocumentExists(ctx context.Context, documentService yoreseedocpb.DocumentServiceClient, internalRPCKey, docID string) (int, error) {
	if documentService == nil {
		return 0, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	md := metadata.New(map[string]string{})
	if internalRPCKey != "" {
		md.Set("x-internal-key", internalRPCKey)
	}
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := documentService.GetDocumentSettings(ctx, &yoreseedocpb.GetDocumentSettingsRequest{
		DocumentExternalId: docID,
	})
	if err != nil {
		return http.StatusServiceUnavailable, fmt.Errorf("service unavailable")
	}
	if resp.GetBase().GetCode() != 0 {
		log.Printf("collab-gateway doc not found docID=%s code=%d msg=%s", docID, resp.GetBase().GetCode(), resp.GetBase().GetMessage())
		return http.StatusNotFound, fmt.Errorf("document not found")
	}
	return 0, nil
}
