package server

import (
	"github.com/TremblingV5/TinyCache/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"

	tinycache "github.com/TremblingV5/TinyCache"
	"github.com/TremblingV5/TinyCache/base"
)

func GetController(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := c.Param("key")

	if bucket := tinycache.GetBucket(bucketName); bucket != nil {
		if value, err := tinycache.Get(bucketName, key); err != nil {
			logger.WriteLogWithCtx(c, "error", base.ErrKeyNotFound.Error(), "bucket", bucketName, "key", key)
			c.JSON(http.StatusOK, base.ErrKeyNotFound)
		} else {
			logger.WriteLogWithCtx(c, "info", "hit cache", "bucket", bucketName, "key", key)
			c.JSON(http.StatusOK, base.Success.WithData(value.String()))
		}
	} else {
		logger.WriteLogWithCtx(c, "error", base.ErrBucketNotFound.Error(), "bucket", bucketName)
		c.JSON(http.StatusOK, base.ErrBucketNotFound)
	}
}

func SetController(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := c.Param("key")
	value := c.Param("value")

	if bucket := tinycache.GetBucket(bucketName); bucket == nil {
		logger.WriteLogWithCtx(c, "info", "add bucket", "bucket", bucketName)
		tinycache.AddBucketLocally(bucketName, config.MaxBytes)
	}

	tinycache.Set(bucketName, key, []byte(value))
	logger.WriteLogWithCtx(c, "info", "set cache", "bucket", bucketName, "key", key, "value", value)
	c.JSON(http.StatusOK, base.Success)
}

func DelController(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := c.Param("key")

	if bucket := tinycache.GetBucket(bucketName); bucket != nil {
		tinycache.Del(bucketName, key)
		logger.WriteLogWithCtx(c, "info", "delete cache", "bucket", bucketName, "key", key)
		c.JSON(http.StatusOK, base.Success)
	} else {
		logger.WriteLogWithCtx(c, "error", base.ErrBucketNotFound.Error(), "bucket", bucketName)
		c.JSON(http.StatusOK, base.ErrBucketNotFound)
	}
}
