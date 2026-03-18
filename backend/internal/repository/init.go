package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/repository/config_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/knowledge_base_repo"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func MustInit() {
	document_repo.Init(storage.KVS)
	knowledge_base_repo.Init(storage.KVS)
	config_repo.Init(storage.KVS)
}
