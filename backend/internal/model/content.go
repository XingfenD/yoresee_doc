package model

type Content struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Content string `gorm:"type:text" json:"content"`
}

func (Content) TableName() string {
	return "contents"
}
