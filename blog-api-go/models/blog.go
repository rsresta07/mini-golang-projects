package models

import "time"

type Blog struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Content   string `gorm:"type:text"`
	UserID    uint   `gorm:"not null"`
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}
