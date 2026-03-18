package scenario

import (
	"testing"
	"time"

	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewScenario(t *testing.T) {
	testAuthorID := uuid.New()

	tests := []struct {
		name              string
		authorID          uuid.UUID
		title             string
		scenarioPrompt    string
		globalStylePrompt string
		scenes            []Scene
		wantErr           error
	}{
		{
			name:              "ValidScenario",
			authorID:          testAuthorID,
			title:             "Test Scenario",
			scenarioPrompt:    "A story about AI",
			globalStylePrompt: "Cinematic style",
			scenes: []Scene{
				{title: "Scene 1", videoPrompt: "AI wakes up"},
			},
			wantErr: nil,
		},
		{
			name:              "EmptyTitle",
			authorID:          testAuthorID,
			title:             "   ",
			scenarioPrompt:    "A story about AI",
			globalStylePrompt: "Cinematic style",
			scenes:            []Scene{{title: "Scene 1", videoPrompt: "AI wakes up"}},
			wantErr:           consts.ErrEmptyTitle,
		},
		{
			name:              "EmptyScenarioPrompt",
			authorID:          testAuthorID,
			title:             "Test",
			scenarioPrompt:    "",
			globalStylePrompt: "Cinematic style",
			scenes:            []Scene{{title: "Scene 1", videoPrompt: "AI wakes up"}},
			wantErr:           consts.ErrEmptyScenarioPrompt,
		},
		{
			name:              "EmptyGlobalStylePrompt",
			authorID:          testAuthorID,
			title:             "Test",
			scenarioPrompt:    "A story",
			globalStylePrompt: "   ",
			scenes:            []Scene{{title: "Scene 1", videoPrompt: "AI wakes up"}},
			wantErr:           consts.ErrEmptyGlobalStylePrompt,
		},
		{
			name:              "NoScenes",
			authorID:          testAuthorID,
			title:             "Test",
			scenarioPrompt:    "A story",
			globalStylePrompt: "Style",
			scenes:            []Scene{},
			wantErr:           consts.ErrNoScenes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Для тестов создаем сцены напрямую (value type)
			if len(tt.scenes) > 0 && tt.scenes[0].title != "" {
				if tt.scenes[0].videoPrompt == "" {
					tt.scenes[0].videoPrompt = "default prompt"
				}
			}

			scenario, err := New(tt.authorID, tt.title, tt.scenarioPrompt, tt.globalStylePrompt, tt.scenes)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, scenario)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, scenario)
				assert.Equal(t, tt.authorID, scenario.AuthorID())
				assert.Equal(t, tt.title, scenario.Title())
				assert.Equal(t, tt.scenarioPrompt, scenario.ScenarioPrompt())
				assert.Equal(t, tt.globalStylePrompt, scenario.GlobalStylePrompt())
				assert.Equal(t, STATUS_GENERATED, scenario.Status())
				assert.NotEmpty(t, scenario.ID())
				assert.WithinDuration(t, time.Now(), scenario.CreatedAt(), 2*time.Second)
				assert.Equal(t, scenario.CreatedAt(), scenario.UpdatedAt())
			}
		})
	}
}

func TestScenarioGetters(t *testing.T) {
	testAuthorID := uuid.New()
	// Создаем валидную сцену через конструктор для чистоты теста
	scene, err := NewScene(1, "S1", 10*time.Second, "P1")
	assert.NoError(t, err)
	
	s, err := New(testAuthorID, "Title", "Prompt", "Style", []Scene{scene})
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.NotEmpty(t, s.ID())
	assert.Equal(t, testAuthorID, s.AuthorID())
	assert.Equal(t, "Title", s.Title())
	assert.Equal(t, "Prompt", s.ScenarioPrompt())
	assert.Equal(t, "Style", s.GlobalStylePrompt())
	assert.Equal(t, STATUS_GENERATED, s.Status())
	assert.Len(t, s.Scenes(), 1)
	assert.False(t, s.CreatedAt().IsZero())
	assert.False(t, s.UpdatedAt().IsZero())
}
