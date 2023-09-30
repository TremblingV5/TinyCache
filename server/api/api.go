package api

import (
	"github.com/gin-gonic/gin"

	"github.com/TremblingV5/TinyCache/base"
)

var config *base.Config

func Init(cfg *base.Config) {
	config = cfg
	config.MaxBytes = 20480
}

func RunApiServer(port string) {
	r := gin.Default()

	group := r.Group("/api")
	group.GET("/:bucket/:key", GetController)
	group.POST("/:bucket/:key/:value", SetController)
	group.DELETE("/:bucket/:key", DelController)

	r.Run("0.0.0.0:" + port)
}
