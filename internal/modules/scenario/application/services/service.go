package scenarioservices

import (
	"context"
	"errors"
	"log/slog"
	"sort"
	"time"

	"github.com/devathh/scene-ai/internal/common/config"
	"github.com/devathh/scene-ai/internal/common/dtos"
	jwtdomain "github.com/devathh/scene-ai/internal/modules/auth/domain/jwt"
	"github.com/devathh/scene-ai/internal/modules/scenario/domain/scenario"
	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
)

type ScenarioService interface {
	Create(ctx context.Context, req *dtos.CreateScenarioRequest, token string) (*dtos.Scenario, error)
	Update(ctx context.Context, id uuid.UUID, req *dtos.UpdateScenarioRequest, token string) (*dtos.Scenario, error)
	Delete(ctx context.Context, id uuid.UUID, token string) error
	GetByID(ctx context.Context, id uuid.UUID) (*dtos.Scenario, error)
	GetList(ctx context.Context, beforeID uuid.UUID, limit int, token string) (*dtos.Scenarios, error)
}

type scenarioService struct {
	cfg          *config.Config
	log          *slog.Logger
	scenarioRepo scenario.ScenarioPersistenceRepository
	jwtManager   jwtdomain.JWTManager
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	scenarioRepo scenario.ScenarioPersistenceRepository,
	jwtManager jwtdomain.JWTManager,
) ScenarioService {
	return &scenarioService{
		cfg:          cfg,
		log:          log,
		scenarioRepo: scenarioRepo,
		jwtManager:   jwtManager,
	}
}

func (s *scenarioService) Create(ctx context.Context, req *dtos.CreateScenarioRequest, token string) (*dtos.Scenario, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userID, err := s.getUserIDFromToken(token)
	if err != nil {
		return nil, err
	}

	scen, err := s.createScenario(req, userID)
	if err != nil {
		return nil, err
	}

	savedScenario, err := s.scenarioRepo.Create(ctx, scen)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return nil, err
		}

		s.log.Error("failed to save scenario",
			slog.String("error", err.Error()),
		)
		return nil, consts.ErrInternalServer
	}

	return s.toDTO(savedScenario), nil
}

func (s *scenarioService) Update(ctx context.Context, id uuid.UUID, req *dtos.UpdateScenarioRequest, token string) (*dtos.Scenario, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userID, err := s.getUserIDFromToken(token)
	if err != nil {
		return nil, err
	}

	existing, err := s.scenarioRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, consts.ErrScenarioNotFound) {
			return nil, err
		}
		s.log.Error("failed to get scenario for update",
			slog.String("error", err.Error()),
		)
		return nil, consts.ErrInternalServer
	}

	if existing.AuthorID() != userID {
		return nil, consts.ErrForbidden
	}

	updatedScen, err := s.updateScenario(existing, req)
	if err != nil {
		return nil, err
	}

	savedScenario, err := s.scenarioRepo.Update(ctx, updatedScen)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}
		if errors.Is(err, consts.ErrScenarioNotFound) {
			return nil, err
		}

		s.log.Error("failed to update scenario",
			slog.String("error", err.Error()),
		)
		return nil, consts.ErrInternalServer
	}

	return s.toDTO(savedScenario), nil
}

func (s *scenarioService) Delete(ctx context.Context, id uuid.UUID, token string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	userID, err := s.getUserIDFromToken(token)
	if err != nil {
		return err
	}

	if err := s.scenarioRepo.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		if errors.Is(err, consts.ErrScenarioNotFound) {
			return err
		}

		s.log.Error("failed to delete scenario",
			slog.String("error", err.Error()),
		)
		return consts.ErrInternalServer
	}

	return nil
}

func (s *scenarioService) GetByID(ctx context.Context, id uuid.UUID) (*dtos.Scenario, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	scenario, err := s.scenarioRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, consts.ErrScenarioNotFound) {
			return nil, err
		}

		s.log.Error("failed to get scenario by id",
			slog.String("error", err.Error()),
		)
		return nil, consts.ErrInternalServer
	}

	return s.toDTO(scenario), nil
}

func (s *scenarioService) GetList(ctx context.Context, beforeID uuid.UUID, limit int, token string) (*dtos.Scenarios, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userID, err := s.getUserIDFromToken(token)
	if err != nil {
		return nil, err
	}

	scenarios, err := s.scenarioRepo.GetList(ctx, userID, beforeID, limit)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return nil, err
		}

		s.log.Error("failed to get list of scenarios",
			slog.String("error", err.Error()),
		)
		return nil, consts.ErrInternalServer
	}

	response := dtos.Scenarios{
		Scenarios: make([]*dtos.Scenario, len(scenarios)),
	}
	for idx, scene := range scenarios {
		response.Scenarios[idx] = s.toDTO(scene)
	}

	return &response, nil
}

func (s *scenarioService) createScenario(req *dtos.CreateScenarioRequest, userID uuid.UUID) (*scenario.Scenario, error) {
	scenes := make([]scenario.Scene, len(req.Scenes))
	for idx, reqScene := range req.Scenes {
		scene, err := scenario.NewScene(
			reqScene.Order,
			reqScene.Title,
			reqScene.Duration,
			reqScene.VideoPrompt,
		)
		if err != nil {
			return nil, err
		}

		scenes[idx] = scene
	}

	return scenario.New(
		userID,
		req.Title,
		req.ScenarioPrompt,
		req.GlobalStylePrompt,
		scenes,
	)
}

func (s *scenarioService) updateScenario(existing *scenario.Scenario, req *dtos.UpdateScenarioRequest) (*scenario.Scenario, error) {
	title := existing.Title()
	if req.Title != nil {
		title = *req.Title
	}

	scenarioPrompt := existing.ScenarioPrompt()
	if req.ScenarioPrompt != nil {
		scenarioPrompt = *req.ScenarioPrompt
	}

	globalStylePrompt := existing.GlobalStylePrompt()
	if req.GlobalStylePrompt != nil {
		globalStylePrompt = *req.GlobalStylePrompt
	}

	scenes := existing.Scenes()
	if len(req.Scenes) > 0 {
		scenes = make([]scenario.Scene, len(req.Scenes))
		for idx, reqScene := range req.Scenes {
			var scene scenario.Scene
			var err error

			idStr := ""
			if reqScene.ID != nil {
				idStr = *reqScene.ID
			}
			order := existing.Scenes()[idx].Order()
			if reqScene.Order != nil {
				order = *reqScene.Order
			}
			titleScene := existing.Scenes()[idx].Title()
			if reqScene.Title != nil {
				titleScene = *reqScene.Title
			}
			duration := existing.Scenes()[idx].Duration()
			if reqScene.Duration != nil {
				duration = *reqScene.Duration
			}
			prompt := existing.Scenes()[idx].VideoPrompt()
			if reqScene.VideoPrompt != nil {
				prompt = *reqScene.VideoPrompt
			}

			if idStr != "" {
				sceneID, errParse := uuid.Parse(idStr)
				if errParse != nil {
					return nil, consts.ErrInvalidID
				}
				scene = scenario.FromScene(sceneID, order, titleScene, duration, prompt)
			} else {
				scene, err = scenario.NewScene(order, titleScene, duration, prompt)
				if err != nil {
					return nil, err
				}
			}
			scenes[idx] = scene
		}
	}

	return scenario.From(
		existing.ID(),
		existing.AuthorID(),
		title,
		scenarioPrompt,
		globalStylePrompt,
		scenario.STATUS_MODIFIED,
		scenes,
		existing.CreatedAt(),
		time.Now().UTC(),
	), nil
}

func (s *scenarioService) toDTO(scen *scenario.Scenario) *dtos.Scenario {
	response := dtos.Scenario{
		ID:                scen.ID().String(),
		AuthorID:          scen.AuthorID().String(),
		Title:             scen.Title(),
		ScenarioPrompt:    scen.ScenarioPrompt(),
		GlobalStylePrompt: scen.GlobalStylePrompt(),
		Status:            scen.Status().Int(),
		Scenes:            make([]dtos.Scene, len(scen.Scenes())),
		CreatedAt:         scen.CreatedAt().UnixMilli(),
		UpdatedAt:         scen.UpdatedAt().UnixMilli(),
	}

	sort.Slice(scen.Scenes(), func(i, j int) bool {
		return scen.Scenes()[i].Order() < scen.Scenes()[j].Order()
	})
	for idx, scene := range scen.Scenes() {
		response.Scenes[idx] = dtos.Scene{
			ID:          scene.ID().String(),
			Order:       scene.Order(),
			Title:       scene.Title(),
			Duration:    scene.Duration(),
			VideoPrompt: scene.VideoPrompt(),
		}
	}
	return &response
}

func (s *scenarioService) getUserIDFromToken(token string) (uuid.UUID, error) {
	claims, err := s.jwtManager.Validate(token)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, consts.ErrInvalidToken
	}

	return userID, nil
}
