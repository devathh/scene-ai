package openrouter

import (
	"context"

	"github.com/devathh/scene-ai/internal/modules/ai/domain/scenario"
	"github.com/google/uuid"
)

type OpenRouterRepository interface {
	GenerateScenario(ctx context.Context, prompt string, authorID uuid.UUID) (*scenario.Scenario, error)
	GenerateScenes(ctx context.Context, start *scenario.Scenario, handler func(scenario.Scene)) error
}
