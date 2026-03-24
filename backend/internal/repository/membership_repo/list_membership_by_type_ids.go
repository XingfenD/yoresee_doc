package membership_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type ListMembershipByTypeAndIDsOperation struct {
	repo           *MembershipRepository
	membershipType model.MembershipType
	membershipIDs  []int64
	tx             *gorm.DB
}

func (r *MembershipRepository) ListMembershipByTypeAndIDs(membershipType model.MembershipType, membershipIDs []int64) *ListMembershipByTypeAndIDsOperation {
	return &ListMembershipByTypeAndIDsOperation{
		repo:           r,
		membershipType: membershipType,
		membershipIDs:  membershipIDs,
	}
}

func (op *ListMembershipByTypeAndIDsOperation) WithTx(tx *gorm.DB) *ListMembershipByTypeAndIDsOperation {
	op.tx = tx
	return op
}

func (op *ListMembershipByTypeAndIDsOperation) Exec() ([]model.MembershipRelation, error) {
	if len(op.membershipIDs) == 0 {
		return []model.MembershipRelation{}, nil
	}
	if op.tx == nil {
		op.tx = storage.DB
	}

	var memberships []model.MembershipRelation
	err := op.tx.
		Where("type = ? AND membership_id IN ?", op.membershipType, op.membershipIDs).
		Find(&memberships).Error
	return memberships, err
}
