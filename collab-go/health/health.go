package health

import (
	"context"
	"encoding/json"
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

func NewHealthChecker(systemService yoreseedocpb.SystemServiceClient) *HealthChecker {
	return &HealthChecker{
		systemService: systemService,
	}
}

func (h *HealthChecker) SetDraining(v bool) {
	h.draining.Store(v)
}

func (h *HealthChecker) Liveness(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *HealthChecker) Readiness(w http.ResponseWriter, r *http.Request) {
	if h.draining.Load() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status":  "not_ready",
			"backend": "unknown",
			"detail":  "server is draining",
		})
		return
	}

	if err := h.checkBackend(); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status":  "not_ready",
			"backend": "unhealthy",
			"detail":  err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"backend": "ok",
	})
}

func (h *HealthChecker) Check(w http.ResponseWriter, r *http.Request) {
	backendStatus := "ok"
	if err := h.checkBackend(); err != nil {
		backendStatus = "unhealthy"
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"backend": backendStatus,
	})
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

func writeJSON(w http.ResponseWriter, status int, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
