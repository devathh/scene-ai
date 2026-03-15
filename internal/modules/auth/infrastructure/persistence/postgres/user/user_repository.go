package authuserpg

import (
	"context"
	"errors"
	"fmt"

	"github.com/devathh/scene-ai/internal/modules/auth/domain/user"
	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Save(ctx context.Context, user *user.User) (*user.User, error) {
	model := ToModel(user)
	if err := ur.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return ToDomain(model), nil
}

func (ur *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := ur.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) || errors.Is(result.Error, context.Canceled) {
			return result.Error
		}

		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return consts.ErrUserNotFound
	}

	return nil
}

func (ur *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var model UserModel
	if err := ur.db.WithContext(ctx).Take(&model, id).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, consts.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to take user: %w", err)
	}

	return ToDomain(&model), nil
}
