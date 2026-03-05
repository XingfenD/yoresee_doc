package dto

type ListInvitationsReq struct {
	CreatorID      *int64     `json:"creator_id"`
	MaxUsedCnt     *int64     `json:"max_used_cnt"`
	ExpiresAtStart *string    `json:"expires_at_start"`
	ExpiresAtEnd   *string    `json:"expires_at_end"`
	CreatedAtStart *string    `json:"created_at_start"`
	CreatedAtEnd   *string    `json:"created_at_end"`
	Disabled       *bool      `json:"disabled"`
	SortArgs       SortArgs   `json:"sort_args"`
	Pagination     Pagination `json:"pagination"`
}
