package scenariopg

import (
	"context"
	"errors"
	"fmt"

	"github.com/devathh/scene-ai/internal/modules/scenario/domain/scenario"
	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	minLimit = 1
	maxLimit = 1000
)

type ScenarioRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *ScenarioRepository {
	return &ScenarioRepository{
		db: db,
	}
}

func (sr *ScenarioRepository) Create(ctx context.Context, scenario *scenario.Scenario) (*scenario.Scenario, error) {
	model := ToModel(scenario)
	if err := sr.db.WithContext(ctx).Create(&model).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		return nil, fmt.Errorf("failed to create scenario in db: %w", err)
	}

	return ToDomain(model), nil
}

func (sr *ScenarioRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	if err := sr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exists bool
		if err := tx.Model(&ScenarioModel{}).
			Where("id = ? AND author_id = ?", id, userID).
			Select("count(*) > 0").
			Find(&exists).Error; err != nil {
			return err
		}
		if !exists {
			return consts.ErrScenarioNotFound
		}

		if err := tx.Where("scenario_id = ?", id).Delete(&SceneModel{}).Error; err != nil {
			return err
		}

		return tx.Where("id = ?", id).Delete(&ScenarioModel{}).Error
	}); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		if errors.Is(err, consts.ErrScenarioNotFound) {
			return err
		}

		return fmt.Errorf("failed to delete scenario from db: %w", err)
	}

	return nil
}

func (sr *ScenarioRepository) Update(ctx context.Context, scenario *scenario.Scenario) (*scenario.Scenario, error) {
	model := ToModel(scenario)
	result := sr.db.WithContext(ctx).Save(&model)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) || errors.Is(result.Error, context.Canceled) {
			return nil, errors.ErrUnsupported
		}

		return nil, fmt.Errorf("failed to update scenario: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, consts.ErrScenarioNotFound
	}

	return ToDomain(model), nil
}

func (sr *ScenarioRepository) GetByID(ctx context.Context, id uuid.UUID) (*scenario.Scenario, error) {
	var model ScenarioModel
	if err := sr.db.WithContext(ctx).Preload("Scenes").Take(&model, id).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, consts.ErrScenarioNotFound
		}

		return nil, fmt.Errorf("failed to get scenario by id from db: %w", err)
	}

	return ToDomain(&model), nil
}

func (sr *ScenarioRepository) GetList(ctx context.Context, userID, beforeID uuid.UUID, limit int) ([]*scenario.Scenario, error) {
	if limit < minLimit || limit > maxLimit {
		return nil, consts.ErrInvalidLimit
	}

	var models []ScenarioModel
	if err := sr.db.WithContext(ctx).
		Where("author_id = ? AND id > ?", userID, beforeID).
		Find(&models).
		Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get list of scenarios: %w", err)
	}

	scenarios := make([]*scenario.Scenario, len(models))
	for idx, scenario := range models {
		scenarios[idx] = ToDomain(&scenario)
	}

	return scenarios, nil
}
