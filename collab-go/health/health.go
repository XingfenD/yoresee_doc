package health

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	yoreseedocpb "github.com/XingfenD/yoresee_doc/collab-go/pkg/gen/yoresee_doc/v1"
)

type HealthChecker struct {
	systemService yoreseedocpb.SystemServiceClient
	draining      atomic.Bool
}

type probeResponse struct {
	Status  string `json:"status"`
	Detail  string `json:"detail,omitempty"`
	Backend string `json:"backend,omitempty"`
}

func NewHealthChecker(systemService yoreseedocpb.SystemServiceClient) *HealthChecker {
	return &HealthChecker{
		systemService: systemService,
	}
}

func (h *HealthChecker) SetDraining(v bool) {
	h.draining.Store(v)
}

func (h *HealthChecker) RegisterProbeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/readyz", h.Readiness)
	mux.HandleFunc("/livez", h.Liveness)
}

func (h *HealthChecker) Liveness(w http.ResponseWriter, _ *http.Request) {
	writeProbeResponse(w, http.StatusOK, "ok", "", "")
}

func (h *HealthChecker) Readiness(w http.ResponseWriter, _ *http.Request) {
	if err := h.readinessErr(); err != nil {
		writeProbeResponse(w, http.StatusServiceUnavailable, "not_ready", err.Error(), probeBackendStatusFromErr(err))
		return
	}
	writeProbeResponse(w, http.StatusOK, "ok", "", "ok")
}

func (h *HealthChecker) Health(w http.ResponseWriter, _ *http.Request) {
	detail := ""
	backend := "ok"
	if err := h.readinessErr(); err != nil {
		detail = err.Error()
		backend = probeBackendStatusFromErr(err)
	}
	writeProbeResponse(w, http.StatusOK, "ok", detail, backend)
}

func (h *HealthChecker) readinessErr() error {
	if h.draining.Load() {
		return fmt.Errorf("server is draining")
	}
	return h.checkBackend()
}

func writeProbeResponse(w http.ResponseWriter, statusCode int, status, detail, backend string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(probeResponse{
		Status:  status,
		Detail:  detail,
		Backend: backend,
	})
}

func probeBackendStatusFromErr(err error) string {
	if err == nil {
		return "ok"
	}
	if err.Error() == "server is draining" {
		return "unknown"
	}
	return "unhealthy"
}

func (h *HealthChecker) checkBackend() error {
	if h.systemService == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := h.systemService.Health(ctx, &yoreseedocpb.HealthRequest{}); err != nil {
		log.Printf("backend health check failed: %v", err)
		return err
	}
	return nil
}
