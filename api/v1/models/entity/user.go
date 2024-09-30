package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	Name       string         `json:"name"`
	Email      string         `json:"email"`
	Password   string         `json:"-" gorm:"column:password"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	VerifiedAt string         `json:"verified_at"`
}
