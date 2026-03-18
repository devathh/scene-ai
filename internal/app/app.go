package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/devathh/scene-ai/internal/common/config"
	"github.com/devathh/scene-ai/internal/infrastructure/cache/redis"
	httpserver "github.com/devathh/scene-ai/internal/infrastructure/http"
	"github.com/devathh/scene-ai/internal/infrastructure/http/handlers"
	jwtkeyloader "github.com/devathh/scene-ai/internal/infrastructure/jwt/keyloader"
	jwtmanager "github.com/devathh/scene-ai/internal/infrastructure/jwt/manager"
	"github.com/devathh/scene-ai/internal/infrastructure/persistence/postgres"
	authservices "github.com/devathh/scene-ai/internal/modules/auth/application/services"
	sessionredis "github.com/devathh/scene-ai/internal/modules/auth/infrastructure/cache/redis/session"
	authuserpg "github.com/devathh/scene-ai/internal/modules/auth/infrastructure/persistence/postgres/user"
	"github.com/devathh/scene-ai/pkg/log"
	"github.com/joho/godotenv"
	redissdk "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	log    *slog.Logger
	server *httpserver.Server
}

func (a *App) Start() error {
	a.log.Info("server is running")
	return a.server.Start()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("server shutdown...")
	return a.server.Shutdown(ctx)
}

func New() (*App, func(), error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, nil, fmt.Errorf("failed to load .env: %w", err)
	}

	cfg, err := provideConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide config: %w", err)
	}

	log, err := provideLogger(cfg.App.Env)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide logger: %w", err)
	}

	log.Info("config was loaded", slog.Any("server", cfg.Server))

	db, err := providePersistence(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide persistence: %w", err)
	}

	redisClient, err := provideCache(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide cache: %w", err)
	}

	server, err := provideServer(log, cfg, db, redisClient)
	if err != nil {
		return nil, nil, err
	}

	return &App{
			log:    log,
			server: server,
		}, func() {
			if err := postgres.Close(db); err != nil {
				log.Error("failed to close connection with postges", slog.String("error", err.Error()))
			} else {
				log.Info("connection with postgres was closed")
			}

			if err := redis.Close(redisClient); err != nil {
				log.Error("failed to close connection redis", slog.String("error", err.Error()))
			} else {
				log.Info("connection with redis was closed")
			}
		}, nil
}

func provideServer(
	log *slog.Logger,
	cfg *config.Config,
	db *gorm.DB,
	redisClient *redissdk.Client,
) (*httpserver.Server, error) {
	authService, err := provideAuthService(
		log,
		cfg,
		db,
		redisClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to provide auth-service: %w", err)
	}

	handler, err := handlers.New(cfg, authService)
	if err != nil {
		return nil, fmt.Errorf("failed to create handler: %w", err)
	}

	return httpserver.New(cfg, handler), nil
}

func provideAuthService(
	log *slog.Logger,
	cfg *config.Config,
	db *gorm.DB,
	redisClient *redissdk.Client,
) (authservices.AuthService, error) {
	jwtManager, err := jwtmanager.New(cfg, jwtkeyloader.New(
		cfg.JWT.PublicKeyPath,
		cfg.JWT.PrivateKeyPath,
	))
	if err != nil {
		return nil, err
	}

	return authservices.New(
		cfg,
		log,
		authuserpg.New(db),
		sessionredis.New(cfg, redisClient),
		jwtManager,
	), nil
}

func provideCache(cfg *config.Config) (*redissdk.Client, error) {
	client, err := redis.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func providePersistence(cfg *config.Config) (*gorm.DB, error) {
	db, err := postgres.Connect(cfg)
	if err != nil {
		return nil, err
	}

	if err := postgres.Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func provideLogger(env string) (*slog.Logger, error) {
	logHandler, err := log.SetupHandler(os.Stdout, env)
	if err != nil {
		return nil, err
	}

	return slog.New(logHandler), nil
}

// .env must be loaded
func provideConfig() (*config.Config, error) {
	constructor := config.NewConstructor()
	cfg, err := constructor.Init(os.Getenv("APP_PATH_CONFIG"))
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
