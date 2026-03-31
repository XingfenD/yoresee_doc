package dto

import "time"

type AttachmentResponse struct {
	ExternalID         string    `json:"external_id"`
	DocumentExternalID string    `json:"document_external_id"`
	Name               string    `json:"name"`
	Size               int64     `json:"size"`
	MimeType           string    `json:"mime_type"`
	Path               string    `json:"path"`
	URL                string    `json:"url"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
