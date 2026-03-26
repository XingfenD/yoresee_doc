package membership_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository/membership_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/sirupsen/logrus"
)

type MembershipService struct {
	repo    *membership_repo.MembershipRepository
	useRepo *user_repo.UserRepository
}

func NewMembershipService() *MembershipService {
	return &MembershipService{
		repo:    membership_repo.MembershipRepo,
		useRepo: user_repo.UserRepo,
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
	userID, err := s.useRepo.GetIDByExternalID(membership.UserExternalID).Exec()
	if err != nil {
		logrus.Errorf("[Service layer: MembershipService] GetIDByExternalID failed, user_external_id=%s, err=%+v", membership.UserExternalID, err)
		return status.StatusUserNotFound
	}

	membershipID, err := s.GetMembershipIDByExternalID(&membership.MembershipBaseRequest)
	if err != nil {
		logrus.Errorf("[Service layer: MembershipService] GetMembershipIDByExternalID failed, membership_external_id=%s, type=%d, err=%+v", membership.MembershipExternalID, membership.Type, err)
		return status.GenErrWithCustomMsg(err, "membership target not found")
	}
	model := &model.MembershipRelation{
		UserID:       userID,
		MembershipID: membershipID,
		Type:         membership.Type,
	}

	if err := s.repo.CreateMembership(model).Exec(); err != nil {
		logrus.Errorf("[Service layer: MembershipService] CreateMembership failed, user_id=%d, membership_id=%d, type=%d, err=%+v", userID, membershipID, membership.Type, err)
		return status.GenErrWithCustomMsg(status.StatusWriteDBError, "create membership relation failed")
	}
	return nil
}

func (s *MembershipService) GetUserGroupMeta(externalID string) (*model.UserGroupMeta, error) {
	meta, err := s.repo.GetUserGroupByExternalID(externalID).Exec()
	if err != nil {
		logrus.Errorf("[Service layer: MembershipService] GetUserGroupMeta failed, external_id=%s, err=%+v", externalID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusMembershipMetaNotFound, "user group not found")
	}
	return meta, nil
}

func (s *MembershipService) GetOrgNodeMeta(externalID string) (*model.OrgNodeMeta, error) {
	meta, err := s.repo.GetOrgNodeByExternalID(externalID).Exec()
	if err != nil {
		logrus.Errorf("[Service layer: MembershipService] GetOrgNodeMeta failed, external_id=%s, err=%+v", externalID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusMembershipMetaNotFound, "org node not found")
	}
	return meta, nil
}

// func (s *MembershipService) GetMembershipRelation(req *dto.MembershipBaseRequest) (*dto.MembershipRelationResponse, error) {
// 	var membershipMetaResponse *dto.MembershipMetaResponse
// 	var membershipID int64
// 	switch req.Type {
// 	case model.MembershipType_UserGroup:
// 		userGroupMeta, err := s.GetUserGroupMeta(req.MembershipExternalID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		membershipMetaResponse = dto.NewMembershipMetaResponseFromUserGroupMetaModel(userGroupMeta)
// 		membershipID = userGroupMeta.ID
// 	case model.MembershipType_OrgNode:
// 		orgNodeMeta, err := s.GetOrgNodeMeta(req.MembershipExternalID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		membershipMetaResponse = dto.NewMembershipMetaResponseFromOrgNodeMetaModel(orgNodeMeta)
// 		membershipID = orgNodeMeta.ID
// 	default:
// 		return nil, status.StatusInvalidMembershipType
// 	}
// 	model := &model.MembershipRelation{
// 		Type:         req.Type,
// 		MembershipID: membershipID,
// 	}
// 	memberships, err := s.repo.ListMembership(model).Exec()
// 	if err != nil {
// 		return nil, err
// 	}
// }

var MembershipSvc = NewMembershipService()
