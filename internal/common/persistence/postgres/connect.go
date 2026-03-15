package postgres

import (
	"fmt"

	"github.com/devathh/scene-ai/internal/common/config"
	authuserpg "github.com/devathh/scene-ai/internal/modules/auth/infrastructure/persistence/postgres/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d sslmode=%s user=%s password=%s dbname=%s",
		cfg.Persistence.Postgres.Host,
		cfg.Persistence.Postgres.Port,
		cfg.Persistence.Postgres.SSLMode,
		cfg.Persistence.Postgres.Auth.User,
		cfg.Persistence.Postgres.Auth.Password,
		cfg.Persistence.Postgres.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to open connection with postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to load sql db: %w", err)
	}

	sqlDB.SetConnMaxIdleTime(cfg.Persistence.Postgres.Conn.MaxIdleTime)
	sqlDB.SetConnMaxLifetime(cfg.Persistence.Postgres.Conn.MaxLifetime)
	sqlDB.SetMaxIdleConns(cfg.Persistence.Postgres.Conn.MaxIdles)
	sqlDB.SetMaxOpenConns(cfg.Persistence.Postgres.Conn.MaxOpens)

	return db, nil
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&authuserpg.UserModel{}); err != nil {
		return fmt.Errorf("failed to migrate postgres: %w", err)
	}

	return nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql db: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close connection with postgres: %w", err)
	}

	return nil
}
