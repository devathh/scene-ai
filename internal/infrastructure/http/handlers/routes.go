package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/devathh/scene-ai/internal/common/dtos"
	aiservices "github.com/devathh/scene-ai/internal/modules/ai/application/services"
	authservices "github.com/devathh/scene-ai/internal/modules/auth/application/services"
	scenarioservices "github.com/devathh/scene-ai/internal/modules/scenario/application/services"
	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Routes struct {
	authService     authservices.AuthService
	scenarioService scenarioservices.ScenarioService
	aiService       aiservices.AIService
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.RegisterRequest true "Registration details"
// @Success 201 {object} dtos.User "User created successfully"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Failure 409 {object} dtos.ErrMsg "User already exists"
// @Router /api/v1/auth/register [post]
func (r *Routes) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.RegisterRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := r.authService.Register(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return access token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginRequest true "Login credentials"
// @Success 200 {object} dtos.Token "Login successful"
// @Failure 401 {object} dtos.ErrMsg "Invalid credentials"
// @Router /api/v1/auth/login [post]
func (r *Routes) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.LoginRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := r.authService.Login(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// UpdateUser godoc
// @Summary Update user profile
// @Description Update current user's information.
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dtos.UpdateUserRequest true "Update details"
// @Success 200 {object} dtos.User "Updated user"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Router /api/v1/auth/user [patch]
func (r *Routes) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := r.getToken(ctx)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		var req dtos.UpdateUserRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := r.authService.UpdateUser(ctx.Request.Context(), token, &req)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// GetUser godoc
// @Summary Get user details
// @Description Retrieve user profile by ID or current token.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Param id query string false "User ID (optional, defaults to token owner)"
// @Success 200 {object} dtos.User "User details"
// @Failure 400 {object} dtos.ErrMsg "Invalid ID"
// @Failure 404 {object} dtos.ErrMsg "User not found"
// @Router /api/v1/auth/user [get]
func (r *Routes) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			resp *dtos.User
			err  error
		)
		id := strings.TrimSpace(ctx.Query("id"))
		if id == "" {
			resp, err = r.getByToken(ctx)
		} else {
			resp, err = r.getByID(ctx, id)
		}

		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// DeleteUser godoc
// @Summary Delete user account
// @Description Permanently delete the current user's account.
// @Tags Auth
// @Security BearerAuth
// @Success 204 "No content"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Router /api/v1/auth/user [delete]
func (r *Routes) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := r.getToken(ctx)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		if err := r.authService.Delete(ctx, token); err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

// CreateScenario godoc
// @Summary Create a new scenario
// @Description Create a manual scenario with scenes.
// @Tags Scenario
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dtos.CreateScenarioRequest true "Scenario details"
// @Success 201 {object} dtos.Scenario "Scenario created"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Router /api/v1/scenario/scenario [post]
func (r *Routes) CreateScenario() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.CreateScenarioRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		token, err := r.getToken(ctx)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		resp, err := r.scenarioService.Create(ctx, &req, token)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	}
}

// UpdateScenario godoc
// @Summary Update a scenario
// @Description Update an existing scenario by ID.
// @Tags Scenario
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Scenario ID"
// @Param request body dtos.UpdateScenarioRequest true "Update details"
// @Success 200 {object} dtos.Scenario "Updated scenario"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Failure 404 {object} dtos.ErrMsg "Scenario not found"
// @Router /api/v1/scenario/scenario/{id} [patch]
func (r *Routes) UpdateScenario() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		var req dtos.UpdateScenarioRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		token, err := r.getToken(ctx)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		resp, err := r.scenarioService.Update(ctx, id, &req, token)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// DeleteScenario godoc
// @Summary Delete a scenario
// @Description Delete a scenario by ID.
// @Tags Scenario
// @Security BearerAuth
// @Param id path string true "Scenario ID"
// @Success 204 "No content"
// @Failure 404 {object} dtos.ErrMsg "Scenario not found"
// @Router /api/v1/scenario/scenario/{id} [delete]
func (r *Routes) DeleteScenario() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		token, err := r.getToken(ctx)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		if err := r.scenarioService.Delete(ctx, id, token); err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

// GetScenarioByID godoc
// @Summary Get scenario by ID
// @Description Retrieve a specific scenario by its ID.
// @Tags Scenario
// @Produce json
// @Param id path string true "Scenario ID"
// @Success 200 {object} dtos.Scenario "Scenario details"
// @Failure 400 {object} dtos.ErrMsg "Invalid ID"
// @Failure 404 {object} dtos.ErrMsg "Scenario not found"
// @Router /api/v1/scenario/scenario/{id} [get]
func (r *Routes) GetScenarioByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		resp, err := r.scenarioService.GetByID(ctx, id)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// GetScenarios godoc
// @Summary List scenarios
// @Description Retrieve a list of scenarios with pagination.
// @Tags Scenario
// @Produce json
// @Security BearerAuth
// @Param before-id query string false "UUID of the last item from previous page"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} dtos.Scenario "List of scenarios"
// @Failure 400 {object} dtos.ErrMsg "Invalid parameters"
// @Router /api/v1/scenario/scenarios [get]
func (r *Routes) GetScenarios() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := r.getToken(ctx)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		beforeID := uuid.Nil
		if raw := ctx.Query("before-id"); raw != "" {
			id, err := uuid.Parse(raw)
			if err != nil {
				ctx.JSON(
					http.StatusBadRequest,
					dtos.ErrMsg{
						Error: "invalid before id",
					},
				)
				return
			}

			beforeID = id
		}

		limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				dtos.ErrMsg{
					Error: "invalid limit",
				},
			)
			return
		}

		resp, err := r.scenarioService.GetList(ctx, beforeID, limit, token)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// GenerateScenario godoc
// @Summary Generate AI scenario
// @Description Trigger AI to generate a new scenario based on a prompt.
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dtos.GenerateScenarioRequest true "Generation prompt"
// @Success 202 {object} map[string]string "Job accepted, returns scenario ID"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Router /api/v1/ai/scenario [post]
func (r *Routes) GenerateScenario() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.GenerateScenarioRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		token, err := r.getToken(ctx)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		id, err := r.aiService.GenerateScenario(ctx, &req, token)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusAccepted, gin.H{
			"id": id.String(),
		})
	}
}

// GetScenes godoc
// @Summary Get generated scenes
// @Description Retrieve the list of generated scenes for a scenario.
// @Tags AI
// @Produce json
// @Param id path string true "Scenario ID"
// @Success 200 {array} dtos.Scene "List of scenes"
// @Failure 400 {object} dtos.ErrMsg "Invalid ID"
// @Router /api/v1/ai/scenario/{id}/scenes [get]
func (r *Routes) GetScenes() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		resp, err := r.aiService.GetScenes(ctx, id)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// GetScenario godoc
// @Summary Get AI scenario status
// @Description Retrieve the status and details of an AI-generated scenario.
// @Tags AI
// @Produce json
// @Param id path string true "Scenario ID"
// @Success 200 {object} dtos.Scenario "Scenario details"
// @Failure 400 {object} dtos.ErrMsg "Invalid ID"
// @Router /api/v1/ai/scenario/{id} [get]
func (r *Routes) GetScenario() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		resp, err := r.aiService.GetScenario(ctx, id)
		if err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// Connect godoc
// @Summary Connect to scene stream
// @Description Establish a WebSocket connection to stream scene generation progress.
// @Tags AI
// @Param id path string true "Scene ID"
// @Success 101 "Switching Protocols"
// @Failure 400 {object} dtos.ErrMsg "Invalid ID"
// @Router /api/v1/ai/scenes/{id} [get]
func (r *Routes) Connect() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dtos.ErrMsg{
				Error: "failed to connect",
			})
			return
		}
		defer conn.Close()

		if err := r.aiService.Connect(ctx.Request.Context(), conn, id); err != nil {
			ctx.JSON(
				r.getCode(err),
				dtos.ErrMsg{
					Error: err.Error(),
				},
			)
			return
		}
	}
}

func (r *Routes) getByID(ctx *gin.Context, raw string) (*dtos.User, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return nil, consts.ErrInvalidUserID
	}

	return r.authService.GetByID(ctx.Request.Context(), id)
}

func (r *Routes) getByToken(ctx *gin.Context) (*dtos.User, error) {
	token, err := r.getToken(ctx)
	if err != nil {
		return nil, err
	}

	return r.authService.GetByToken(ctx.Request.Context(), token)
}

func (r *Routes) getToken(ctx *gin.Context) (string, error) {
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		return "", consts.ErrInvalidToken
	}

	return authorization, nil
}

func (r *Routes) getCode(err error) int {
	switch {
	case errors.Is(err, consts.ErrInvalidFirstname),
		errors.Is(err, consts.ErrInvalidLastname),
		errors.Is(err, consts.ErrInvalidPassword),
		errors.Is(err, consts.ErrInvalidUserID):
		return http.StatusBadRequest
	case errors.Is(err, consts.ErrUserNotFound),
		errors.Is(err, consts.ErrSessionNotFound),
		errors.Is(err, consts.ErrScenarioNotFound):
		return http.StatusNotFound
	case errors.Is(err, consts.ErrInvalidToken):
		return http.StatusUnauthorized
	case errors.Is(err, consts.ErrInvalidIP),
		errors.Is(err, consts.ErrInvalidUserAgent):
		return http.StatusBadRequest
	case errors.Is(err, consts.ErrUserAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, consts.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, consts.ErrEmptyTitle),
		errors.Is(err, consts.ErrEmptyScenarioPrompt),
		errors.Is(err, consts.ErrEmptyGlobalStylePrompt),
		errors.Is(err, consts.ErrEmptyVideoPrompt),
		errors.Is(err, consts.ErrNoScenes),
		errors.Is(err, consts.ErrInvalidLimit):
		return http.StatusBadRequest
	case errors.Is(err, consts.ErrForbidden):
		return http.StatusForbidden
	}

	return http.StatusInternalServerError
}
