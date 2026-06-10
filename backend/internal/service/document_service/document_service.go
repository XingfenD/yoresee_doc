package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/repository/attachment_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_version_repo"
	doc_yjs_snapshot_repo "github.com/XingfenD/yoresee_doc/internal/repository/document_yjs_snapshot_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/knowledge_base_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/template_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"gorm.io/gorm"
)

type DocumentService struct {
	db             *gorm.DB
	documentRepo   *document_repo.DocumentRepository
	userRepo       *user_repo.UserRepository
	kbRepo         *knowledge_base_repo.KnowledgeBaseRepository
	docVersionRepo *document_version_repo.DocumentVersionRepository
	snapshotRepo   *doc_yjs_snapshot_repo.DocumentYjsSnapshotRepository
	templateRepo   *template_repo.TemplateRepository
	attachmentRepo *attachment_repo.AttachmentRepository
}

func NewDocumentService(repos *repository.Repositories) *DocumentService {
	return &DocumentService{
		db:             repos.DB,
		documentRepo:   repos.Document,
		userRepo:       repos.User,
		kbRepo:         repos.KnowledgeBase,
		docVersionRepo: repos.DocumentVersion,
		snapshotRepo:   repos.DocumentYjsSnapshot,
		templateRepo:   repos.Template,
		attachmentRepo: repos.Attachment,
	}
}

var DocumentSvc *DocumentService

func (s *DocumentService) DeleteDocument(id int64) error {
	return s.documentRepo.Delete(id).Exec()
}

func (s *DocumentService) DeleteDocumentByExternalID(ctx context.Context, externalID string) error {
	docID, err := s.documentRepo.GetIDByExternalID(externalID).Exec(ctx)
	if err != nil {
		return err
	}
	return s.DeleteDocument(docID)
}
