package membership_service

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/auth"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/bytedance/gg/gslice"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 200
)

func normalizePagination(p dto.Pagination) dto.Pagination {
	if p.Page <= 0 {
		p.Page = defaultPage
	}
	if p.PageSize <= 0 {
		p.PageSize = defaultPageSize
	}
	if p.PageSize > maxPageSize {
		p.PageSize = maxPageSize
	}
	return p
}

func (s *MembershipService) ListUsers(req *dto.ListUsersRequest) ([]*dto.UserResponse, int64, error) {
	if req == nil {
		return nil, 0, status.StatusParamError
	}
	req.Pagination = normalizePagination(req.Pagination)

	users, total, err := s.useRepo.QueryUsers().
		WithKeyword(req.Keyword).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize).
		ExecWithTotal()
	if err != nil {
		return nil, 0, status.StatusReadDBError
	}

	resp := make([]*dto.UserResponse, 0, len(users))
	for i := range users {
		user := users[i]
		resp = append(resp, dto.NewUserResponseFromModel(&user))
	}
	return resp, total, nil
}

func (s *MembershipService) UpdateUser(req *dto.UpdateUserRequest) error {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" {
		return status.StatusParamError
	}
	if req.Username == nil && req.Email == nil && req.Nickname == nil && req.Status == nil {
		return status.StatusParamError
	}

	user, err := s.useRepo.GetByExternalID(req.ExternalID).Exec()
	if err != nil {
		return status.StatusUserNotFound
	}
	oldStatus := user.Status

	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if username == "" {
			return status.StatusParamError
		}
		user.Username = username
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email == "" {
			return status.StatusParamError
		}
		user.Email = email
	}
	if req.Nickname != nil {
		user.Nickname = strings.TrimSpace(*req.Nickname)
	}
	if req.Status != nil {
		user.Status = int(*req.Status)
	}

	if err := s.useRepo.Update(user).Exec(); err != nil {
		return status.StatusWriteDBError
	}
	if req.Status != nil && oldStatus > 0 && user.Status <= 0 {
		if err := auth.BlacklistUserJWTs(user.ExternalID); err != nil {
			logrus.Errorf("[Service layer: MembershipService] blacklist user jwt tokens failed, user_external_id=%s, err=%+v", user.ExternalID, err)
			return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "ban user failed: blacklist jwt tokens failed")
		}
	}
	return nil
}

func (s *MembershipService) ListUserGroupMembers(req *dto.ListUserGroupMembersRequest) ([]*dto.UserResponse, int64, error) {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" {
		return nil, 0, status.StatusParamError
	}
	req.Pagination = normalizePagination(req.Pagination)

	group, err := s.repo.GetUserGroupByExternalID(req.ExternalID).Exec()
	if err != nil {
		return nil, 0, status.StatusMembershipMetaNotFound
	}

	relations, err := s.repo.ListMembership(&model.MembershipRelation{
		Type:         model.MembershipType_UserGroup,
		MembershipID: group.ID,
	}).Exec()
	if err != nil {
		return nil, 0, status.StatusReadDBError
	}
	if len(relations) == 0 {
		return []*dto.UserResponse{}, 0, nil
	}

	memberIDs := make([]int64, 0, len(relations))
	for _, relation := range relations {
		memberIDs = append(memberIDs, relation.UserID)
	}

	query := s.useRepo.QueryUsers().
		WithPagination(req.Pagination.Page, req.Pagination.PageSize)
	if req.Keyword != nil && strings.TrimSpace(*req.Keyword) != "" {
		query = query.WithKeyword(req.Keyword)
	}
	query = query.WithUserIDs(memberIDs)

	users, total, err := query.ExecWithTotal()
	if err != nil {
		return nil, 0, status.StatusReadDBError
	}

	resp := make([]*dto.UserResponse, 0, len(users))
	for i := range users {
		user := users[i]
		resp = append(resp, dto.NewUserResponseFromModel(&user))
	}

	return resp, total, nil
}

func (s *MembershipService) ListUserGroups(req *dto.ListUserGroupsRequest) ([]*dto.UserGroupResponse, int64, error) {
	if req == nil {
		return nil, 0, status.StatusParamError
	}
	req.Pagination = normalizePagination(req.Pagination)

	groups, total, err := s.repo.QueryUserGroup().
		WithKeyword(req.Keyword).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize).
		ExecWithTotal()
	if err != nil {
		return nil, 0, status.StatusReadDBError
	}

	responses, err := s.buildUserGroupResponses(groups, false)
	if err != nil {
		logrus.Errorf("[Service layer: MembershipService] buildUserGroupResponses failed, err=%+v", err)
		return nil, 0, status.GenErrWithCustomMsg(err, "build user group response failed")
	}
	return responses, total, nil
}

func (s *MembershipService) GetUserGroup(req *dto.GetUserGroupRequest) (*dto.UserGroupResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" {
		return nil, status.StatusParamError
	}
	group, err := s.repo.GetUserGroupByExternalID(req.ExternalID).Exec()
	if err != nil {
		return nil, status.StatusMembershipMetaNotFound
	}

	responses, err := s.buildUserGroupResponses([]model.UserGroupMeta{*group}, false)
	if err != nil {
		logrus.Errorf("[Service layer: MembershipService] buildUserGroupResponses failed, external_id=%s, err=%+v", req.ExternalID, err)
		return nil, status.GenErrWithCustomMsg(err, "build user group response failed")
	}
	if len(responses) == 0 {
		return nil, status.StatusMembershipMetaNotFound
	}
	return responses[0], nil
}

func (s *MembershipService) CreateUserGroup(req *dto.CreateUserGroupRequest) (string, error) {
	if req == nil || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.CreatorUserExternalID) == "" {
		return "", status.StatusParamError
	}

	var createdExternalID string
	err := utils.WithTransaction(func(tx *gorm.DB) error {
		creatorID, err := s.useRepo.GetIDByExternalID(req.CreatorUserExternalID).WithTx(tx).Exec()
		if err != nil {
			return status.StatusUserNotFound
		}

		group := &model.UserGroupMeta{
			ExternalID:  utils.GenerateExternalID(utils.ExternalIDContextUserGroup),
			Name:        strings.TrimSpace(req.Name),
			Description: strings.TrimSpace(req.Description),
			CreatorID:   creatorID,
		}
		if err := s.repo.CreateUserGroup(group).WithTx(tx).Exec(); err != nil {
			return status.StatusWriteDBError
		}

		userIDs, err := s.resolveUserExternalIDsToIDs(req.MemberUserExternalIDs, tx)
		if err != nil {
			logrus.Errorf("[Service layer: MembershipService] resolveUserExternalIDsToIDs failed, err=%+v", err)
			return status.GenErrWithCustomMsg(err, "resolve user ids failed")
		}
		if err := s.syncUserGroupMembers(tx, group.ID, userIDs); err != nil {
			logrus.Errorf("[Service layer: MembershipService] syncUserGroupMembers failed, group_id=%d, err=%+v", group.ID, err)
			return status.GenErrWithCustomMsg(err, "sync user group members failed")
		}

		createdExternalID = group.ExternalID
		return nil
	})
	if err != nil {
		logrus.Errorf("[Service layer: MembershipService] CreateUserGroup transaction failed, err=%+v", err)
		return "", status.GenErrWithCustomMsg(err, "create user group failed")
	}
	return createdExternalID, nil
}

func (s *MembershipService) UpdateUserGroup(req *dto.UpdateUserGroupRequest) error {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" {
		return status.StatusParamError
	}
	if req.Name == nil && req.Description == nil && !req.SyncMembers {
		return status.StatusParamError
	}

	return utils.WithTransaction(func(tx *gorm.DB) error {
		group, err := s.repo.GetUserGroupByExternalID(req.ExternalID).WithTx(tx).Exec()
		if err != nil {
			return status.StatusMembershipMetaNotFound
		}

		shouldSave := false
		if req.Name != nil {
			newName := strings.TrimSpace(*req.Name)
			if newName == "" {
				return status.StatusParamError
			}
			group.Name = newName
			shouldSave = true
		}
		if req.Description != nil {
			group.Description = strings.TrimSpace(*req.Description)
			shouldSave = true
		}
		if shouldSave {
			if err := s.repo.UpdateUserGroup(group).WithTx(tx).Exec(); err != nil {
				return status.StatusWriteDBError
			}
		}

		if req.SyncMembers {
			userIDs, err := s.resolveUserExternalIDsToIDs(req.MemberUserExternalIDs, tx)
			if err != nil {
				logrus.Errorf("[Service layer: MembershipService] resolveUserExternalIDsToIDs failed, external_id=%s, err=%+v", req.ExternalID, err)
				return status.GenErrWithCustomMsg(err, "resolve user ids failed")
			}
			if err := s.syncUserGroupMembers(tx, group.ID, userIDs); err != nil {
				logrus.Errorf("[Service layer: MembershipService] syncUserGroupMembers failed, group_id=%d, err=%+v", group.ID, err)
				return status.GenErrWithCustomMsg(err, "sync user group members failed")
			}
		}
		return nil
	})
}

func (s *MembershipService) DeleteUserGroup(req *dto.DeleteUserGroupRequest) error {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" {
		return status.StatusParamError
	}

	return utils.WithTransaction(func(tx *gorm.DB) error {
		group, err := s.repo.GetUserGroupByExternalID(req.ExternalID).WithTx(tx).Exec()
		if err != nil {
			return status.StatusMembershipMetaNotFound
		}

		if err := tx.Where("type = ? AND membership_id = ?", model.MembershipType_UserGroup, group.ID).
			Delete(&model.MembershipRelation{}).Error; err != nil {
			return status.StatusWriteDBError
		}
		if err := s.repo.DeleteUserGroup(group).WithTx(tx).Exec(); err != nil {
			return status.StatusWriteDBError
		}
		return nil
	})
}

func (s *MembershipService) buildUserGroupResponses(groups []model.UserGroupMeta, includeMembers bool) ([]*dto.UserGroupResponse, error) {
	if len(groups) == 0 {
		return []*dto.UserGroupResponse{}, nil
	}

	groupIDs := make([]int64, 0, len(groups))
	creatorIDs := make([]int64, 0, len(groups))
	for _, g := range groups {
		groupIDs = append(groupIDs, g.ID)
		creatorIDs = append(creatorIDs, g.CreatorID)
	}
	creatorIDs = gslice.Uniq(creatorIDs)

	relations, err := s.repo.ListMembershipByTypeAndIDs(model.MembershipType_UserGroup, groupIDs).Exec()
	if err != nil {
		return nil, status.StatusReadDBError
	}

	creatorMap, err := s.useRepo.MGetUserByID(creatorIDs).Exec()
	if err != nil {
		return nil, status.StatusReadDBError
	}

	groupMemberUserIDs := make(map[int64][]int64, len(groupIDs))
	groupMemberSeen := make(map[int64]map[int64]struct{}, len(groupIDs))
	for _, relation := range relations {
		if _, ok := groupMemberSeen[relation.MembershipID]; !ok {
			groupMemberSeen[relation.MembershipID] = map[int64]struct{}{}
		}
		if _, exists := groupMemberSeen[relation.MembershipID][relation.UserID]; exists {
			continue
		}
		groupMemberSeen[relation.MembershipID][relation.UserID] = struct{}{}

		groupMemberUserIDs[relation.MembershipID] = append(groupMemberUserIDs[relation.MembershipID], relation.UserID)
	}

	responses := make([]*dto.UserGroupResponse, 0, len(groups))
	for _, g := range groups {
		resp := &dto.UserGroupResponse{
			ExternalID:  g.ExternalID,
			Name:        g.Name,
			Description: g.Description,
			MemberCount: len(groupMemberUserIDs[g.ID]),
		}
		if creator := creatorMap[g.CreatorID]; creator != nil {
			resp.CreatorUserExternalID = creator.ExternalID
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func (s *MembershipService) resolveUserExternalIDsToIDs(externalIDs []string, tx *gorm.DB) ([]int64, error) {
	if len(externalIDs) == 0 {
		return []int64{}, nil
	}

	deduped := make([]string, 0, len(externalIDs))
	seen := map[string]struct{}{}
	for _, ext := range externalIDs {
		ext = strings.TrimSpace(ext)
		if ext == "" {
			return nil, status.StatusParamError
		}
		if _, ok := seen[ext]; ok {
			continue
		}
		seen[ext] = struct{}{}
		deduped = append(deduped, ext)
	}

	users, err := s.useRepo.ListByExternal(deduped).WithTx(tx).Exec()
	if err != nil {
		return nil, status.StatusReadDBError
	}
	if len(users) != len(deduped) {
		return nil, status.StatusUserNotFound
	}

	userIDs := make([]int64, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	return userIDs, nil
}

func (s *MembershipService) syncUserGroupMembers(tx *gorm.DB, groupID int64, targetUserIDs []int64) error {
	existingRelations, err := s.repo.ListMembership(&model.MembershipRelation{
		Type:         model.MembershipType_UserGroup,
		MembershipID: groupID,
	}).WithTx(tx).Exec()
	if err != nil {
		return status.StatusReadDBError
	}

	existingUserIDs := map[int64]struct{}{}
	for _, relation := range existingRelations {
		existingUserIDs[relation.UserID] = struct{}{}
	}
	targetUserIDSet := map[int64]struct{}{}
	for _, userID := range targetUserIDs {
		targetUserIDSet[userID] = struct{}{}
	}

	toCreate := make([]*model.MembershipRelation, 0)
	for userID := range targetUserIDSet {
		if _, ok := existingUserIDs[userID]; ok {
			continue
		}
		toCreate = append(toCreate, &model.MembershipRelation{
			Type:         model.MembershipType_UserGroup,
			UserID:       userID,
			MembershipID: groupID,
		})
	}

	toDeleteUserIDs := make([]int64, 0)
	for userID := range existingUserIDs {
		if _, ok := targetUserIDSet[userID]; ok {
			continue
		}
		toDeleteUserIDs = append(toDeleteUserIDs, userID)
	}

	if len(toCreate) > 0 {
		if err := s.repo.BatchCreateMembership(toCreate).WithTx(tx).Exec(); err != nil {
			return status.StatusWriteDBError
		}
	}
	if len(toDeleteUserIDs) > 0 {
		if err := tx.Where("type = ? AND membership_id = ? AND user_id IN ?", model.MembershipType_UserGroup, groupID, toDeleteUserIDs).
			Delete(&model.MembershipRelation{}).Error; err != nil {
			return status.StatusWriteDBError
		}
	}
	return nil
}
