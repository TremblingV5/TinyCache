package transmit

import (
	"context"
	tinycache "github.com/TremblingV5/TinyCache"

	"github.com/TremblingV5/TinyCache/base"
	"github.com/TremblingV5/TinyCache/pb"
)

type CacheServer struct{}

func (c *CacheServer) Get(ctx context.Context, req *pb.GetKeyRequest) (*pb.GetKeyResponse, error) {
	bucketName := req.GetBucket()
	key := req.GetKey()

	resp := &pb.GetKeyResponse{}

	if bucket := tinycache.GetBucket(bucketName); bucket != nil {
		if value, err := bucket.GetLocally(key); err != nil {
			resp.ErrCode = tinycache.ErrKeyNotFound.Code()
			return resp, err
		} else {
			resp.Value = value.B
			resp.ErrCode = tinycache.Success.Code()
			return resp, nil
		}
	} else {
		resp.ErrCode = tinycache.ErrBucketNotFound.Code()
		return resp, nil
	}
}

func (c *CacheServer) Set(ctx context.Context, req *pb.SetKeyRequest) (*pb.SetKeyResponse, error) {
	bucketName := req.GetBucket()
	key := req.GetKey()
	value := req.GetValue()

	if bucket := tinycache.GetBucket(bucketName); bucket == nil {
		tinycache.AddBucketLocally(bucketName, 20480)
	}

	bucket := tinycache.GetBucket(bucketName)
	bucket.SetLocally(key, base.ByteView{
		B: []byte(value),
	})

	resp := &pb.SetKeyResponse{
		ErrCode: tinycache.Success.Code(),
	}

	return resp, nil
}

func (c *CacheServer) Delete(ctx context.Context, req *pb.DeleteKeyRequest) (*pb.DeleteKeyResponse, error) {
	bucketName := req.GetBucket()
	key := req.GetKey()

	resp := &pb.DeleteKeyResponse{}

	if bucket := tinycache.GetBucket(bucketName); bucket != nil {
		bucket.DelLocally(key)
		resp.ErrCode = tinycache.Success.Code()
		return resp, nil
	} else {
		resp.ErrCode = tinycache.ErrBucketNotFound.Code()
		return resp, nil
	}
}

func (c *CacheServer) DeleteBucket(ctx context.Context, req *pb.DeleteBucketRequest) (*pb.DeleteBucketResponse, error) {
	bucketName := req.GetBucket()

	tinycache.RemoveBucketLocally(bucketName)

	resp := &pb.DeleteBucketResponse{
		ErrCode: tinycache.Success.Code(),
	}
	return resp, nil
}
