package authservices

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/devathh/scene-ai/internal/common/config"
	"github.com/devathh/scene-ai/internal/common/dtos"
	jwtdomain "github.com/devathh/scene-ai/internal/modules/auth/domain/jwt"
	"github.com/devathh/scene-ai/internal/modules/auth/domain/session"
	"github.com/devathh/scene-ai/internal/modules/auth/domain/user"
	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.Token, error)
}

type authService struct {
	cfg         *config.Config
	log         *slog.Logger
	userRepo    user.UserPersistenceRepository
	sessionRepo session.SessionRepository
	jwtManager  jwtdomain.JWTManager
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	userRepo user.UserPersistenceRepository,
	sessionRepo session.SessionRepository,
	jwtManager jwtdomain.JWTManager,
) AuthService {
	return &authService{
		cfg:         cfg,
		log:         log,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtManager:  jwtManager,
	}
}

func (as *authService) Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.Token, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	user, err := as.createNewUser(req)
	if err != nil {
		return nil, err
	}

	savedUser, err := as.userRepo.Save(ctx, user)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		as.log.Error("failed to save user into db",
			slog.String("error", err.Error()),
			slog.String("user_id", user.ID().String()),
		)
		return nil, consts.ErrInternalServer
	}

	access, refresh, err := as.generatePairOfTokens(savedUser)
	if err != nil {
		as.log.Warn("failed to generate a pair of tokens",
			slog.String("error", err.Error()),
		)
		return nil, consts.ErrInternalServer
	}

	session, err := as.createSession(ctx, savedUser.ID())
	if err != nil {
		return nil, err
	}

	if err := as.sessionRepo.Set(ctx, session, refresh); err != nil {
		as.log.Error("failed to save session into cache",
			slog.String("error", err.Error()),
		)
		return nil, consts.ErrInternalServer
	}

	return &dtos.Token{
		Access:     access,
		AccessTTL:  as.cfg.Cache.AccessTTL.Milliseconds(),
		Refresh:    refresh,
		RefreshTTL: as.cfg.Cache.AccessTTL.Milliseconds(),
	}, nil
}

func (as *authService) createSession(ctx context.Context, userID uuid.UUID) (*session.Session, error) {
	ip, err := as.getIPFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	userAgent, err := as.getUserAgentFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	hashPrint := sha256.Sum256([]byte(fmt.Sprintf("%s:%s", ip, userAgent)))
	fingerPrint := hex.EncodeToString(
		hashPrint[:],
	)

	return &session.Session{
		UserID:      userID,
		FingerPrint: fingerPrint,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func (as *authService) getIPFromCtx(ctx context.Context) (string, error) {
	raw := ctx.Value("x-client-ip")
	if ip, ok := raw.(string); ok {
		return ip, nil
	}

	return "", consts.ErrInvalidIP
}

func (as *authService) getUserAgentFromCtx(ctx context.Context) (string, error) {
	raw := ctx.Value("x-user-agent")
	if userAgent, ok := raw.(string); ok {
		return userAgent, nil
	}

	return "", consts.ErrInvalidUserAgent
}

// returns access and refresh
func (as *authService) generatePairOfTokens(user *user.User) (string, string, error) {
	access, err := as.jwtManager.GenerateAccess(user.ID())
	if err != nil {
		return "", "", err
	}

	refresh, err := as.jwtManager.GenerateRefresh()
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (as *authService) createNewUser(req *dtos.RegisterRequest) (*user.User, error) {
	pwdHash, err := user.NewPasswordHash(req.Password)
	if err != nil {
		return nil, err
	}

	return user.New(
		req.Firstname,
		req.Lastname,
		pwdHash,
	)
}
