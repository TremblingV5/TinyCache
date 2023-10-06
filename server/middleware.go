package server

import (
	"github.com/gin-gonic/gin"

	"github.com/TremblingV5/TinyCache/internal/snowflake"
)

func generateRequestId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("requestId", snowflake.GetSnowFlakeId())
	}
}
