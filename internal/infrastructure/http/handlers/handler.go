package handlers

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/devathh/scene-ai/internal/common/config"
	"github.com/devathh/scene-ai/internal/infrastructure/http/middlewares"
	aiservices "github.com/devathh/scene-ai/internal/modules/ai/application/services"
	authservices "github.com/devathh/scene-ai/internal/modules/auth/application/services"
	scenarioservices "github.com/devathh/scene-ai/internal/modules/scenario/application/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var ErrInvalidEnvironment = errors.New("invalid global environment")

func New(
	cfg *config.Config,
	authService authservices.AuthService,
	scenarioService scenarioservices.ScenarioService,
	aiService aiservices.AIService,
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	frontendPath := "./frontend"
	router.Static("/assets", filepath.Join(frontendPath, "assets"))
	
	router.GET("/generating.html", func(c *gin.Context) {
		c.File(filepath.Join(frontendPath, "generating.html"))
	})
	router.GET("/profile.html", func(c *gin.Context) {
		c.File(filepath.Join(frontendPath, "profile.html"))
	})

	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		if len(path) >= 8 && path[:8] == "/swagger" {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		if path == "/generating.html" || path == "/profile.html" {
			c.File(filepath.Join(frontendPath, path[1:]))
			return
		}

		c.File(filepath.Join(frontendPath, "index.html"))
	})

	router.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(frontendPath, "index.html"))
	})

	routes := Routes{
		authService:     authService,
		scenarioService: scenarioService,
		aiService:       aiService,
	}

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/register", routes.Register())
				auth.POST("/login", routes.Login())

				auth.PATCH("/user", routes.UpdateUser())
				auth.GET("/user", routes.GetUser())
				auth.DELETE("/user", routes.DeleteUser())
			}

			scenario := v1.Group("/scenario")
			{
				scenario.POST("/scenario", routes.CreateScenario())
				scenario.PATCH("/scenario/:id", routes.UpdateScenario())
				scenario.DELETE("/scenario/:id", routes.DeleteScenario())
				scenario.GET("/scenario/:id", routes.GetScenarioByID())

				scenario.GET("/scenarios", routes.GetScenarios())
			}

			ai := v1.Group("/ai")
			{
				ai.POST("/scenario", routes.GenerateScenario())
				ai.GET("/scenario/:id", routes.GetScenario())
				ai.GET("/scenario/:id/scenes", routes.GetScenes())
				ai.GET("/scenes/:id", routes.Connect())
			}
		}
	}

	return router, nil
}