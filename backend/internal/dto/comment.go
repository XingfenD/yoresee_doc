package dto

type CreateCommentRequest struct {
	DocumentExternalID        string
	Content                   string
	ParentExternalID          *string
	AnchorID                  *string
	Quote                     *string
	CreatorExternalID         string
	MentionedUserExternalIDs  []string
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

type UpdateCommentRequest struct {
	ExternalID         string
	Content            string
	OperatorExternalID string
}
