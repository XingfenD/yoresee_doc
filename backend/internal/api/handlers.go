package api

var HealthHandlerImpl = &HealthHandler{}
var TestProtectedHandlerImpl = &TestProtectedHandler{}
var TestPostHandlerImpl = &TestPostHandler{}
var AuthRegisterHandlerImpl = &AuthRegisterHandler{}
var AuthLoginHandlerImpl = &AuthLoginHandler{}
var SystemInfoHandlerImpl = &SystemInfoHandler{}
var GetDocumentContentHandlerImpl = &GetDocumentContentHandler{}
var ListDocumentsHandlerImpl = &ListDocumentsHandler{}
var ListKnowledgeBasesHandlerImpl = &ListKnowledgeBasesHandler{}
var GetKnowledgeBaseHandlerImpl = &GetKnowledgeBaseHandler{}

type HealthHandler struct{}
type TestProtectedHandler struct{}
type TestPostHandler struct{}
type AuthRegisterHandler struct{}
type AuthLoginHandler struct{}
type SystemInfoHandler struct{}
type GetDocumentContentHandler struct{}
type ListDocumentsHandler struct{}
type ListKnowledgeBasesHandler struct{}
type GetKnowledgeBaseHandler struct{}
