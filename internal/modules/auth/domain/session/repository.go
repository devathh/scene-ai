package session

import (
	"context"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Set(ctx context.Context, session *Session, refresh string) error
	Get(ctx context.Context, refresh string) (*Session, error)
	Del(ctx context.Context, refresh string) error
	DelAll(ctx context.Context, userID uuid.UUID) error
}
