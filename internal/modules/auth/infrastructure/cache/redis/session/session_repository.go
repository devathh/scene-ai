package sessionredis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/devathh/scene-ai/internal/common/config"
	"github.com/devathh/scene-ai/internal/modules/auth/domain/session"
	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	cfg    *config.Config
	client *redis.Client
}

func New(cfg *config.Config, client *redis.Client) *SessionRepository {
	return &SessionRepository{
		cfg:    cfg,
		client: client,
	}
}

func (sr *SessionRepository) Set(ctx context.Context, session *session.Session, refresh string) error {
	bytes, err := ToBytes(session)
	if err != nil {
		return err
	}

	pipe := sr.client.TxPipeline()
	pipe.Set(ctx, sr.getRefreshKey(refresh), bytes, sr.cfg.JWT.RefreshTTL)

	key := sr.getSessionsKey(session.UserID)
	pipe.SAdd(ctx, key, refresh)
	pipe.Expire(ctx, key, sr.cfg.JWT.RefreshTTL)

	if _, err := pipe.Exec(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to exec pipe: %w", err)
	}

	return nil
}

func (sr *SessionRepository) Get(ctx context.Context, refresh string) (*session.Session, error) {
	bytes, err := sr.client.Get(ctx, sr.getRefreshKey(refresh)).Bytes()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, redis.Nil) {
			return nil, consts.ErrSessionNotFound
		}

		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return ToModel(bytes)
}

func (sr *SessionRepository) Del(ctx context.Context, refresh string) error {
	if err := sr.client.Del(ctx, sr.getRefreshKey(refresh)).Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (sr *SessionRepository) DelAll(ctx context.Context, userID uuid.UUID) error {
	oldKeys, err := sr.client.SMembers(ctx, sr.getSessionsKey(userID)).Result()
	if err != nil {
		return fmt.Errorf("failed to get all sessions: %w", err)
	}

	keys := make([]string, len(oldKeys))
	for i, key := range oldKeys {
		keys[i] = sr.getRefreshKey(key)
	}

	pipe := sr.client.TxPipeline()
	pipe.Del(ctx, keys...)
	pipe.Del(ctx, sr.getSessionsKey(userID))

	if _, err := pipe.Exec(context.Background()); err != nil {
		return fmt.Errorf("failed to exec del's pipe: %w", err)
	}

	return nil
}

func (sr *SessionRepository) getRefreshKey(refresh string) string {
	hash := sha256.Sum256([]byte(refresh))
	return fmt.Sprintf("rs:%s", hex.EncodeToString(hash[:]))
}

func (sr *SessionRepository) getSessionsKey(userID uuid.UUID) string {
	return fmt.Sprintf("us:%s", userID.String())
}
