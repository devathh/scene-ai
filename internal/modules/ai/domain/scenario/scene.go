package scenario

import (
	"strings"
	"time"

	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
)

type Scene struct {
	id          uuid.UUID
	order       int
	title       string
	duration    time.Duration
	videoPrompt string
}

func NewScene(
	order int,
	title string,
	duration time.Duration,
	videoPrompt string,
) (Scene, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Scene{}, consts.ErrEmptyTitle
	}

	videoPrompt = strings.TrimSpace(videoPrompt)
	if videoPrompt == "" {
		return Scene{}, consts.ErrEmptyVideoPrompt
	}

	return Scene{
		id:          uuid.New(),
		order:       order,
		title:       title,
		duration:    duration,
		videoPrompt: videoPrompt,
	}, nil
}

func FromScene(
	id uuid.UUID,
	order int,
	title string,
	duration time.Duration,
	videoPrompt string,
) Scene {
	return Scene{
		id:          id,
		order:       order,
		title:       title,
		duration:    duration,
		videoPrompt: videoPrompt,
	}
}

func (s *Scene) ID() uuid.UUID {
	return s.id
}

func (s *Scene) Order() int {
	return s.order
}

func (s *Scene) Title() string {
	return s.title
}

func (s *Scene) Duration() time.Duration {
	return s.duration
}

func (s *Scene) VideoPrompt() string {
	return s.videoPrompt
}
