package dto

type CreateNotificationRequest struct {
	ReceiverExternalIDs []string `json:"receiver_external_ids"`
	Type                string   `json:"type"`
	Title               string   `json:"title"`
	Content             string   `json:"content"`
	PayloadJSON         string   `json:"payload_json"`
}

type ListNotificationsRequest struct {
	UserExternalID string     `json:"user_external_id"`
	Status         *string    `json:"status,omitempty"`
	Pagination     Pagination `json:"pagination"`
}

type MarkNotificationsReadRequest struct {
	UserExternalID string   `json:"user_external_id"`
	ExternalIDs    []string `json:"external_ids"`
}
