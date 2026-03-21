package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devathh/scene-ai/internal/app"
	_ "github.com/devathh/scene-ai/docs"
)

// @title Scene AI API
// @version 1.0
// @description API for generating scenarios and scenes using Artificial Intelligence.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@scene-ai.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	appInstance, cleanup, err := app.New()
	if err != nil {
		slog.Error("failed to setup app", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer cleanup()

	go func() {
		if err := appInstance.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", slog.String("error", err.Error()))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if err := appInstance.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
}