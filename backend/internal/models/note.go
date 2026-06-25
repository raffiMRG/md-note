package models

import "time"

type Note struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"size:255" json:"title"`
	Content       string    `gorm:"type:mediumtext" json:"content"`
	CreatedBy     *uint64   `json:"created_by"`
	UpdatedBy     *uint64   `json:"updated_by"`
	CreatedByUser *User     `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	UpdatedByUser *User     `gorm:"foreignKey:UpdatedBy" json:"updated_by_user,omitempty"`
	Tags          []Tag     `gorm:"many2many:note_tags;" json:"tags"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
