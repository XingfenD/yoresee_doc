package dto

type OrgNodeResponse struct {
	ExternalID            string             `json:"external_id"`
	ParentExternalID      string             `json:"parent_external_id"`
	Name                  string             `json:"name"`
	Path                  string             `json:"path"`
	Description           string             `json:"description"`
	CreatorUserExternalID string             `json:"creator_user_external_id"`
	MemberCount           int                `json:"member_count"`
	Children              []*OrgNodeResponse `json:"children,omitempty"`
}

type ListOrgNodesRequest struct {
	ParentExternalID string     `json:"parent_external_id"`
	Keyword          *string    `json:"keyword"`
	Pagination       Pagination `json:"pagination"`
	IncludeChildren  bool       `json:"include_children"`
}

type GetOrgNodeRequest struct {
	ExternalID      string `json:"external_id"`
	IncludeChildren bool   `json:"include_children"`
}

type CreateOrgNodeRequest struct {
	CreatorUserExternalID string   `json:"creator_user_external_id"`
	ParentExternalID      string   `json:"parent_external_id"`
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	MemberUserExternalIDs []string `json:"member_user_external_ids"`
}

type UpdateOrgNodeRequest struct {
	ExternalID            string   `json:"external_id"`
	Name                  *string  `json:"name"`
	Description           *string  `json:"description"`
	SyncMembers           bool     `json:"sync_members"`
	MemberUserExternalIDs []string `json:"member_user_external_ids"`
}

type DeleteOrgNodeRequest struct {
	ExternalID string `json:"external_id"`
}

type MoveOrgNodeRequest struct {
	ExternalID          string `json:"external_id"`
	NewParentExternalID string `json:"new_parent_external_id"`
}

type ListOrgNodeMembersRequest struct {
	ExternalID string     `json:"external_id"`
	Keyword    *string    `json:"keyword"`
	Pagination Pagination `json:"pagination"`
}
