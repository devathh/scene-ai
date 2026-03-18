package scenariopg

import (
	"time"

	"github.com/google/uuid"
)

type ScenarioModel struct {
	ID                uuid.UUID    `gorm:"primarykey"`
	AuthorID          uuid.UUID    `gorm:"index"`
	Title             string       `gorm:"not null"`
	ScenarioPrompt    string       `gorm:"not null"`
	GlobalStylePrompt string       `gorm:"not null"`
	Status            int          `gorm:"status"`
	Scenes            []SceneModel `gorm:"foreignkey:ScenarioID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type SceneModel struct {
	ID          uuid.UUID `gorm:"primarykey"`
	ScenarioID  uuid.UUID `gorm:"not null;index"`
	Order       int
	Title       string        `gorm:"not null"`
	Duration    time.Duration `gorm:"not null"`
	VideoPrompt string        `gorm:"not null"`
}
