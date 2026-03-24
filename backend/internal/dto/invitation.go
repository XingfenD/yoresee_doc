package dto

import "time"

type ListInvitationsReq struct {
	CreatorID      *int64     `json:"creator_id"`
	MaxUsedCnt     *int64     `json:"max_used_cnt"`
	ExpiresAtStart *string    `json:"expires_at_start"`
	ExpiresAtEnd   *string    `json:"expires_at_end"`
	CreatedAtStart *string    `json:"created_at_start"`
	CreatedAtEnd   *string    `json:"created_at_end"`
	Disabled       *bool      `json:"disabled"`
	OnlyMine       *bool      `json:"only_mine"`
	Keyword        *string    `json:"keyword"`
	SortArgs       SortArgs   `json:"sort_args"`
	Pagination     Pagination `json:"pagination"`
}

type CreateInvitationRequest struct {
	CreatorExternalID string     `json:"creator_external_id"`
	MaxUsedCnt        *int64     `json:"max_used_cnt"`
	ExpiresAt         *time.Time `json:"expires_at"`
	Note              *string    `json:"note"`
}

type UpdateInvitationRequest struct {
	Code       string     `json:"code"`
	MaxUsedCnt *int64     `json:"max_used_cnt"`
	ExpiresAt  *time.Time `json:"expires_at"`
	Disabled   *bool      `json:"disabled"`
	Note       *string    `json:"note"`
}

type DeleteInvitationRequest struct {
	Code string `json:"code"`
}

type InvitationResponse struct {
	ID                  int64      `json:"id"`
	Code                string     `json:"code"`
	CreatedByExternalID string     `json:"created_by_external_id"`
	CreatedByName       string     `json:"created_by_name"`
	UsedCnt             int64      `json:"used_cnt"`
	MaxUsedCnt          *int64     `json:"max_used_cnt"`
	ExpiresAt           *time.Time `json:"expires_at"`
	CreatedAt           time.Time  `json:"created_at"`
	Disabled            bool       `json:"disabled"`
	Note                *string    `json:"note"`
}

type ListInvitationRecordsRequest struct {
	Code        *string    `json:"code"`
	Status      *string    `json:"status"`
	UsedAtStart *string    `json:"used_at_start"`
	UsedAtEnd   *string    `json:"used_at_end"`
	CreatorID   *int64     `json:"creator_id"`
	OnlyMine    *bool      `json:"only_mine"`
	Keyword     *string    `json:"keyword"`
	Pagination  Pagination `json:"pagination"`
}

type InvitationRecordResponse struct {
	ID               int64     `json:"id"`
	Code             string    `json:"code"`
	UsedBy           string    `json:"used_by"`
	UsedByExternalID string    `json:"used_by_external_id"`
	UsedAt           time.Time `json:"used_at"`
	Status           string    `json:"status"`
}
