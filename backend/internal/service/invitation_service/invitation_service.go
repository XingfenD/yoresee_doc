package invitation_service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository/invitation_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type InvitationService struct {
	invitationRepo *invitation_repo.InvitationRepository
	userRepo       *user_repo.UserRepository
}

func NewInvitationService() *InvitationService {
	return &InvitationService{
		invitationRepo: invitation_repo.InvitationRepo,
		userRepo:       user_repo.UserRepo,
	}
}

func (s *InvitationService) Generate(userID int64, maxUsedCnt *int64, expiresAt *time.Time, note *string) (*model.Invitation, error) {
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
		Note:       note,
	}

	if err := storage.DB.Create(invitation).Error; err != nil {
		return nil, status.StatusWriteDBError
	}
	return invitation, nil
}

func (s *InvitationService) ListByCreator(userID int64) ([]model.Invitation, error) {
	return s.invitationRepo.List(&model.Invitation{CreatedBy: userID}).Exec()
}

func (s *InvitationService) ListInvitations(req *dto.ListInvitationsReq) ([]model.Invitation, int64, error) {

	return s.invitationRepo.List(&model.Invitation{}).
		WithCreatorID(req.CreatorID).
		WithMaxUsedCnt(req.MaxUsedCnt).
		WithExpiresAtRange(req.ExpiresAtStart, req.ExpiresAtEnd).
		WithCreatedAtRange(req.CreatedAtStart, req.CreatedAtEnd).
		WithDisabled(req.Disabled).
		WithSort(req.SortArgs.Field, req.SortArgs.Desc).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize).
		ExecWithTotal()
}

func (s *InvitationService) CreateInvitation(req *dto.CreateInvitationRequest) (*model.Invitation, error) {
	if req == nil || req.CreatorExternalID == "" {
		return nil, status.StatusParamError
	}
	// TODO: normal user can't create invitation codes
	userID, err := s.userRepo.GetIDByExternalID(req.CreatorExternalID).Exec()
	if err != nil {
		return nil, status.StatusUserNotFound
	}
	return s.Generate(userID, req.MaxUsedCnt, req.ExpiresAt, req.Note)
}

func (s *InvitationService) UpdateInvitation(req *dto.UpdateInvitationRequest) error {
	if req == nil || req.Code == "" {
		return status.StatusParamError
	}
	return s.UpdateCode(req.Code, req.ExpiresAt, req.MaxUsedCnt, req.Disabled, req.Note)
}

func (s *InvitationService) DeleteInvitation(req *dto.DeleteInvitationRequest) error {
	if req == nil || req.Code == "" {
		return status.StatusParamError
	}
	inv, err := s.invitationRepo.GetByCode(req.Code).Exec()
	if err != nil {
		return status.StatusInvitationInvalid
	}
	if err := s.invitationRepo.Delete(inv.ID).Exec(); err != nil {
		return status.StatusWriteDBError
	}
	return nil
}

func (s *InvitationService) CreateInvitationRecord(record *model.InvitationRecord) error {
	if record == nil || record.Code == "" || record.Status == "" {
		return status.StatusParamError
	}
	if record.UsedAt.IsZero() {
		record.UsedAt = time.Now()
	}
	if err := s.invitationRepo.CreateRecord(record).Exec(); err != nil {
		return status.StatusWriteDBError
	}
	return nil
}

func (s *InvitationService) CreateInvitationRecordWithTx(tx *gorm.DB, record *model.InvitationRecord) error {
	if record == nil || record.Code == "" || record.Status == "" {
		return status.StatusParamError
	}
	if record.UsedAt.IsZero() {
		record.UsedAt = time.Now()
	}
	if err := s.invitationRepo.CreateRecord(record).WithTx(tx).Exec(); err != nil {
		return status.StatusWriteDBError
	}
	return nil
}

func (s *InvitationService) ListInvitationRecords(req *dto.ListInvitationRecordsRequest) ([]model.InvitationRecord, int64, error) {
	if req == nil {
		return nil, 0, status.StatusParamError
	}
	return s.invitationRepo.ListRecords().
		WithCode(req.Code).
		WithStatus(req.Status).
		WithUsedAtRange(req.UsedAtStart, req.UsedAtEnd).
		WithCreatorID(req.CreatorID).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize).
		ExecWithTotal()
}

func (s *InvitationService) ValidateAndUse(code string) error {
	ok, err := s.invitationRepo.ValidateAndUse(code).Exec()
	if err != nil {
		return status.StatusWriteDBError
	}
	if ok {
		return nil
	}

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

	return status.StatusInvitationInvalid
}

func (s *InvitationService) ValidateAndUseWithTx(tx *gorm.DB, code string) error {
	ok, err := s.invitationRepo.ValidateAndUse(code).WithTx(tx).Exec()
	if err != nil {
		return status.StatusWriteDBError
	}
	if ok {
		return nil
	}

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

	return status.StatusInvitationInvalid
}

func (s *InvitationService) UpdateCode(code string, newExpireTime *time.Time, newMaxUsedCnt *int64, isDisabled *bool, note *string) error {
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
	if note != nil {
		invitation.Note = note
	}
	if err := s.invitationRepo.Update(invitation).Exec(); err != nil {
		return status.StatusWriteDBError
	}
	return nil
}

var InvitationSvc = NewInvitationService()
