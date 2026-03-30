package connectserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

type probeResponse struct {
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func registerProbeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", handleHealthProbe)
	mux.HandleFunc("/readyz", handleReadinessProbe)
	mux.HandleFunc("/livez", handleLivenessProbe)
}

func handleHealthProbe(w http.ResponseWriter, _ *http.Request) {
	writeProbeResponse(w, http.StatusOK, "ok", "")
}

func handleLivenessProbe(w http.ResponseWriter, _ *http.Request) {
	writeProbeResponse(w, http.StatusOK, "ok", "")
}

func handleReadinessProbe(w http.ResponseWriter, _ *http.Request) {
	if err := readinessErr(); err != nil {
		writeProbeResponse(w, http.StatusServiceUnavailable, "not_ready", err.Error())
		return
	}
	writeProbeResponse(w, http.StatusOK, "ok", "")
}

func readinessErr() error {
	if draining.Load() {
		return fmt.Errorf("server is draining")
	}
	if storage.DB == nil {
		return fmt.Errorf("database is not initialized")
	}
	if storage.KVS == nil {
		return fmt.Errorf("redis is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	sqlDB, err := storage.DB.DB()
	if err != nil {
		return fmt.Errorf("database handle unavailable: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	if err := storage.KVS.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}

func writeProbeResponse(w http.ResponseWriter, statusCode int, status, detail string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(probeResponse{
		Status: status,
		Detail: detail,
	})
}
