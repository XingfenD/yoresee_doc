package membership_service

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/bytedance/gg/gslice"
	"gorm.io/gorm"
)

func (s *MembershipService) ListOrgNodes(req *dto.ListOrgNodesRequest) ([]*dto.OrgNodeResponse, int64, error) {
	if req == nil {
		return nil, 0, status.StatusParamError
	}
	req.Pagination = normalizePagination(req.Pagination)

	var parentID int64 = 0
	if req.ParentExternalID != "" {
		parent, err := s.repo.GetOrgNodeByExternalID(req.ParentExternalID).Exec()
		if err != nil {
			return nil, 0, status.StatusMembershipMetaNotFound
		}
		parentID = parent.ID
	}

	nodes, total, err := s.repo.QueryOrgNode().
		WithParentID(parentID).
		WithKeyword(req.Keyword).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize).
		ExecWithTotal()
	if err != nil {
		return nil, 0, status.StatusReadDBError
	}

	responses, err := s.buildOrgNodeResponses(nodes, req.IncludeChildren)
	if err != nil {
		return nil, 0, err
	}
	return responses, total, nil
}

func (s *MembershipService) GetOrgNode(req *dto.GetOrgNodeRequest) (*dto.OrgNodeResponse, error) {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" {
		return nil, status.StatusParamError
	}
	node, err := s.repo.GetOrgNodeByExternalID(req.ExternalID).Exec()
	if err != nil {
		return nil, status.StatusMembershipMetaNotFound
	}

	if req.IncludeChildren {
		return s.buildOrgNodeWithChildren(node)
	}

	responses, err := s.buildOrgNodeResponses([]model.OrgNodeMeta{*node}, false)
	if err != nil {
		return nil, err
	}
	if len(responses) == 0 {
		return nil, status.StatusMembershipMetaNotFound
	}
	return responses[0], nil
}

func (s *MembershipService) CreateOrgNode(req *dto.CreateOrgNodeRequest) (string, error) {
	if req == nil || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.CreatorUserExternalID) == "" {
		return "", status.StatusParamError
	}

	var createdExternalID string
	err := utils.WithTransaction(func(tx *gorm.DB) error {
		creatorID, err := s.useRepo.GetIDByExternalID(req.CreatorUserExternalID).WithTx(tx).Exec()
		if err != nil {
			return status.StatusUserNotFound
		}

		var parentID int64 = 0
		var parentPath string = ""
		if req.ParentExternalID != "" {
			parent, err := s.repo.GetOrgNodeByExternalID(req.ParentExternalID).WithTx(tx).Exec()
			if err != nil {
				return status.StatusMembershipMetaNotFound
			}
			parentID = parent.ID
			parentPath = parent.Path
		}

		node := &model.OrgNodeMeta{
			ExternalID:  utils.GenerateExternalID(utils.ExternalIDContextOrgNode),
			ParentID:    parentID,
			Name:        strings.TrimSpace(req.Name),
			Description: strings.TrimSpace(req.Description),
			CreatorID:   creatorID,
		}

		if err := s.repo.CreateOrgNode(node).WithTx(tx).Exec(); err != nil {
			return status.StatusWriteDBError
		}

		node.Path = buildOrgNodePath(parentPath, node.ID)
		if err := s.repo.UpdateOrgNode(node).WithTx(tx).Exec(); err != nil {
			return status.StatusWriteDBError
		}

		userIDs, err := s.resolveUserExternalIDsToIDs(req.MemberUserExternalIDs, tx)
		if err != nil {
			return err
		}
		if err := s.syncOrgNodeMembers(tx, node.ID, userIDs); err != nil {
			return err
		}

		createdExternalID = node.ExternalID
		return nil
	})
	if err != nil {
		return "", err
	}
	return createdExternalID, nil
}

func (s *MembershipService) UpdateOrgNode(req *dto.UpdateOrgNodeRequest) error {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" {
		return status.StatusParamError
	}
	if req.Name == nil && req.Description == nil && !req.SyncMembers {
		return status.StatusParamError
	}

	return utils.WithTransaction(func(tx *gorm.DB) error {
		node, err := s.repo.GetOrgNodeByExternalID(req.ExternalID).WithTx(tx).Exec()
		if err != nil {
			return status.StatusMembershipMetaNotFound
		}

		shouldSave := false
		if req.Name != nil {
			newName := strings.TrimSpace(*req.Name)
			if newName == "" {
				return status.StatusParamError
			}
			node.Name = newName
			shouldSave = true
		}
		if req.Description != nil {
			node.Description = strings.TrimSpace(*req.Description)
			shouldSave = true
		}
		if shouldSave {
			if err := s.repo.UpdateOrgNode(node).WithTx(tx).Exec(); err != nil {
				return status.StatusWriteDBError
			}
		}

		if req.SyncMembers {
			userIDs, err := s.resolveUserExternalIDsToIDs(req.MemberUserExternalIDs, tx)
			if err != nil {
				return err
			}
			if err := s.syncOrgNodeMembers(tx, node.ID, userIDs); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *MembershipService) DeleteOrgNode(req *dto.DeleteOrgNodeRequest) error {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" {
		return status.StatusParamError
	}

	return utils.WithTransaction(func(tx *gorm.DB) error {
		node, err := s.repo.GetOrgNodeByExternalID(req.ExternalID).WithTx(tx).Exec()
		if err != nil {
			return status.StatusMembershipMetaNotFound
		}

		children, err := s.repo.QueryOrgNode().WithParentID(node.ID).Exec()
		if err != nil {
			return status.StatusReadDBError
		}
		if len(children) > 0 {
			return status.StatusOrgNodeHasChildren
		}

		if err := tx.Where("type = ? AND membership_id = ?", model.MembershipType_OrgNode, node.ID).
			Delete(&model.MembershipRelation{}).Error; err != nil {
			return status.StatusWriteDBError
		}
		if err := s.repo.DeleteOrgNode(node).WithTx(tx).Exec(); err != nil {
			return status.StatusWriteDBError
		}
		return nil
	})
}

func (s *MembershipService) MoveOrgNode(req *dto.MoveOrgNodeRequest) error {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" {
		return status.StatusParamError
	}

	return utils.WithTransaction(func(tx *gorm.DB) error {
		node, err := s.repo.GetOrgNodeByExternalID(req.ExternalID).WithTx(tx).Exec()
		if err != nil {
			return status.StatusMembershipMetaNotFound
		}

		var newParentID int64 = 0
		var newParentPath string = ""
		if req.NewParentExternalID != "" {
			newParent, err := s.repo.GetOrgNodeByExternalID(req.NewParentExternalID).WithTx(tx).Exec()
			if err != nil {
				return status.StatusMembershipMetaNotFound
			}
			newParentID = newParent.ID
			newParentPath = newParent.Path

			if newParentID == node.ID {
				return status.StatusParamError
			}

			descendants, err := s.repo.QueryOrgNodeByPathPrefix(node.Path).WithTx(tx).Exec()
			if err != nil {
				return status.StatusReadDBError
			}
			for _, descendant := range descendants {
				if descendant.ID == newParentID {
					return status.StatusOrgNodeCannotMoveToDescendant
				}
			}
		}

		node.ParentID = newParentID
		oldPath := node.Path
		node.Path = buildOrgNodePath(newParentPath, node.ID)

		if err := s.repo.UpdateOrgNode(node).WithTx(tx).Exec(); err != nil {
			return status.StatusWriteDBError
		}

		descendants, err := s.repo.QueryOrgNodeByPathPrefix(oldPath).WithTx(tx).Exec()
		if err != nil {
			return status.StatusReadDBError
		}
		for _, descendant := range descendants {
			if descendant.ID == node.ID {
				continue
			}
			descendant.Path = strings.Replace(descendant.Path, oldPath, node.Path, 1)
			if err := s.repo.UpdateOrgNode(&descendant).WithTx(tx).Exec(); err != nil {
				return status.StatusWriteDBError
			}
		}

		return nil
	})
}

func (s *MembershipService) ListOrgNodeMembers(req *dto.ListOrgNodeMembersRequest) ([]*dto.UserResponse, int64, error) {
	if req == nil || strings.TrimSpace(req.ExternalID) == "" {
		return nil, 0, status.StatusParamError
	}
	req.Pagination = normalizePagination(req.Pagination)

	node, err := s.repo.GetOrgNodeByExternalID(req.ExternalID).Exec()
	if err != nil {
		return nil, 0, status.StatusMembershipMetaNotFound
	}

	relations, err := s.repo.ListMembership(&model.MembershipRelation{
		Type:         model.MembershipType_OrgNode,
		MembershipID: node.ID,
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

func buildOrgNodePath(parentPath string, nodeID int64) string {
	label := "n" + utils.Int64ToString(nodeID)
	if strings.TrimSpace(parentPath) == "" {
		return label
	}
	return parentPath + "." + label
}

func (s *MembershipService) buildOrgNodeResponses(nodes []model.OrgNodeMeta, includeChildren bool) ([]*dto.OrgNodeResponse, error) {
	if len(nodes) == 0 {
		return []*dto.OrgNodeResponse{}, nil
	}

	nodeIDs := make([]int64, 0, len(nodes))
	parentIDs := make([]int64, 0, len(nodes))
	creatorIDs := make([]int64, 0, len(nodes))
	for _, n := range nodes {
		nodeIDs = append(nodeIDs, n.ID)
		parentIDs = append(parentIDs, n.ParentID)
		creatorIDs = append(creatorIDs, n.CreatorID)
	}
	parentIDs = gslice.Uniq(parentIDs)
	creatorIDs = gslice.Uniq(creatorIDs)

	relations, err := s.repo.ListMembershipByTypeAndIDs(model.MembershipType_OrgNode, nodeIDs).Exec()
	if err != nil {
		return nil, status.StatusReadDBError
	}

	parentMap, err := s.repo.MGetOrgNodeByID(parentIDs).Exec()
	if err != nil {
		return nil, status.StatusReadDBError
	}

	creatorMap, err := s.useRepo.MGetUserByID(creatorIDs).Exec()
	if err != nil {
		return nil, status.StatusReadDBError
	}

	nodeMemberUserIDs := make(map[int64][]int64, len(nodeIDs))
	nodeMemberSeen := make(map[int64]map[int64]struct{}, len(nodeIDs))
	for _, relation := range relations {
		if _, ok := nodeMemberSeen[relation.MembershipID]; !ok {
			nodeMemberSeen[relation.MembershipID] = map[int64]struct{}{}
		}
		if _, exists := nodeMemberSeen[relation.MembershipID][relation.UserID]; exists {
			continue
		}
		nodeMemberSeen[relation.MembershipID][relation.UserID] = struct{}{}

		nodeMemberUserIDs[relation.MembershipID] = append(nodeMemberUserIDs[relation.MembershipID], relation.UserID)
	}

	responses := make([]*dto.OrgNodeResponse, 0, len(nodes))
	for _, n := range nodes {
		resp := &dto.OrgNodeResponse{
			ExternalID:  n.ExternalID,
			Name:        n.Name,
			Path:        n.Path,
			Description: n.Description,
			MemberCount: len(nodeMemberUserIDs[n.ID]),
		}
		if parent := parentMap[n.ParentID]; parent != nil {
			resp.ParentExternalID = parent.ExternalID
		}
		if creator := creatorMap[n.CreatorID]; creator != nil {
			resp.CreatorUserExternalID = creator.ExternalID
		}
		responses = append(responses, resp)
	}

	if includeChildren {
		for i, resp := range responses {
			node := nodes[i]
			children, err := s.repo.QueryOrgNode().WithParentID(node.ID).Exec()
			if err != nil {
				return nil, status.StatusReadDBError
			}
			if len(children) > 0 {
				childResponses, err := s.buildOrgNodeResponses(children, true)
				if err != nil {
					return nil, err
				}
				resp.Children = childResponses
			}
		}
	}

	return responses, nil
}

func (s *MembershipService) buildOrgNodeWithChildren(node *model.OrgNodeMeta) (*dto.OrgNodeResponse, error) {
	responses, err := s.buildOrgNodeResponses([]model.OrgNodeMeta{*node}, true)
	if err != nil {
		return nil, err
	}
	if len(responses) == 0 {
		return nil, status.StatusMembershipMetaNotFound
	}
	return responses[0], nil
}

func (s *MembershipService) syncOrgNodeMembers(tx *gorm.DB, orgNodeID int64, targetUserIDs []int64) error {
	existingRelations, err := s.repo.ListMembership(&model.MembershipRelation{
		Type:         model.MembershipType_OrgNode,
		MembershipID: orgNodeID,
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
			Type:         model.MembershipType_OrgNode,
			UserID:       userID,
			MembershipID: orgNodeID,
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
		if err := tx.Where("type = ? AND membership_id = ? AND user_id IN ?", model.MembershipType_OrgNode, orgNodeID, toDeleteUserIDs).
			Delete(&model.MembershipRelation{}).Error; err != nil {
			return status.StatusWriteDBError
		}
	}
	return nil
}
