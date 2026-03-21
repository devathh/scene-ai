package scenarioredis

import (
	"time"

	"github.com/google/uuid"
)

type ScenarioModel struct {
	ID                uuid.UUID    `json:"id"`
	AuthorID          uuid.UUID    `json:"author_id"`
	Title             string       `json:"title"`
	ScenarioPrompt    string       `json:"scenario_prompt"`
	GlobalStylePrompt string       `json:"global_style_prompt"`
	Status            int          `json:"status"`
	Scenes            []SceneModel `json:"scenes"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

type SceneModel struct {
	ID          uuid.UUID     `json:"id"`
	Order       int           `json:"order"`
	Title       string        `json:"title"`
	Duration    time.Duration `json:"duration"`
	VideoPrompt string        `json:"video_prompt"`
}
