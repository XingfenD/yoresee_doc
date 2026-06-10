package service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/service/auth_service"
	"github.com/XingfenD/yoresee_doc/internal/service/comment_service"
	"github.com/XingfenD/yoresee_doc/internal/service/config_service"
	"github.com/XingfenD/yoresee_doc/internal/service/document_service"
	svc_iface "github.com/XingfenD/yoresee_doc/internal/service/interface"
	"github.com/XingfenD/yoresee_doc/internal/service/invitation_service"
	"github.com/XingfenD/yoresee_doc/internal/service/knowledge_base_service"
	"github.com/XingfenD/yoresee_doc/internal/service/membership_service"
	"github.com/XingfenD/yoresee_doc/internal/service/mq_service"
	"github.com/XingfenD/yoresee_doc/internal/service/notification_service"
	"github.com/XingfenD/yoresee_doc/internal/service/setting_service"
	"github.com/XingfenD/yoresee_doc/internal/service/user_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/sirupsen/logrus"
)

func RegisterTopicConsumer(h svc_iface.TopicConsumer) error {
	return mq_service.MQSvc.Consume(
		context.Background(),
		mq.BackendRedis,
		mq.ConsumeOptions{
			Topic:   h.Topic(),
			Mode:    mq.ConsumeModeFanout,
			AutoAck: true,
			OnError: mq.ErrorActionDrop,
		},
		func(ctx context.Context, message mq.Message) error {
			return h.Consume()(message.Body)
		},
	)
}

func Init(cfg *config.Config, repos *repository.Repositories) error {
	auth_service.AuthSvc = auth_service.NewAuthService(repos)
	user_service.UserSvc = user_service.NewUserService(repos)
	document_service.DocumentSvc = document_service.NewDocumentService(repos)
	comment_service.CommentSvc = comment_service.NewCommentService(repos)
	knowledge_base_service.KnowledgeBaseSvc = knowledge_base_service.NewKnowledgeBaseService(repos)
	membership_service.MembershipSvc = membership_service.NewMembershipService(repos)
	invitation_service.InvitationSvc = invitation_service.NewInvitationService(repos)
	notification_service.NotificationSvc = notification_service.NewNotificationService(repos)

	config_service.InitConfigService()
	setting_service.InitSettingService()
	if err := InitMQTopicConsumer(); err != nil {
		logrus.Errorf("[Service layer] InitMQTopicConsumer failed, err=%+v", err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "initialize message queue topic consumer failed")
	}

	return nil
}

func InitMQTopicConsumer() error {
	return nil
}
