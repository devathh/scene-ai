package scenarioredis

import (
	"encoding/json"

	"github.com/devathh/scene-ai/internal/modules/ai/domain/scenario"
)

func ToBytes(domain *scenario.Scenario) ([]byte, error) {
	model := ScenarioModel{
		ID:                domain.ID(),
		AuthorID:          domain.AuthorID(),
		Title:             domain.Title(),
		ScenarioPrompt:    domain.ScenarioPrompt(),
		GlobalStylePrompt: domain.GlobalStylePrompt(),
		Status:            domain.Status().Int(),
		Scenes:            make([]SceneModel, len(domain.Scenes())),
		CreatedAt:         domain.CreatedAt(),
		UpdatedAt:         domain.UpdatedAt(),
	}

	for idx, scene := range domain.Scenes() {
		model.Scenes[idx] = SceneModel{
			ID:          scene.ID(),
			Order:       scene.Order(),
			Title:       scene.Title(),
			Duration:    scene.Duration(),
			VideoPrompt: scene.VideoPrompt(),
		}
	}

	return json.Marshal(&model)
}

func ToDomain(bytes []byte) (*scenario.Scenario, error) {
	var model ScenarioModel
	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, err
	}

	scenes := make([]scenario.Scene, len(model.Scenes))
	for idx, scene := range model.Scenes {
		scenes[idx] = scenario.FromScene(
			scene.ID,
			scene.Order,
			scene.Title,
			scene.Duration,
			scene.VideoPrompt,
		)
	}

	status, _ := scenario.NewStatus(model.Status)
	return scenario.From(
		model.ID,
		model.AuthorID,
		model.Title,
		model.ScenarioPrompt,
		model.GlobalStylePrompt,
		status,
		scenes,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}

func ToBytesScene(scene *scenario.Scene) ([]byte, error) {
	return json.Marshal(SceneModel{
		ID:          scene.ID(),
		Order:       scene.Order(),
		Title:       scene.Title(),
		Duration:    scene.Duration(),
		VideoPrompt: scene.VideoPrompt(),
	})
}

func ToDomainScene(bytes []byte) (*scenario.Scene, error) {
	var model SceneModel
	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, err
	}

	scene := scenario.FromScene(
		model.ID,
		model.Order,
		model.Title,
		model.Duration,
		model.VideoPrompt,
	)

	return &scene, nil
}
