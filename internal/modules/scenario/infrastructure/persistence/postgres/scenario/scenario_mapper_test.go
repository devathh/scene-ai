package scenariopg

import (
	"testing"
	"time"

	"github.com/devathh/scene-ai/internal/modules/scenario/domain/scenario"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToModel(t *testing.T) {
	authorID := uuid.New()
	scenarioID := uuid.New()
	now := time.Now().UTC()

	domainScene, err := scenario.NewScene(1, "Scene Title", 10*time.Second, "Video Prompt")
	assert.NoError(t, err)

	domainScenario := scenario.From(
		scenarioID,
		authorID,
		"Test Title",
		"Scenario Prompt",
		"Global Style",
		scenario.STATUS_MODIFIED,
		[]scenario.Scene{domainScene},
		now,
		now.Add(1*time.Hour),
	)

	model := ToModel(domainScenario)

	assert.NotNil(t, model)
	assert.Equal(t, scenarioID, model.ID)
	assert.Equal(t, authorID, model.AuthorID)
	assert.Equal(t, "Test Title", model.Title)
	assert.Equal(t, "Scenario Prompt", model.ScenarioPrompt)
	assert.Equal(t, "Global Style", model.GlobalStylePrompt)
	assert.Equal(t, scenario.STATUS_MODIFIED.Int(), model.Status)
	assert.Equal(t, now, model.CreatedAt)
	assert.Equal(t, now.Add(1*time.Hour), model.UpdatedAt)

	assert.Len(t, model.Scenes, 1)
	sceneModel := model.Scenes[0]
	assert.Equal(t, domainScene.ID(), sceneModel.ID)
	assert.Equal(t, scenarioID, sceneModel.ScenarioID)
	assert.Equal(t, 1, sceneModel.Order)
	assert.Equal(t, "Scene Title", sceneModel.Title)
	assert.Equal(t, 10*time.Second, sceneModel.Duration)
	assert.Equal(t, "Video Prompt", sceneModel.VideoPrompt)
}

func TestToDomain(t *testing.T) {
	authorID := uuid.New()
	scenarioID := uuid.New()
	sceneID := uuid.New()
	now := time.Now().UTC()

	model := &ScenarioModel{
		ID:                scenarioID,
		AuthorID:          authorID,
		Title:             "DB Title",
		ScenarioPrompt:    "DB Scenario Prompt",
		GlobalStylePrompt: "DB Global Style",
		Status:            scenario.STATUS_GENERATED.Int(),
		CreatedAt:         now,
		UpdatedAt:         now.Add(2 * time.Hour),
		Scenes: []SceneModel{
			{
				ID:          sceneID,
				ScenarioID:  scenarioID,
				Order:       5,
				Title:       "DB Scene Title",
				Duration:    30 * time.Second,
				VideoPrompt: "DB Video Prompt",
			},
		},
	}

	domain := ToDomain(model)

	assert.NotNil(t, domain)
	assert.Equal(t, scenarioID, domain.ID())
	assert.Equal(t, authorID, domain.AuthorID())
	assert.Equal(t, "DB Title", domain.Title())
	assert.Equal(t, "DB Scenario Prompt", domain.ScenarioPrompt())
	assert.Equal(t, "DB Global Style", domain.GlobalStylePrompt())
	assert.Equal(t, scenario.STATUS_GENERATED, domain.Status())
	assert.Equal(t, now, domain.CreatedAt())
	assert.Equal(t, now.Add(2*time.Hour), domain.UpdatedAt())

	scenes := domain.Scenes()
	assert.Len(t, scenes, 1)
	dbScene := scenes[0]
	assert.Equal(t, sceneID, dbScene.ID())
	assert.Equal(t, 5, dbScene.Order())
	assert.Equal(t, "DB Scene Title", dbScene.Title())
	assert.Equal(t, 30*time.Second, dbScene.Duration())
	assert.Equal(t, "DB Video Prompt", dbScene.VideoPrompt())
}

func TestMapperRoundTrip(t *testing.T) {
	authorID := uuid.New()
	originalID := uuid.New()
	now := time.Now().UTC()

	domainScene, _ := scenario.NewScene(1, "Round Trip Scene", 5*time.Second, "Prompt")
	originalDomain := scenario.From(
		originalID,
		authorID,
		"Round Trip Title",
		"Round Trip Scenario",
		"Round Trip Style",
		scenario.STATUS_MODIFIED,
		[]scenario.Scene{domainScene},
		now,
		now,
	)

	model := ToModel(originalDomain)

	restoredDomain := ToDomain(model)

	assert.Equal(t, originalDomain.ID(), restoredDomain.ID())
	assert.Equal(t, originalDomain.AuthorID(), restoredDomain.AuthorID())
	assert.Equal(t, originalDomain.Title(), restoredDomain.Title())
	assert.Equal(t, originalDomain.ScenarioPrompt(), restoredDomain.ScenarioPrompt())
	assert.Equal(t, originalDomain.GlobalStylePrompt(), restoredDomain.GlobalStylePrompt())
	assert.Equal(t, originalDomain.Status(), restoredDomain.Status())
	assert.Equal(t, originalDomain.CreatedAt(), restoredDomain.CreatedAt())
	assert.Equal(t, originalDomain.UpdatedAt(), restoredDomain.UpdatedAt())

	origScenes := originalDomain.Scenes()
	restScenes := restoredDomain.Scenes()
	assert.Len(t, restScenes, len(origScenes))
	if len(origScenes) > 0 {
		assert.Equal(t, origScenes[0].ID(), restScenes[0].ID())
		assert.Equal(t, origScenes[0].Order(), restScenes[0].Order())
		assert.Equal(t, origScenes[0].Title(), restScenes[0].Title())
		assert.Equal(t, origScenes[0].Duration(), restScenes[0].Duration())
		assert.Equal(t, origScenes[0].VideoPrompt(), restScenes[0].VideoPrompt())
	}
}

func TestToDomainInvalidStatus(t *testing.T) {
	model := &ScenarioModel{
		ID:                uuid.New(),
		AuthorID:          uuid.New(),
		Title:             "Test",
		ScenarioPrompt:    "Prompt",
		GlobalStylePrompt: "Style",
		Status:            999,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Scenes:            []SceneModel{},
	}

	domain := ToDomain(model)
	assert.NotNil(t, domain)
	assert.Equal(t, scenario.STATUS_UNKNOWN, domain.Status())
}
