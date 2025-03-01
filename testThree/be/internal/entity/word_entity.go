package entity

// Word entity
type Word struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Lang      string `gorm:"column:lang;size:2;not null"`
	Content   string `gorm:"column:content;type:text;not null"`
	Translate string `gorm:"column:translate;type:text;not null"`
}

func (w *Word) TableName() string {
	return "word"
}
