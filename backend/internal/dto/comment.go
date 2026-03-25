package dto

type CreateCommentRequest struct {
	DocumentExternalID string
	Content            string
	ParentExternalID   *string
	CreatorExternalID  string
}

type ListCommentsRequest struct {
	DocumentExternalID string
	Page               int
	PageSize           int
}

type DeleteCommentRequest struct {
	ExternalID         string
	OperatorExternalID string
}
