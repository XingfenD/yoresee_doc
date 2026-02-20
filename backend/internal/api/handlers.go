package api

var HealthHandlerImpl = &HealthHandler{}

type HealthHandler struct {
}

type HealthResponse struct {
	BaseResponse
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	Version   string `json:"version"`
}

type HealthRequest struct {
	BaseRequest
}
