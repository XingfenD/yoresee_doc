package dto

type CreateCommentRequest struct {
	DocumentExternalID string
	Content            string
	ParentExternalID   *string
	AnchorID           *string
	Quote              *string
	CreatorExternalID  string
}

const (
	CommentScopeUnspecified int32 = 0
	CommentScopeAll         int32 = 1
	CommentScopeNormal      int32 = 2
	CommentScopeInline      int32 = 3
)

type ListCommentsRequest struct {
	DocumentExternalID string
	Page               int
	PageSize           int
	Scope              int32
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
