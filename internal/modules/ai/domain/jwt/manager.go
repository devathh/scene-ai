package jwtdomain

import (
	"github.com/devathh/scene-ai/internal/common/claims"
)

type JWTManager interface {
	Validate(tokenString string) (*claims.Claims, error)
}
