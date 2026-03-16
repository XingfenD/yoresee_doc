package repository

import (
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func MustInit() {
	initDocumentRepo(storage.KVS)
}
