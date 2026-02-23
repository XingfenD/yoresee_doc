package model

type DocKnowledgeRelation struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	DocumentID  int64  `gorm:"not null;index" json:"document_id"`
	KnowledgeID *int64 `gorm:"index" json:"knowledge_id"`      // in knowledge_base
	OwnerID     *int64 `gorm:"not null;index" json:"owner_id"` // or in user own doc
}

func (DocKnowledgeRelation) TableName() string {
	return "doc_knowledge_relations"
}
