package models

type Tag struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100" json:"name"`
	Slug string `gorm:"size:100;uniqueIndex" json:"slug"`
}
