package api

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
		if value, err := bucket.Get(key); err != nil {
			c.JSON(http.StatusOK, tinycache.ErrKeyNotFound)
		} else {
			c.JSON(http.StatusOK, tinycache.Success.WithData(value.String()))
		}
	} else {
		c.JSON(http.StatusOK, tinycache.ErrBucketNotFound)
	}
}

func SetController(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := c.Param("key")
	value := c.Param("value")

	if bucket := tinycache.GetBucket(bucketName); bucket == nil {
		tinycache.AddBucketLocally(bucketName, config.MaxBytes)
	}

	bucket := tinycache.GetBucket(bucketName)
	bucket.Set(key, base.ByteView{
		B: []byte(value),
	})
	c.JSON(http.StatusOK, tinycache.Success)
}

func DelController(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := c.Param("key")

	if bucket := tinycache.GetBucket(bucketName); bucket != nil {
		bucket.Del(key)
		c.JSON(http.StatusOK, tinycache.Success)
	} else {
		c.JSON(http.StatusOK, tinycache.ErrBucketNotFound)
	}
}

func RemoveBucketController(c *gin.Context) {
	bucketName := c.Param("bucket")
	tinycache.RemoveBucketLocally(bucketName)
	c.JSON(http.StatusOK, tinycache.Success)
}
