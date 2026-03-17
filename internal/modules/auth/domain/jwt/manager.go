package jwtdomain

import "github.com/google/uuid"

type JWTManager interface {
	GenerateAccess(userID uuid.UUID) (string, error)
	GenerateRefresh() (string, error)
	Validate(tokenString string) (*Claims, error)
}
