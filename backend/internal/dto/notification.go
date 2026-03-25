package dto

type CreateNotificationRequest struct {
	ReceiverExternalIDs []string
	Type                string
	Title               string
	Content             string
	PayloadJSON         string
}

type ListNotificationsRequest struct {
	UserExternalID string
	Status         *string
	Pagination     Pagination
}

type MarkNotificationsReadRequest struct {
	UserExternalID string
	IDs            []int64
}
