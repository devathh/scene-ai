package user

import (
	"context"

	"github.com/google/uuid"
)

type UserPersistenceRepository interface {
	Save(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}
