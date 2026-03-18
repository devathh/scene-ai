package jwtdomain

import (
	"github.com/devathh/scene-ai/internal/common/claims"
	"github.com/google/uuid"
)

type JWTManager interface {
	GenerateAccess(userID uuid.UUID) (string, error)
	GenerateRefresh() (string, error)
	Validate(tokenString string) (*claims.Claims, error)
}
