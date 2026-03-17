package jwtmanager

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/devathh/scene-ai/internal/common/config"
	jwtdomain "github.com/devathh/scene-ai/internal/modules/auth/domain/jwt"
	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	cfg     *config.Config
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

func New(cfg *config.Config, loader jwtdomain.KeyLoader) (*JWTManager, error) {
	private, err := loader.LoadPrivate()
	if err != nil {
		return nil, err
	}

	public, err := loader.LoadPublic()
	if err != nil {
		return nil, err
	}

	return &JWTManager{
		cfg:     cfg,
		private: private,
		public:  public,
	}, nil
}

func (jm *JWTManager) GenerateAccess(userID uuid.UUID) (string, error) {
	claims := jwtdomain.Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "scene-ai",
			Subject:   "scene-user",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(jm.cfg.Cache.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(jm.private)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return tokenString, nil
}

func (jm *JWTManager) GenerateRefresh() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(buf), nil
}

func (jm *JWTManager) Validate(tokenString string) (*jwtdomain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtdomain.Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, consts.ErrInvalidToken
		}

		return jm.public, nil
	})
	if err != nil {
		return nil, consts.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*jwtdomain.Claims); ok {
		return claims, nil
	}

	return nil, consts.ErrInvalidToken
}
