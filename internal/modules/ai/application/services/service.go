package aiservices

import (
	"context"
	"errors"
	"log/slog"
	"sort"

	"github.com/devathh/scene-ai/internal/common/config"
	"github.com/devathh/scene-ai/internal/common/dtos"
	jwtdomain "github.com/devathh/scene-ai/internal/modules/ai/domain/jwt"
	"github.com/devathh/scene-ai/internal/modules/ai/domain/openrouter"
	"github.com/devathh/scene-ai/internal/modules/ai/domain/scenario"
	scenarioservices "github.com/devathh/scene-ai/internal/modules/scenario/application/services"
	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type AIService interface {
	GenerateScenario(ctx context.Context, req *dtos.GenerateScenarioRequest, token string) (uuid.UUID, error)
	GetScenes(ctx context.Context, scenarioID uuid.UUID) (*dtos.Scenes, error)
	Connect(ctx context.Context, conn *websocket.Conn, scenarioID uuid.UUID) error
	GetScenario(ctx context.Context, scenarioID uuid.UUID) (*dtos.Scenario, error)
}

type aiService struct {
	cfg             *config.Config
	log             *slog.Logger
	scenarioRepo    scenario.ScenarioCacheRepository
	orRepo          openrouter.OpenRouterRepository
	jwtMngr         jwtdomain.JWTManager
	scenarioService scenarioservices.ScenarioService
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	scenarioRepo scenario.ScenarioCacheRepository,
	orRepo openrouter.OpenRouterRepository,
	jwtMngr jwtdomain.JWTManager,
	scenarioService scenarioservices.ScenarioService,
) AIService {
	return &aiService{
		cfg:             cfg,
		log:             log,
		scenarioRepo:    scenarioRepo,
		orRepo:          orRepo,
		jwtMngr:         jwtMngr,
		scenarioService: scenarioService,
	}
}

func (as *aiService) GenerateScenario(ctx context.Context, req *dtos.GenerateScenarioRequest, token string) (uuid.UUID, error) {
	if err := ctx.Err(); err != nil {
		return uuid.Nil, err
	}

	userID, err := as.getUserIDFromToken(token)
	if err != nil {
		return uuid.Nil, err
	}

	scenario, err := as.orRepo.GenerateScenario(ctx, req.Prompt, userID)
	if err != nil {
		as.log.Error("failed to generate scenario",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return uuid.Nil, nil
	}

	if err := as.scenarioRepo.SetScenario(ctx, scenario); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return uuid.Nil, err
		}

		as.log.Error("failed to set scenario",
			slog.String("error", err.Error()),
		)
		return uuid.Nil, consts.ErrInternalServer
	}

	go as.handleNewScenario(context.Background(), scenario, token)
	return scenario.ID(), nil
}

func (as *aiService) GetScenes(ctx context.Context, scenarioID uuid.UUID) (*dtos.Scenes, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	scenes, err := as.scenarioRepo.GetScenes(ctx, scenarioID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, consts.ErrScenarioNotFound) {
			return nil, err
		}

		as.log.Error("failed to get scenes",
			slog.String("error", err.Error()),
			slog.String("scenario_id", scenarioID.String()),
		)
		return nil, consts.ErrInternalServer
	}

	response := make([]dtos.Scene, len(scenes))
	for idx, scene := range scenes {
		response[idx] = dtos.Scene{
			ID:          scene.ID().String(),
			Order:       scene.Order(),
			Title:       scene.Title(),
			Duration:    scene.Duration(),
			VideoPrompt: scene.VideoPrompt(),
		}
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].Order < response[j].Order
	})

	return &dtos.Scenes{
		Scenes: response,
	}, nil
}

func (as *aiService) Connect(ctx context.Context, conn *websocket.Conn, scenarioID uuid.UUID) error {
	if err := as.scenarioRepo.SubscribeScene(ctx, scenarioID, func(scene *scenario.Scene) {
		err := conn.WriteJSON(dtos.Scene{
			ID:          scene.ID().String(),
			Order:       scene.Order(),
			Title:       scene.Title(),
			Duration:    scene.Duration(),
			VideoPrompt: scene.VideoPrompt(),
		})
		if err != nil {
			as.log.Error("failed to send msg into websocket conn", slog.String("error", err.Error()))
		}
	}); err != nil {
		as.log.Error("failed to subscribe conn", slog.String("error", err.Error()))
		return consts.ErrInternalServer
	}

	return nil
}

func (as *aiService) GetScenario(ctx context.Context, scenarioID uuid.UUID) (*dtos.Scenario, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	scenario, err := as.scenarioRepo.GetScenario(ctx, scenarioID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, consts.ErrScenarioNotFound) {
			return nil, err
		}

		as.log.Error("failed to get scenario",
			slog.String("error", err.Error()),
			slog.String("scenario_id", scenarioID.String()),
		)
		return nil, consts.ErrInternalServer
	}

	return &dtos.Scenario{
		ID:                scenario.ID().String(),
		AuthorID:          scenario.AuthorID().String(),
		Title:             scenario.Title(),
		ScenarioPrompt:    scenario.ScenarioPrompt(),
		GlobalStylePrompt: scenario.GlobalStylePrompt(),
		Status:            scenario.Status().Int(),
		CreatedAt:         scenario.CreatedAt().UnixMilli(),
		UpdatedAt:         scenario.UpdatedAt().UnixMilli(),
	}, nil
}

func (as *aiService) handleNewScenario(ctx context.Context, rawScenario *scenario.Scenario, token string) {
	dtoScenes := []dtos.CreateSceneRequest{}

	if err := as.orRepo.GenerateScenes(ctx, rawScenario, func(s scenario.Scene) {
		if err := as.scenarioRepo.AddScene(ctx, &s, rawScenario.ID()); err != nil {
			as.log.Error("failed to add scene",
				slog.String("error", err.Error()),
				slog.String("scenario_id", rawScenario.ID().String()))
		}

		if err := as.scenarioRepo.PublishScene(ctx, &s, rawScenario.ID()); err != nil {
			as.log.Error("failed to publish scene",
				slog.String("error", err.Error()),
				slog.String("scenario_id", rawScenario.ID().String()))
		}

		dtoScenes = append(dtoScenes, dtos.CreateSceneRequest{
			Order:       s.Order(),
			Title:       s.Title(),
			Duration:    s.Duration(),
			VideoPrompt: s.VideoPrompt(),
		})
	}); err != nil {
		as.log.Error("failed to generate scenes",
			slog.String("error", err.Error()),
			slog.String("scenario_id", err.Error()),
		)
	}

	if _, err := as.scenarioService.Create(ctx, &dtos.CreateScenarioRequest{
		Title:             rawScenario.Title(),
		ScenarioPrompt:    rawScenario.ScenarioPrompt(),
		GlobalStylePrompt: rawScenario.GlobalStylePrompt(),
		Scenes:            dtoScenes,
	}, token); err != nil {
		as.log.Error("failed to create scenario",
			slog.String("error", err.Error()),
		)
	}

	as.log.Info("generated")
}

func (as *aiService) getUserIDFromToken(token string) (uuid.UUID, error) {
	claims, err := as.jwtMngr.Validate(token)
	if err != nil {
		return uuid.Nil, consts.ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, consts.ErrInvalidToken
	}

	return userID, nil
}
