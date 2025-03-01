package entity

// WordDialog entity for many-to-many relationship
type WordDialog struct {
	DialogID int64 `gorm:"column:dialog_id;primaryKey"`
	WordID   int64 `gorm:"column:word_id;primaryKey"`
}

func (wd *WordDialog) TableName() string {
	return "word_dialog"
}
