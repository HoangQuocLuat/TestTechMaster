package entity

type Dialog struct {
	ID      int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Lang    string `gorm:"column:lang;size:50;not null"`
	Content string `gorm:"column:content;type:text;not null"`
}

func (d *Dialog) TableName() string {
	return "dialog"
}
