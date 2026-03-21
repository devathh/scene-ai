package scenario

import (
	"context"

	"github.com/google/uuid"
)

type ScenarioCacheRepository interface {
	SetScenario(ctx context.Context, scenario *Scenario) error
	GetScenario(ctx context.Context, id uuid.UUID) (*Scenario, error)
	AddScene(ctx context.Context, scene *Scene, scenarioID uuid.UUID) error
	GetScenes(ctx context.Context, scenarioID uuid.UUID) ([]*Scene, error)
	PublishScene(ctx context.Context, scene *Scene, scenarioID uuid.UUID) error
	SubscribeScene(ctx context.Context, scenarioID uuid.UUID, f func(scene *Scene)) error
}
