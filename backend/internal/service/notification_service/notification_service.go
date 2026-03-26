package notification_service

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository/notification_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
)

type NotificationService struct {
	notificationRepo *notification_repo.NotificationRepository
	userRepo         *user_repo.UserRepository
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		notificationRepo: notification_repo.NotificationRepo,
		userRepo:         user_repo.UserRepo,
	}
}

func (s *NotificationService) CreateNotifications(req *dto.CreateNotificationRequest) error {
	if req == nil || len(req.ReceiverExternalIDs) == 0 {
		return status.StatusParamError
	}
	if strings.TrimSpace(req.Type) == "" {
		return status.StatusParamError
	}

	users, err := s.userRepo.ListByExternal(req.ReceiverExternalIDs).Exec()
	if err != nil {
		return status.StatusReadDBError
	}
	if len(users) == 0 {
		return status.StatusUserNotFound
	}

	userIDMap := make(map[string]int64, len(users))
	for _, user := range users {
		userIDMap[user.ExternalID] = user.ID
	}

	items := make([]model.Notification, 0, len(req.ReceiverExternalIDs))
	for _, externalID := range req.ReceiverExternalIDs {
		id, ok := userIDMap[externalID]
		if !ok {
			return status.StatusUserNotFound
		}
		items = append(items, model.Notification{
			ExternalID: utils.GenerateExternalID(utils.ExternalIDContextNotification),
			ReceiverID: id,
			Type:       req.Type,
			Status:     "unread",
			Title:      req.Title,
			Content:    req.Content,
			Payload:    req.PayloadJSON,
		})
	}

	if err := s.notificationRepo.CreateBatch(items).Exec(); err != nil {
		return status.StatusWriteDBError
	}
	return nil
}

func (s *NotificationService) ListNotifications(req *dto.ListNotificationsRequest) ([]model.Notification, int64, error) {
	if req == nil || strings.TrimSpace(req.UserExternalID) == "" {
		return nil, 0, status.StatusParamError
	}
	userID, err := s.userRepo.GetIDByExternalID(req.UserExternalID).Exec()
	if err != nil {
		return nil, 0, status.StatusUserNotFound
	}
	list, total, err := s.notificationRepo.List(userID).
		WithStatus(req.Status).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize).
		ExecWithTotal()
	if err != nil {
		return nil, 0, err
	}

	for i := range list {
		if strings.TrimSpace(list[i].ExternalID) == "" {
			externalID := utils.GenerateExternalID(utils.ExternalIDContextNotification)
			list[i].ExternalID = externalID
			_ = s.notificationRepo.UpdateExternalID(list[i].ID, externalID).Exec()
		}
	}
	return list, total, nil
}

func (s *NotificationService) MarkRead(req *dto.MarkNotificationsReadRequest) error {
	if req == nil || strings.TrimSpace(req.UserExternalID) == "" {
		return status.StatusParamError
	}
	userID, err := s.userRepo.GetIDByExternalID(req.UserExternalID).Exec()
	if err != nil {
		return status.StatusUserNotFound
	}
	if err := s.notificationRepo.MarkRead(userID, req.ExternalIDs).Exec(); err != nil {
		return status.StatusWriteDBError
	}
	return nil
}

func (s *NotificationService) MarkAllRead(userExternalID string) error {
	if strings.TrimSpace(userExternalID) == "" {
		return status.StatusParamError
	}
	userID, err := s.userRepo.GetIDByExternalID(userExternalID).Exec()
	if err != nil {
		return status.StatusUserNotFound
	}
	if err := s.notificationRepo.MarkAllRead(userID).Exec(); err != nil {
		return status.StatusWriteDBError
	}
	return nil
}

var NotificationSvc = NewNotificationService()
