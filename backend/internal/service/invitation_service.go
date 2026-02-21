package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type InvitationService struct {
	invitationRepo *repository.InvitationRepository
}

func NewInvitationService() *InvitationService {
	return &InvitationService{
		invitationRepo: repository.InvitationRepo,
	}
}

func (s *InvitationService) Generate(userID int64, maxUsedCnt *int64, expiresAt *time.Time) (*model.Invitation, error) {
	// Simple random code
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	code := hex.EncodeToString(bytes)

	invitation := &model.Invitation{
		Code:       code,
		CreatedBy:  userID,
		ExpiresAt:  expiresAt,
		MaxUsedCnt: maxUsedCnt,
	}

	if err := storage.DB.Create(invitation).Error; err != nil {
		return nil, status.StatusWriteDBError
	}
	return invitation, nil
}

func (s *InvitationService) ListByCreator(userID int64) ([]model.Invitation, error) {
	return s.invitationRepo.ListByCreatedBy(userID).Exec()
}

func (s *InvitationService) ValidateAndUse(code string) error {
	invitation, err := s.invitationRepo.GetByCode(code).Exec()
	if err != nil {
		return status.StatusInvitationInvalid
	}

	if invitation.DeletedAt != nil {
		return status.GenErrWithCustomMsg(status.StatusInvitationInvalid, fmt.Sprintf("invitation code: %s deleted", code))
	}
	if invitation.Disabled == true {
		return status.GenErrWithCustomMsg(status.StatusInvitationInvalid, fmt.Sprintf("invitation code: %s disabled", code))
	}
	if invitation.ExpiresAt != nil && invitation.ExpiresAt.Before(time.Now()) {
		return status.GenErrWithCustomMsg(status.StatusInvitationInvalid, fmt.Sprintf("invitation code: %s expired", code))
	}
	if invitation.MaxUsedCnt != nil && invitation.UsedCnt >= *invitation.MaxUsedCnt {
		return status.GenErrWithCustomMsg(status.StatusInvitationInvalid, fmt.Sprintf("invitation code: %s used up", code))
	}

	invitation.UsedCnt++
	return s.invitationRepo.Update(invitation).Exec()
}

// 事务版本：验证并使用邀请码
func (s *InvitationService) ValidateAndUseWithTx(tx *gorm.DB, code string) error {
	invitation, err := s.invitationRepo.GetByCode(code).WithTx(tx).Exec()
	if err != nil {
		return status.StatusInvitationInvalid
	}

	if invitation.DeletedAt != nil {
		return status.GenErrWithCustomMsg(status.StatusInvitationInvalid, fmt.Sprintf("invitation code: %s deleted", code))
	}
	if invitation.Disabled == true {
		return status.GenErrWithCustomMsg(status.StatusInvitationInvalid, fmt.Sprintf("invitation code: %s disabled", code))
	}
	if invitation.ExpiresAt != nil && invitation.ExpiresAt.Before(time.Now()) {
		return status.GenErrWithCustomMsg(status.StatusInvitationInvalid, fmt.Sprintf("invitation code: %s expired", code))
	}
	if invitation.MaxUsedCnt != nil && invitation.UsedCnt >= *invitation.MaxUsedCnt {
		return status.GenErrWithCustomMsg(status.StatusInvitationInvalid, fmt.Sprintf("invitation code: %s used up", code))
	}

	invitation.UsedCnt++
	return s.invitationRepo.Update(invitation).WithTx(tx).Exec()
}

func (s *InvitationService) UpdateCode(code string, newExpireTime *time.Time, newMaxUsedCnt *int64, isDisabled *bool) error {
	invitation, err := s.invitationRepo.GetByCode(code).Exec()
	if err != nil {
		return err
	}
	if newExpireTime != nil {
		invitation.ExpiresAt = newExpireTime
	}
	if newMaxUsedCnt != nil {
		invitation.MaxUsedCnt = newMaxUsedCnt
	}
	if isDisabled != nil {
		invitation.Disabled = *isDisabled
	}
	if err := s.invitationRepo.Update(invitation).Exec(); err != nil {
		return status.StatusWriteDBError
	}
	return nil
}

var InvitationSvc = &InvitationService{}
