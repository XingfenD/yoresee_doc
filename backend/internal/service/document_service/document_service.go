package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/repository/attachment_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_version_repo"
	doc_yjs_snapshot_repo "github.com/XingfenD/yoresee_doc/internal/repository/document_yjs_snapshot_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/knowledge_base_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/template_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
)

type DocumentService struct {
	documentRepo   *document_repo.DocumentRepository
	userRepo       *user_repo.UserRepository
	kbRepo         *knowledge_base_repo.KnowledgeBaseRepository
	docVersionRepo *document_version_repo.DocumentVersionRepository
	snapshotRepo   *doc_yjs_snapshot_repo.DocumentYjsSnapshotRepository
	templateRepo   *template_repo.TemplateRepository
	attachmentRepo *attachment_repo.AttachmentRepository
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentRepo:   &document_repo.DocumentRepo,
		userRepo:       user_repo.UserRepo,
		kbRepo:         knowledge_base_repo.KnowledgeBaseRepo,
		docVersionRepo: document_version_repo.DocumentVersionRepo,
		snapshotRepo:   doc_yjs_snapshot_repo.DocumentYjsSnapshotRepo,
		templateRepo:   template_repo.TemplateRepo,
		attachmentRepo: attachment_repo.AttachmentRepo,
	}
}

var DocumentSvc = NewDocumentService()

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
