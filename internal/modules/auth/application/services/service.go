package authservices

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
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
	Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.Token, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dtos.User, error)
	GetByToken(ctx context.Context, token string) (*dtos.User, error)
	UpdateUser(ctx context.Context, token string, req *dtos.UpdateUserRequest) (*dtos.User, error)
	Delete(ctx context.Context, token string) error
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

	u, err := as.createNewUser(req)
	if err != nil {
		return nil, err
	}

	savedUser, err := as.userRepo.Save(ctx, u)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, consts.ErrUserAlreadyExists) {
			return nil, err
		}

		as.log.Error("failed to save user into db",
			slog.String("error", err.Error()),
			slog.String("user_id", u.ID().String()),
		)
		return nil, consts.ErrInternalServer
	}

	access, refresh, err := as.approveTokens(ctx, savedUser)
	if err != nil {
		return nil, err
	}

	return &dtos.Token{
		Access:     access,
		AccessTTL:  as.cfg.JWT.AccessTTL.Milliseconds(),
		Refresh:    refresh,
		RefreshTTL: as.cfg.JWT.RefreshTTL.Milliseconds(),
	}, nil
}

func (as *authService) Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.Token, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	email := user.Email(strings.TrimSpace(req.Email))
	if !email.IsValid() {
		return nil, consts.ErrInvalidEmail
	}

	u, err := as.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, consts.ErrUserNotFound) {
			return nil, consts.ErrInvalidCredentials
		}

		as.log.Error("failed to get user by email",
			slog.String("error", err.Error()),
			slog.String("email", email.String()),
		)
		return nil, consts.ErrInternalServer
	}

	if !u.PasswordHash().Compare(req.Password) {
		return nil, consts.ErrInvalidCredentials
	}

	access, refresh, err := as.approveTokens(ctx, u)
	if err != nil {
		return nil, err
	}

	return &dtos.Token{
		Access:     access,
		AccessTTL:  as.cfg.JWT.AccessTTL.Milliseconds(),
		Refresh:    refresh,
		RefreshTTL: as.cfg.JWT.RefreshTTL.Milliseconds(),
	}, nil
}

func (as *authService) GetByID(ctx context.Context, id uuid.UUID) (*dtos.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	u, err := as.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, consts.ErrUserNotFound) {
			return nil, err
		}

		as.log.Error("failed to get user by id",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return nil, consts.ErrInternalServer
	}

	return &dtos.User{
		ID:        u.ID().String(),
		Email:     u.Email().String(),
		Firstname: u.Firstname(),
		Lastname:  u.Lastname(),
		CreatedAt: u.CreatedAt().UnixMilli(),
		UpdatedAt: u.UpdatedAt().UnixMilli(),
	}, nil
}

func (as *authService) GetByToken(ctx context.Context, token string) (*dtos.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	claims, err := as.jwtManager.Validate(token)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, consts.ErrInvalidToken
	}

	u, err := as.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, consts.ErrUserNotFound) {
			return nil, err
		}

		as.log.Error("failed to get user by id",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return nil, consts.ErrInternalServer
	}

	return &dtos.User{
		ID:        u.ID().String(),
		Email:     u.Email().String(),
		Firstname: u.Firstname(),
		Lastname:  u.Lastname(),
		CreatedAt: u.CreatedAt().UnixMilli(),
		UpdatedAt: u.UpdatedAt().UnixMilli(),
	}, nil
}

func (as *authService) UpdateUser(ctx context.Context, token string, req *dtos.UpdateUserRequest) (*dtos.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	claims, err := as.jwtManager.Validate(token)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, consts.ErrInvalidToken
	}

	u, err := as.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return nil, err
		}
		if errors.Is(err, consts.ErrUserNotFound) {
			return nil, err
		}
		as.log.Error("failed to get user for update",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, consts.ErrInternalServer
	}

	mask := user.UpdateMask{}

	if req.Email != nil {
		email := user.Email(*req.Email)
		mask.Email = &email
	}
	if req.Firstname != nil {
		mask.Firstname = req.Firstname
	}
	if req.Lastname != nil {
		mask.Lastname = req.Lastname
	}
	if req.Password != nil {
		mask.Password = req.Password
	}

	if err := u.Apply(mask); err != nil {
		return nil, err
	}

	updatedUser, err := as.userRepo.Update(ctx, u)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return nil, err
		}
		if errors.Is(err, consts.ErrUserAlreadyExists) {
			return nil, err
		}
		as.log.Error("failed to update user in db",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, consts.ErrInternalServer
	}

	return &dtos.User{
		ID:        updatedUser.ID().String(),
		Email:     updatedUser.Email().String(),
		Firstname: updatedUser.Firstname(),
		Lastname:  updatedUser.Lastname(),
		CreatedAt: updatedUser.CreatedAt().UnixMilli(),
		UpdatedAt: updatedUser.UpdatedAt().UnixMilli(),
	}, nil
}

func (as *authService) Delete(ctx context.Context, token string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	claims, err := as.jwtManager.Validate(token)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(claims.ID)
	if err != nil {
		return consts.ErrInvalidToken
	}

	if err := as.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return err
		}

		if errors.Is(err, consts.ErrUserNotFound) {
			return err
		}

		return err
	}

	return nil
}

func (as *authService) approveTokens(ctx context.Context, u *user.User) (string, string, error) {
	access, refresh, err := as.generatePairOfTokens(u)
	if err != nil {
		as.log.Warn("failed to generate a pair of tokens",
			slog.String("error", err.Error()),
		)
		return "", "", consts.ErrInternalServer
	}

	session, err := as.createSession(ctx, u.ID())
	if err != nil {
		return "", "", err
	}

	if err := as.sessionRepo.Set(ctx, session, refresh); err != nil {
		as.log.Error("failed to save session into cache",
			slog.String("error", err.Error()),
		)
		return "", "", consts.ErrInternalServer
	}

	return access, refresh, nil
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

	hashPrint := sha256.Sum256(fmt.Appendf(nil, "%s:%s", ip, userAgent))
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
	raw := ctx.Value(consts.ClientIP)
	if ip, ok := raw.(string); ok {
		return ip, nil
	}

	return "", consts.ErrInvalidIP
}

func (as *authService) getUserAgentFromCtx(ctx context.Context) (string, error) {
	raw := ctx.Value(consts.UserAgent)
	if userAgent, ok := raw.(string); ok {
		return userAgent, nil
	}

	return "", consts.ErrInvalidUserAgent
}

// returns access and refresh
func (as *authService) generatePairOfTokens(u *user.User) (string, string, error) {
	access, err := as.jwtManager.GenerateAccess(u.ID())
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
		user.Email(req.Email),
		req.Firstname,
		req.Lastname,
		pwdHash,
	)
}
