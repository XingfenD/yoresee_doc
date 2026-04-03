package dto

type GetDocumentSettingsRequest struct {
	ExternalID string `json:"external_id"`
}

type DocumentSettingsResponse struct {
	IsPublic bool `json:"is_public"`
}

type UpdateDocumentSettingsRequest struct {
	ExternalID string `json:"external_id"`
	IsPublic   bool   `json:"is_public"`
}
