package scenariopg

import "github.com/devathh/scene-ai/internal/modules/scenario/domain/scenario"

func ToModel(domain *scenario.Scenario) *ScenarioModel {
	scenes := make([]SceneModel, len(domain.Scenes()))
	for idx, scene := range domain.Scenes() {
		scenes[idx] = SceneModel{
			ID:          scene.ID(),
			ScenarioID:  domain.ID(),
			Order:       scene.Order(),
			Title:       scene.Title(),
			Duration:    scene.Duration(),
			VideoPrompt: scene.VideoPrompt(),
		}
	}

	return &ScenarioModel{
		ID:                domain.ID(),
		AuthorID:          domain.AuthorID(),
		Title:             domain.Title(),
		ScenarioPrompt:    domain.ScenarioPrompt(),
		GlobalStylePrompt: domain.GlobalStylePrompt(),
		Status:            domain.Status().Int(),
		Scenes:            scenes,
		CreatedAt:         domain.CreatedAt(),
		UpdatedAt:         domain.UpdatedAt(),
	}
}

func ToDomain(model *ScenarioModel) *scenario.Scenario {
	status, _ := scenario.NewStatus(model.Status)
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
	)
}
