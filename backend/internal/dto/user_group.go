package dto

type UserGroupResponse struct {
	ExternalID            string          `json:"external_id"`
	Name                  string          `json:"name"`
	Description           string          `json:"description"`
	CreatorUserExternalID string          `json:"creator_user_external_id"`
	MemberCount           int             `json:"member_count"`
	Members               []*UserResponse `json:"members"`
}

type ListUserGroupsRequest struct {
	Keyword    *string    `json:"keyword"`
	Pagination Pagination `json:"pagination"`
}

type GetUserGroupRequest struct {
	ExternalID string `json:"external_id"`
}

type CreateUserGroupRequest struct {
	CreatorUserExternalID string   `json:"creator_user_external_id"`
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	MemberUserExternalIDs []string `json:"member_user_external_ids"`
}

type UpdateUserGroupRequest struct {
	ExternalID            string   `json:"external_id"`
	Name                  *string  `json:"name"`
	Description           *string  `json:"description"`
	SyncMembers           bool     `json:"sync_members"`
	MemberUserExternalIDs []string `json:"member_user_external_ids"`
}

type DeleteUserGroupRequest struct {
	ExternalID string `json:"external_id"`
}
