package dto

type CreateCommentRequest struct {
	DocumentExternalID       string   `json:"document_external_id"`
	Content                  string   `json:"content"`
	ParentExternalID         *string  `json:"parent_external_id,omitempty"`
	AnchorID                 *string  `json:"anchor_id,omitempty"`
	Quote                    *string  `json:"quote,omitempty"`
	CreatorExternalID        string   `json:"creator_external_id"`
	MentionedUserExternalIDs []string `json:"mentioned_user_external_ids,omitempty"`
}

type ListCommentsRequest struct {
	DocumentExternalID string `json:"document_external_id"`
	Page               int    `json:"page"`
	PageSize           int    `json:"page_size"`
}

type DeleteCommentRequest struct {
	ExternalID         string `json:"external_id"`
	OperatorExternalID string `json:"operator_external_id"`
}

type UpdateCommentRequest struct {
	ExternalID         string `json:"external_id"`
	Content            string `json:"content"`
	OperatorExternalID string `json:"operator_external_id"`
}
