package authuserpg

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID           uuid.UUID `gorm:"primarykey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Firstname    string    `gorm:"not null"`
	Lastname     string
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
	DeletedAt    gorm.DeletedAt
}
