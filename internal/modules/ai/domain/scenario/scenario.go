package scenario

import (
	"strings"
	"time"

	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
)

type Scenario struct {
	id                uuid.UUID // uuid v7
	authorID          uuid.UUID
	title             string
	scenarioPrompt    string
	globalStylePrompt string
	status            status
	scenes            []Scene
	createdAt         time.Time
	updatedAt         time.Time
}

func New(
	authorID uuid.UUID,
	title string,
	scenarioPrompt, globalStylePrompt string,
	scenes []Scene,
) (*Scenario, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, consts.ErrEmptyTitle
	}

	scenarioPrompt = strings.TrimSpace(scenarioPrompt)
	if scenarioPrompt == "" {
		return nil, consts.ErrEmptyScenarioPrompt
	}

	globalStylePrompt = strings.TrimSpace(globalStylePrompt)
	if globalStylePrompt == "" {
		return nil, consts.ErrEmptyGlobalStylePrompt
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, consts.ErrInvalidID
	}

	now := time.Now().UTC()
	return &Scenario{
		id:                id,
		authorID:          authorID,
		title:             title,
		scenarioPrompt:    scenarioPrompt,
		globalStylePrompt: globalStylePrompt,
		status:            STATUS_GENERATING,
		scenes:            scenes,
		createdAt:         now,
		updatedAt:         now,
	}, nil
}

func From(
	id uuid.UUID,
	authorID uuid.UUID,
	title string,
	scenarioPrompt, globalStylePrompt string,
	status status,
	scenes []Scene,
	createdAt, updatedAt time.Time,
) *Scenario {
	return &Scenario{
		id:                id,
		authorID:          authorID,
		title:             title,
		scenarioPrompt:    scenarioPrompt,
		globalStylePrompt: globalStylePrompt,
		status:            status,
		scenes:            scenes,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}

func (s *Scenario) ID() uuid.UUID {
	return s.id
}

func (s *Scenario) AuthorID() uuid.UUID {
	return s.authorID
}

func (s *Scenario) Title() string {
	return s.title
}

func (s *Scenario) ScenarioPrompt() string {
	return s.scenarioPrompt
}

func (s *Scenario) GlobalStylePrompt() string {
	return s.globalStylePrompt
}

func (s *Scenario) Status() status {
	return s.status
}

func (s *Scenario) Scenes() []Scene {
	return s.scenes
}

func (s *Scenario) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Scenario) UpdatedAt() time.Time {
	return s.updatedAt
}
