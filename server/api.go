package server

import (
	"github.com/gin-gonic/gin"

	"github.com/TremblingV5/TinyCache/base"
	"github.com/TremblingV5/TinyCache/internal/logger"
	"github.com/TremblingV5/TinyCache/internal/snowflake"
)

var config *base.Config

func Init(cfg *base.Config) {
	config = cfg
	snowflake.InitSnowFlake(int64(config.SnowFlakeNodeNum))
	logger.InitLogger(cfg)
}

func RunApiServer(port string) {
	r := gin.Default()

	group := r.Group("/api")
	group.GET("/:bucket/:key", GetController)
	group.POST("/:bucket/:key/:value", SetController)
	group.DELETE("/:bucket/:key", DelController)

	r.Use(logger.GinLogger())
	r.Use(logger.GinRecovery(true))
	r.Use(generateRequestId())

	r.Run("0.0.0.0:" + port)
}
