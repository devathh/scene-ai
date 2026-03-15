package authuserpg

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID           uuid.UUID `gorm:"primarykey"`
	Firstname    string    `gorm:"not null"`
	Lastname     string
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
	DeletedAt    time.Time
}
