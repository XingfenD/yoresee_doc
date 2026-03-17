package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func MustInit() {
	document_repo.Init(storage.KVS)
}
