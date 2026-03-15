package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/devathh/scene-ai/internal/common/config"
	"github.com/devathh/scene-ai/internal/common/persistence/postgres"
	authuserpg "github.com/devathh/scene-ai/internal/modules/auth/infrastructure/persistence/postgres/user"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	constr := config.NewConstructor()
	cfg, err := constr.Init(os.Getenv("APP_CONFIG_PATH"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	db, err := postgres.Connect(cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := postgres.Migrate(db); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	ur := authuserpg.New(db)
	fmt.Println(ur.Delete(context.Background(), uuid.MustParse("ad237575-7fb0-4025-a89d-8d46e015529f")))
}
