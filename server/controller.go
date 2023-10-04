package server

import (
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
			c.JSON(http.StatusOK, base.ErrKeyNotFound)
		} else {
			c.JSON(http.StatusOK, base.Success.WithData(value.String()))
		}
	} else {
		c.JSON(http.StatusOK, base.ErrBucketNotFound)
	}
}

func SetController(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := c.Param("key")
	value := c.Param("value")

	if bucket := tinycache.GetBucket(bucketName); bucket == nil {
		tinycache.AddBucketLocally(bucketName, config.MaxBytes)
	}

	tinycache.Set(bucketName, key, []byte(value))
	c.JSON(http.StatusOK, base.Success)
}

func DelController(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := c.Param("key")

	if bucket := tinycache.GetBucket(bucketName); bucket != nil {
		tinycache.Del(bucketName, key)
		c.JSON(http.StatusOK, base.Success)
	} else {
		c.JSON(http.StatusOK, base.ErrBucketNotFound)
	}
}
