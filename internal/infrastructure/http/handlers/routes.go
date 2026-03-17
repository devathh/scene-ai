package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/devathh/scene-ai/internal/common/dtos"
	authservices "github.com/devathh/scene-ai/internal/modules/auth/application/services"
	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Routes struct {
	authService authservices.AuthService
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

func (r *Routes) Delete() gin.HandlerFunc {
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
		errors.Is(err, consts.ErrSessionNotFound):
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
	}

	return http.StatusInternalServerError
}