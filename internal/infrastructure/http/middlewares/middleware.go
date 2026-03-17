package middlewares

import (
	"context"

	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/gin-gonic/gin"
)

func BaseMiddleware(ctx *gin.Context) {
	clientIP := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	baseCtx := ctx.Request.Context()
	enrichedCtx := context.WithValue(baseCtx, consts.ClientIP, clientIP)
	enrichedCtx = context.WithValue(enrichedCtx, consts.UserAgent, userAgent)

	ctx.Request = ctx.Request.WithContext(enrichedCtx)

	ctx.Next()
}
