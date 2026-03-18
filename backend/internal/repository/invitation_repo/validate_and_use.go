package invitation_repo

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type InvitationValidateAndUseOperation struct {
	repo *InvitationRepository
	code string
	tx   *gorm.DB
}

func (r *InvitationRepository) ValidateAndUse(code string) *InvitationValidateAndUseOperation {
	return &InvitationValidateAndUseOperation{
		repo: r,
		code: code,
	}
}

func (op *InvitationValidateAndUseOperation) WithTx(tx *gorm.DB) *InvitationValidateAndUseOperation {
	op.tx = tx
	return op
}

func (op *InvitationValidateAndUseOperation) Exec() (bool, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	now := time.Now()
	result := db.Model(&model.Invitation{}).
		Where("code = ? AND deleted_at IS NULL AND disabled = ?", op.code, false).
		Where("(expires_at IS NULL OR expires_at > ?)", now).
		Where("(max_used_cnt IS NULL OR used_cnt < max_used_cnt)").
		UpdateColumn("used_cnt", gorm.Expr("used_cnt + ?", 1))

	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
