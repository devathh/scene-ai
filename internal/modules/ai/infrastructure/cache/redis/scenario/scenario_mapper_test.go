package scenarioredis

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/devathh/scene-ai/internal/modules/ai/domain/scenario"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToBytes(t *testing.T) {
	authorID := uuid.New()
	id := uuid.New()
	now := time.Now().UTC()

	domainScenes := []scenario.Scene{
		scenario.FromScene(uuid.New(), 1, "Scene 1", 5*time.Second, "Prompt 1"),
		scenario.FromScene(uuid.New(), 2, "Scene 2", 10*time.Second, "Prompt 2"),
	}

	domainScenario := scenario.From(
		id,
		authorID,
		"Test Title",
		"Test Scenario Prompt",
		"Test Global Style",
		scenario.STATUS_GENERATED,
		domainScenes,
		now,
		now,
	)

	bytes, err := ToBytes(domainScenario)

	require.NoError(t, err)
	require.NotNil(t, bytes)

	var model ScenarioModel
	err = json.Unmarshal(bytes, &model)
	require.NoError(t, err)

	assert.Equal(t, id, model.ID)
	assert.Equal(t, authorID, model.AuthorID)
	assert.Equal(t, "Test Title", model.Title)
	assert.Equal(t, "Test Scenario Prompt", model.ScenarioPrompt)
	assert.Equal(t, "Test Global Style", model.GlobalStylePrompt)
	assert.Equal(t, int(scenario.STATUS_GENERATED), model.Status)
	assert.Equal(t, now, model.CreatedAt)
	assert.Equal(t, now, model.UpdatedAt)
	assert.Len(t, model.Scenes, 2)

	assert.Equal(t, domainScenes[0].ID(), model.Scenes[0].ID)
	assert.Equal(t, domainScenes[0].Order(), model.Scenes[0].Order)
	assert.Equal(t, domainScenes[0].Title(), model.Scenes[0].Title)
	assert.Equal(t, domainScenes[0].Duration(), model.Scenes[0].Duration)
	assert.Equal(t, domainScenes[0].VideoPrompt(), model.Scenes[0].VideoPrompt)
}

func TestToDomain(t *testing.T) {
	authorID := uuid.New()
	id := uuid.New()
	now := time.Now().UTC()
	sceneID := uuid.New()

	model := ScenarioModel{
		ID:                id,
		AuthorID:          authorID,
		Title:             "Test Title",
		ScenarioPrompt:    "Test Scenario Prompt",
		GlobalStylePrompt: "Test Global Style",
		Status:            int(scenario.STATUS_MODIFIED),
		CreatedAt:         now,
		UpdatedAt:         now,
		Scenes: []SceneModel{
			{
				ID:          sceneID,
				Order:       1,
				Title:       "Scene Title",
				Duration:    15 * time.Second,
				VideoPrompt: "Scene Video Prompt",
			},
		},
	}

	bytes, err := json.Marshal(model)
	require.NoError(t, err)

	domainScenario, err := ToDomain(bytes)

	require.NoError(t, err)
	require.NotNil(t, domainScenario)

	assert.Equal(t, id, domainScenario.ID())
	assert.Equal(t, authorID, domainScenario.AuthorID())
	assert.Equal(t, "Test Title", domainScenario.Title())
	assert.Equal(t, "Test Scenario Prompt", domainScenario.ScenarioPrompt())
	assert.Equal(t, "Test Global Style", domainScenario.GlobalStylePrompt())
	assert.Equal(t, scenario.STATUS_MODIFIED, domainScenario.Status())
	assert.Equal(t, now, domainScenario.CreatedAt())
	assert.Equal(t, now, domainScenario.UpdatedAt())

	scenes := domainScenario.Scenes()
	require.Len(t, scenes, 1)
	assert.Equal(t, sceneID, scenes[0].ID())
	assert.Equal(t, 1, scenes[0].Order())
	assert.Equal(t, "Scene Title", scenes[0].Title())
	assert.Equal(t, 15*time.Second, scenes[0].Duration())
	assert.Equal(t, "Scene Video Prompt", scenes[0].VideoPrompt())
}

func TestToDomain_InvalidJSON(t *testing.T) {
	invalidBytes := []byte("{ invalid json }")

	domainScenario, err := ToDomain(invalidBytes)

	assert.Error(t, err)
	assert.Nil(t, domainScenario)
}

func TestToDomain_InvalidStatus(t *testing.T) {
	authorID := uuid.New()
	id := uuid.New()
	now := time.Now().UTC()

	model := ScenarioModel{
		ID:                id,
		AuthorID:          authorID,
		Title:             "Test Title",
		ScenarioPrompt:    "Test Scenario Prompt",
		GlobalStylePrompt: "Test Global Style",
		Status:            999,
		CreatedAt:         now,
		UpdatedAt:         now,
		Scenes: []SceneModel{
			{
				ID:          uuid.New(),
				Order:       1,
				Title:       "Scene Title",
				Duration:    5 * time.Second,
				VideoPrompt: "Prompt",
			},
		},
	}

	bytes, err := json.Marshal(model)
	require.NoError(t, err)

	domainScenario, err := ToDomain(bytes)

	require.NoError(t, err)
	assert.Equal(t, scenario.STATUS_UNKNOWN, domainScenario.Status())
}

func TestRoundTrip(t *testing.T) {
	authorID := uuid.New()
	id := uuid.New()
	now := time.Now().UTC()

	originalScenes := []scenario.Scene{
		scenario.FromScene(uuid.New(), 1, "First", 5*time.Second, "Prompt 1"),
		scenario.FromScene(uuid.New(), 2, "Second", 10*time.Second, "Prompt 2"),
	}

	original := scenario.From(
		id,
		authorID,
		"Round Trip Title",
		"Round Trip Scenario",
		"Round Trip Style",
		scenario.STATUS_GENERATING,
		originalScenes,
		now,
		now,
	)

	bytes, err := ToBytes(original)
	require.NoError(t, err)

	restored, err := ToDomain(bytes)
	require.NoError(t, err)

	assert.Equal(t, original.ID(), restored.ID())
	assert.Equal(t, original.AuthorID(), restored.AuthorID())
	assert.Equal(t, original.Title(), restored.Title())
	assert.Equal(t, original.ScenarioPrompt(), restored.ScenarioPrompt())
	assert.Equal(t, original.GlobalStylePrompt(), restored.GlobalStylePrompt())
	assert.Equal(t, original.Status(), restored.Status())
	assert.Equal(t, original.CreatedAt(), restored.CreatedAt())
	assert.Equal(t, original.UpdatedAt(), restored.UpdatedAt())

	origScenes := original.Scenes()
	restScenes := restored.Scenes()
	require.Len(t, restScenes, len(origScenes))

	for i := range origScenes {
		assert.Equal(t, origScenes[i].ID(), restScenes[i].ID())
		assert.Equal(t, origScenes[i].Order(), restScenes[i].Order())
		assert.Equal(t, origScenes[i].Title(), restScenes[i].Title())
		assert.Equal(t, origScenes[i].Duration(), restScenes[i].Duration())
		assert.Equal(t, origScenes[i].VideoPrompt(), restScenes[i].VideoPrompt())
	}
}

func TestToBytesScene(t *testing.T) {
	sceneID := uuid.New()
	originalScene := scenario.FromScene(
		sceneID,
		3,
		"Single Scene",
		25*time.Second,
		"Detailed video prompt",
	)

	bytes, err := ToBytesScene(&originalScene)

	require.NoError(t, err)
	require.NotNil(t, bytes)

	var model SceneModel
	err = json.Unmarshal(bytes, &model)
	require.NoError(t, err)

	assert.Equal(t, sceneID, model.ID)
	assert.Equal(t, 3, model.Order)
	assert.Equal(t, "Single Scene", model.Title)
	assert.Equal(t, 25*time.Second, model.Duration)
	assert.Equal(t, "Detailed video prompt", model.VideoPrompt)
}

func TestToDomainScene(t *testing.T) {
	sceneID := uuid.New()

	model := SceneModel{
		ID:          sceneID,
		Order:       5,
		Title:       "Restored Scene",
		Duration:    40 * time.Second,
		VideoPrompt: "Restored prompt",
	}

	bytes, err := json.Marshal(model)
	require.NoError(t, err)

	domainScene, err := ToDomainScene(bytes)

	require.NoError(t, err)
	require.NotNil(t, domainScene)

	assert.Equal(t, sceneID, domainScene.ID())
	assert.Equal(t, 5, domainScene.Order())
	assert.Equal(t, "Restored Scene", domainScene.Title())
	assert.Equal(t, 40*time.Second, domainScene.Duration())
	assert.Equal(t, "Restored prompt", domainScene.VideoPrompt())
}

func TestToDomainScene_InvalidJSON(t *testing.T) {
	invalidBytes := []byte("{ bad json }")

	scene, err := ToDomainScene(invalidBytes)

	assert.Error(t, err)
	assert.Nil(t, scene)
}

func TestSceneRoundTrip(t *testing.T) {
	sceneID := uuid.New()
	original := scenario.FromScene(
		sceneID,
		10,
		"Round Trip Scene",
		100*time.Millisecond,
		"Prompt for round trip",
	)

	bytes, err := ToBytesScene(&original)
	require.NoError(t, err)

	restored, err := ToDomainScene(bytes)
	require.NoError(t, err)

	assert.Equal(t, original.ID(), restored.ID())
	assert.Equal(t, original.Order(), restored.Order())
	assert.Equal(t, original.Title(), restored.Title())
	assert.Equal(t, original.Duration(), restored.Duration())
	assert.Equal(t, original.VideoPrompt(), restored.VideoPrompt())
}
