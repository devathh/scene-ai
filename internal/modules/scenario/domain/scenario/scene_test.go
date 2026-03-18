package scenario

import (
	"testing"
	"time"

	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewScene(t *testing.T) {
	tests := []struct {
		name           string
		order          int
		title          string
		duration       time.Duration
		videoPrompt    string
		wantErr        error
		validateResult func(*testing.T, Scene)
	}{
		{
			name:        "ValidScene",
			order:       1,
			title:       "Intro Scene",
			duration:    10 * time.Second,
			videoPrompt: "A dark forest at night",
			wantErr:     nil,
			validateResult: func(t *testing.T, scene Scene) {
				assert.Equal(t, 1, scene.Order())
				assert.Equal(t, "Intro Scene", scene.Title())
				assert.Equal(t, 10*time.Second, scene.Duration())
				assert.Equal(t, "A dark forest at night", scene.VideoPrompt())
				assert.NotEmpty(t, scene.ID())
			},
		},
		{
			name:        "EmptyTitle",
			order:       1,
			title:       "   ",
			duration:    5 * time.Second,
			videoPrompt: "Prompt",
			wantErr:     consts.ErrEmptyTitle,
		},
		{
			name:        "EmptyVideoPrompt",
			order:       1,
			title:       "Scene",
			duration:    5 * time.Second,
			videoPrompt: "",
			wantErr:     consts.ErrEmptyVideoPrompt,
		},
		{
			name:        "ZeroDurationAllowed",
			order:       2,
			title:       "Flash",
			duration:    0,
			videoPrompt: "Quick flash",
			wantErr:     nil,
			validateResult: func(t *testing.T, scene Scene) {
				assert.Equal(t, time.Duration(0), scene.Duration())
			},
		},
		{
			name:        "NegativeOrderAllowed",
			order:       -1,
			title:       "Prequel",
			duration:    30 * time.Second,
			videoPrompt: "Backstory",
			wantErr:     nil,
			validateResult: func(t *testing.T, scene Scene) {
				assert.Equal(t, -1, scene.Order())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scene, err := NewScene(tt.order, tt.title, tt.duration, tt.videoPrompt)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				// Проверяем, что возвращен пустой Scene (zero value)
				assert.Equal(t, Scene{}, scene)
			} else {
				assert.NoError(t, err)
				// Scene теперь значение, а не указатель, поэтому он не может быть nil
				if tt.validateResult != nil {
					tt.validateResult(t, scene)
				}
			}
		})
	}
}

func TestSceneGetters(t *testing.T) {
	duration := 15 * time.Second
	scene, err := NewScene(5, "Main Event", duration, "Epic battle")
	assert.NoError(t, err)

	assert.NotEmpty(t, scene.ID())
	assert.Equal(t, 5, scene.Order())
	assert.Equal(t, "Main Event", scene.Title())
	assert.Equal(t, duration, scene.Duration())
	assert.Equal(t, "Epic battle", scene.VideoPrompt())
}

func TestFromScene(t *testing.T) {
	testID := uuid.New()
	order := 10
	title := "Reconstructed Scene"
	duration := 20 * time.Second
	prompt := "Reconstruction prompt"

	scene := FromScene(testID, order, title, duration, prompt)

	assert.Equal(t, testID, scene.ID())
	assert.Equal(t, order, scene.Order())
	assert.Equal(t, title, scene.Title())
	assert.Equal(t, duration, scene.Duration())
	assert.Equal(t, prompt, scene.VideoPrompt())
}
