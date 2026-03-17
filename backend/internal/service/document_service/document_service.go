package document_service

import (
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_version_repo"
	doc_yjs_snapshot_repo "github.com/XingfenD/yoresee_doc/internal/repository/document_yjs_snapshot_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/knowledge_base_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
)

type DocumentService struct {
	documentRepo   *document_repo.DocumentRepository
	userRepo       *user_repo.UserRepository
	kbRepo         *knowledge_base_repo.KnowledgeBaseRepository
	docVersionRepo *document_version_repo.DocumentVersionRepository
	snapshotRepo   *doc_yjs_snapshot_repo.DocumentYjsSnapshotRepository
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentRepo: &document_repo.DocumentRepo,
		userRepo:     user_repo.UserRepo,
		kbRepo:       knowledge_base_repo.KnowledgeBaseRepo,
		snapshotRepo: doc_yjs_snapshot_repo.DocumentYjsSnapshotRepo,
	}
}

var DocumentSvc = NewDocumentService()
