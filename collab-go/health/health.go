package health

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	yoreseedocpb "github.com/XingfenD/yoresee_doc/collab-go/pkg/gen/yoresee_doc/v1"
)

type HealthChecker struct {
	systemService yoreseedocpb.SystemServiceClient
}

func NewHealthChecker(systemService yoreseedocpb.SystemServiceClient) *HealthChecker {
	return &HealthChecker{
		systemService: systemService,
	}
}

func (h *HealthChecker) Check(w http.ResponseWriter, r *http.Request) {
	selfStatus := "ok"

	backendStatus := "unknown"
	if h.systemService != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := h.systemService.Health(ctx, &yoreseedocpb.HealthRequest{})
		if err != nil {
			log.Printf("backend health check failed: %v", err)
			backendStatus = "unhealthy"
		} else {
			backendStatus = "ok"
		}
	}

	response := fmt.Sprintf(`{
	"status": "%s",
	"backend": "%s"
}`, selfStatus, backendStatus)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
