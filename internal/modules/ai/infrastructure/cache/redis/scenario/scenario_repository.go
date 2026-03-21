package scenarioredis

import (
	"context"
	"errors"
	"fmt"

	"github.com/devathh/scene-ai/internal/common/config"
	"github.com/devathh/scene-ai/internal/modules/ai/domain/scenario"
	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ScenarioCacheRepository struct {
	cfg    *config.Config
	client *redis.Client
}

func NewClient(cfg *config.Config, client *redis.Client) *ScenarioCacheRepository {
	return &ScenarioCacheRepository{
		cfg:    cfg,
		client: client,
	}
}

func (scr *ScenarioCacheRepository) SetScenario(ctx context.Context, scenario *scenario.Scenario) error {
	key := scr.getScenarioKey(scenario.ID())
	bytes, err := ToBytes(scenario)
	if err != nil {
		return fmt.Errorf("failed to convert domain into bytes: %w", err)
	}

	if err := scr.client.Set(ctx, key, bytes, scr.cfg.Cache.ScenarioGenerationTTL).Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to set scenario into cache: %w", err)
	}

	return nil
}

func (scr *ScenarioCacheRepository) GetScenario(ctx context.Context, id uuid.UUID) (*scenario.Scenario, error) {
	key := scr.getScenarioKey(id)
	bytes, err := scr.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, redis.Nil) {
			return nil, consts.ErrScenarioNotFound
		}

		return nil, fmt.Errorf("failed to get scenario: %w", err)
	}

	domain, err := ToDomain(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert bytes to domain: %w", err)
	}

	return domain, nil
}

func (scr *ScenarioCacheRepository) AddScene(ctx context.Context, scene *scenario.Scene, scenarioID uuid.UUID) error {
	bytes, err := ToBytesScene(scene)
	if err != nil {
		return fmt.Errorf("failed to convert bytes into model: %w", err)
	}

	key := scr.getScenesKey(scenarioID)

	pipe := scr.client.TxPipeline()
	pipe.SAdd(ctx, key, bytes)
	pipe.Expire(ctx, key, scr.cfg.Cache.ScenarioGenerationTTL)

	if _, err := pipe.Exec(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to exec pipe: %w", err)
	}

	return nil
}

func (scr *ScenarioCacheRepository) GetScenes(ctx context.Context, scenarioID uuid.UUID) ([]*scenario.Scene, error) {
	key := scr.getScenesKey(scenarioID)

	rawBytes, err := scr.client.SMembers(ctx, key).Result()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, redis.Nil) {
			return nil, consts.ErrScenarioNotFound
		}

		return nil, fmt.Errorf("failed to get scenes: %w", err)
	}

	domains := make([]*scenario.Scene, len(rawBytes))
	for idx, bytes := range rawBytes {
		domains[idx], err = ToDomainScene([]byte(bytes))
		if err != nil {
			return nil, fmt.Errorf("failed to convert bytes into scene: %w", err)
		}
	}

	return domains, nil
}

func (scr *ScenarioCacheRepository) PublishScene(ctx context.Context, scene *scenario.Scene, scenarioID uuid.UUID) error {
	bytes, err := ToBytesScene(scene)
	if err != nil {
		return fmt.Errorf("failed to convert bytes into model: %w", err)
	}

	key := scr.pubsubKey(scenarioID)

	if err := scr.client.Publish(ctx, key, bytes).Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to publish scene: %w", err)
	}

	return nil
}

func (scr *ScenarioCacheRepository) SubscribeScene(ctx context.Context, scenarioID uuid.UUID, f func(scene *scenario.Scene)) error {
	key := scr.pubsubKey(scenarioID)

	pubsub := scr.client.Subscribe(ctx, key)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		domain, err := ToDomainScene([]byte(msg.Payload))
		if err != nil {
			return fmt.Errorf("failed to convert msg to domain: %w", err)
		}

		f(domain)
	}

	return nil
}

func (scr *ScenarioCacheRepository) getScenarioKey(scenarioID uuid.UUID) string {
	return fmt.Sprintf("sk:%s", scenarioID.String())
}

func (scr *ScenarioCacheRepository) getScenesKey(scenarioID uuid.UUID) string {
	return fmt.Sprintf("scenes:%s", scenarioID.String())
}

func (scr *ScenarioCacheRepository) pubsubKey(scenarioID uuid.UUID) string {
	return fmt.Sprintf("pbscenes:%s", scenarioID.String())
}
