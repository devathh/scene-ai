package handlers

import (
	"errors"
	"net/http"

	"github.com/devathh/scene-ai/internal/common/config"
	"github.com/devathh/scene-ai/internal/infrastructure/http/middlewares"
	authservices "github.com/devathh/scene-ai/internal/modules/auth/application/services"
	"github.com/gin-gonic/gin"
)

var ErrInvalidEnvironment = errors.New("invalid global environment")

func New(
	cfg *config.Config,
	authService authservices.AuthService,
) (http.Handler, error) {
	var router *gin.Engine
	switch cfg.App.Env {
	case "prod":
		router = gin.New()
		router.Use(gin.Recovery())
	case "dev", "local":
		router = gin.Default()
	default:
		return nil, ErrInvalidEnvironment
	}

	router.Use(middlewares.BaseMiddleware)

	routes := Routes{
		authService: authService,
	}

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/register", routes.Register())
				auth.POST("/login", routes.Login())

				auth.PUT("/user", routes.UpdateUser())
				auth.GET("/user", routes.GetUser())
				auth.DELETE("/user", routes.Delete())
			}
		}
	}

	return router, nil
}
