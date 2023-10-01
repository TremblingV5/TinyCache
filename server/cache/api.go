package cache

import (
	"github.com/TremblingV5/TinyCache/base"
	"github.com/gin-gonic/gin"
)

var config *base.Config

func Init(cfg *base.Config) {
	config = cfg
}

func RunCacheServer(port string) {
	r := gin.Default()

	group := r.Group("/cache")
	group.GET("/:bucket/:key", GetController)
	group.POST("/:bucket/:key/:value", SetController)
	group.DELETE("/:bucket/:key", DelController)

	group.DELETE("/:bucket", RemoveBucketController)

	group.GET("/ping/:name/:site", PingController)

	r.Run("0.0.0.0:" + port)
}
