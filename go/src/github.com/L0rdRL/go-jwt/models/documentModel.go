package models

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	UserID    uint // ID пользователя, создавшего документ
	User      User `gorm:"foreignKey:UserID"`
	Name      string
	S3URL     string
	Type      string // Может быть "payment", "transport", "finance"
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
