package dto

type UserGroupResponse struct {
	ExternalID            string `json:"external_id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	CreatorUserExternalID string `json:"creator_user_external_id"`
	MemberCount           int    `json:"member_count"`
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

type ListUsersRequest struct {
	Keyword    *string    `json:"keyword"`
	Pagination Pagination `json:"pagination"`
}

type UpdateUserRequest struct {
	ExternalID string  `json:"external_id"`
	Username   *string `json:"username"`
	Email      *string `json:"email"`
	Nickname   *string `json:"nickname"`
	Status     *int32  `json:"status"`
}

type ListUserGroupMembersRequest struct {
	ExternalID string     `json:"external_id"`
	Keyword    *string    `json:"keyword"`
	Pagination Pagination `json:"pagination"`
}
