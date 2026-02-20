package api

var HealthHandlerImpl = &HealthHandler{}
var TestProtectedHandlerImpl = &TestProtectedHandler{}
var TestPostHandlerImpl = &TestPostHandler{}

type HealthHandler struct{}
type TestProtectedHandler struct{}
type TestPostHandler struct{}

type HealthResponse struct {
	BaseResponse
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	Version   string `json:"version"`
}

type HealthRequest struct {
	BaseRequest
}

type TestProtectedResponse struct {
	BaseResponse
	Message string `json:"message"`
}

type TestProtectedRequest struct {
	BaseRequest
}

type TestPostResponse struct {
	BaseResponse
	Message string `json:"message"`
}

type TestPostRequest struct {
	BaseRequest
	Message string `json:"message"`
}
