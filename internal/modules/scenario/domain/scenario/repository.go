package scenario

import (
	"context"

	"github.com/google/uuid"
)

type ScenarioPersistenceRepository interface {
	Create(ctx context.Context, scenario *Scenario) (*Scenario, error)
	Update(ctx context.Context, scenario *Scenario) (*Scenario, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Scenario, error)
	GetList(ctx context.Context, userID uuid.UUID, beforeID uuid.UUID, limit int) ([]*Scenario, error)
}
