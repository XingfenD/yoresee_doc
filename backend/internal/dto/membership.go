package dto

import "github.com/XingfenD/yoresee_doc/internal/model"

type MembershipRelationBase struct {
	Type         model.MembershipType
	UserID       int64
	MembershipID int64
}

type MembershipBase struct {
	Type                 model.MembershipType
	MembershipExternalID string
	MembershipName       string
}

type MembershipBaseRequest struct {
	Type                 model.MembershipType
	MembershipExternalID string
}

type MembershipMetaResponse struct {
	MembershipBase
	Description string
}

func NewMembershipMetaResponseFromUserGroupMetaModel(m *model.UserGroupMeta) *MembershipMetaResponse {
	if m == nil {
		return nil
	}

	return &MembershipMetaResponse{
		MembershipBase: MembershipBase{
			Type:                 model.MembershipType_UserGroup,
			MembershipName:       m.Name,
			MembershipExternalID: m.ExternalID,
		},
		Description: m.Description,
	}
}

func NewMembershipMetaResponseFromOrgNodeMetaModel(m *model.OrgNodeMeta) *MembershipMetaResponse {
	if m == nil {
		return nil
	}

	return &MembershipMetaResponse{
		MembershipBase: MembershipBase{
			Type:                 model.MembershipType_OrgNode,
			MembershipName:       m.Name,
			MembershipExternalID: m.ExternalID,
		},
		Description: m.Description,
	}
}

type CreateMembershipRelationRequest struct {
	MembershipBaseRequest
	UserExternalID string
}

type MembershipRelationResponse struct {
	MembershipBase
	UserList []UserBase
}

// type MembershipAuthorityResponse struct {
// 	MembershipBase
// 	CreatorExternalID    string
// 	CurrentUserAuthority []model.Permission
// }
