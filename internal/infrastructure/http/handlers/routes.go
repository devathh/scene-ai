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
