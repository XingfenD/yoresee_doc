package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/repository/attachment_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/comment_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_version_repo"
	doc_yjs_snapshot_repo "github.com/XingfenD/yoresee_doc/internal/repository/document_yjs_snapshot_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/invitation_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/knowledge_base_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/membership_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/notification_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/template_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repositories struct {
	DB                  *gorm.DB
	Redis               *redis.Client
	Document            *document_repo.DocumentRepository
	KnowledgeBase       *knowledge_base_repo.KnowledgeBaseRepository
	User                *user_repo.UserRepository
	Comment             *comment_repo.CommentRepository
	Invitation          *invitation_repo.InvitationRepository
	Membership          *membership_repo.MembershipRepository
	Template            *template_repo.TemplateRepository
	Attachment          *attachment_repo.AttachmentRepository
	DocumentVersion     *document_version_repo.DocumentVersionRepository
	Notification        *notification_repo.NotificationRepository
	DocumentYjsSnapshot *doc_yjs_snapshot_repo.DocumentYjsSnapshotRepository
}

func NewRepositories(db *gorm.DB, redis *redis.Client) *Repositories {
	return &Repositories{
		DB:                  db,
		Redis:               redis,
		Document:            document_repo.NewDocumentRepository(db, redis),
		KnowledgeBase:       knowledge_base_repo.NewKnowledgeBaseRepository(db, redis),
		User:                user_repo.NewUserRepository(db, redis),
		Comment:             comment_repo.NewCommentRepository(db),
		Invitation:          invitation_repo.NewInvitationRepository(db),
		Membership:          membership_repo.NewMembershipRepository(db),
		Template:            template_repo.NewTemplateRepository(db),
		Attachment:          attachment_repo.NewAttachmentRepository(db),
		DocumentVersion:     document_version_repo.NewDocumentVersionRepository(db),
		Notification:        notification_repo.NewNotificationRepository(db),
		DocumentYjsSnapshot: doc_yjs_snapshot_repo.NewDocumentYjsSnapshotRepository(db),
	}
}
