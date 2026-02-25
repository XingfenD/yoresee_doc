package service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

type MembershipService struct {
	repo        *repository.MembershipRepository
	userService *UserService
}

func NewMembershipService() *MembershipService {
	return &MembershipService{
		repo: repository.MembershipRepo,
	}
}

func (s *MembershipService) GetMembershipIDByExternalID(req *dto.MembershipBaseRequest) (int64, error) {
	var membershipID int64
	if req.Type == model.MembershipType_OrgNode {
		orgNodeID, err := s.repo.GetOrgNodeIDByExternalID(req.MembershipExternalID).Exec()
		if err != nil {
			return 0, status.StatusMembershipMetaNotFound
		}
		membershipID = orgNodeID
	} else if req.Type == model.MembershipType_UserGroup {
		userGroupID, err := s.repo.GetUserGroupIDByExternalID(req.MembershipExternalID).Exec()
		if err != nil {
			return 0, status.StatusMembershipMetaNotFound
		}
		membershipID = userGroupID
	} else {
		return 0, status.StatusInvalidMembershipType
	}
	return membershipID, nil
}

func (s *MembershipService) CreateMembershipRelation(membership *dto.CreateMembershipRelationRequest) error {
	userID, err := s.userService.GetIDByExternalID(membership.UserExternalID)
	if err != nil {
		return status.StatusUserNotFound
	}

	membershipID, err := s.GetMembershipIDByExternalID(&membership.MembershipBaseRequest)
	if err != nil {
		return err
	}
	model := &model.MembershipRelation{
		UserID:       userID,
		MembershipID: membershipID,
		Type:         membership.Type,
	}

	return s.repo.CreateMembership(model).Exec()
}

func (s *MembershipService) GetUserGroupMeta(externalID string) (*model.UserGroupMeta, error) {
	return s.repo.GetUserGroupByExternalID(externalID).Exec()
}

func (s *MembershipService) GetOrgNodeMeta(externalID string) (*model.OrgNodeMeta, error) {
	return s.repo.GetOrgNodeByExternalID(externalID).Exec()
}

func (s *MembershipService) GetMembershipRelation(req *dto.MembershipBaseRequest) (*dto.MembershipRelationResponse, error) {
	var membershipMetaResponse *dto.MembershipMetaResponse
	var membershipID int64
	switch req.Type {
	case model.MembershipType_UserGroup:
		userGroupMeta, err := s.GetUserGroupMeta(req.MembershipExternalID)
		if err != nil {
			return nil, err
		}
		membershipMetaResponse = dto.NewMembershipMetaResponseFromUserGroupMetaModel(userGroupMeta)
		membershipID = userGroupMeta.ID
	case model.MembershipType_OrgNode:
		orgNodeMeta, err := s.GetOrgNodeMeta(req.MembershipExternalID)
		if err != nil {
			return nil, err
		}
		membershipMetaResponse = dto.NewMembershipMetaResponseFromOrgNodeMetaModel(orgNodeMeta)
		membershipID = orgNodeMeta.ID
	default:
		return nil, status.StatusInvalidMembershipType
	}
	model := &model.MembershipRelation{
		Type:         req.Type,
		MembershipID: membershipID,
	}
	memberships, err := s.repo.ListMembership(model).Exec()
	if err != nil {
		return nil, err
	}
}

var MembershipSvc = NewMembershipService()
