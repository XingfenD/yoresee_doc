package dto

import "time"

type DocumentVersionResponse struct {
	Version       int       `json:"version"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	ChangeSummary string    `json:"change_summary"`
	CreatedAt     time.Time `json:"created_at"`
}
