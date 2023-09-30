package cache

import "github.com/gin-gonic/gin"

func RunCacheServer(port string) {
	r := gin.Default()

	group := r.Group("/cache")
	group.GET("/:bucket/:key", GetController)
	group.POST("/:bucket/:key/:value", SetController)
	group.DELETE("/:bucket/:key", DelController)
	group.POST("/:bucket", CreateBucketController)
	group.DELETE("/:bucket", RemoveBucketController)
	group.GET("/ping", PingController)

	r.Run("0.0.0.0:" + port)
}
