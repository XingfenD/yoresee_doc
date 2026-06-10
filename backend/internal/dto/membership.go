package dto

// MembershipType mirrors model.MembershipType to keep DTO layer independent.
type MembershipType int64

const (
	MembershipType_UserGroup MembershipType = 1
	MembershipType_OrgNode   MembershipType = 2
)

type MembershipRelationBase struct {
	Type         MembershipType `json:"type"`
	UserID       int64          `json:"user_id"`
	MembershipID int64          `json:"membership_id"`
}

type MembershipBase struct {
	Type                 MembershipType `json:"type"`
	MembershipExternalID string         `json:"membership_external_id"`
	MembershipName       string         `json:"membership_name"`
}

type MembershipBaseRequest struct {
	Type                 MembershipType `json:"type"`
	MembershipExternalID string         `json:"membership_external_id"`
}

type MembershipMetaResponse struct {
	MembershipBase
	Description string `json:"description"`
}

type CreateMembershipRelationRequest struct {
	MembershipBaseRequest
	UserExternalID string `json:"user_external_id"`
}

type MembershipRelationResponse struct {
	MembershipBase
	UserList []UserBase `json:"user_list"`
}

// type MembershipAuthorityResponse struct {
// 	MembershipBase
// 	CreatorExternalID    string
// 	CurrentUserAuthority []model.Permission
// }
