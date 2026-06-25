package models

import "time"

type CORSOrigin struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Origin    string    `gorm:"size:255;uniqueIndex" json:"origin"`
	CreatedAt time.Time `json:"created_at"`
}
