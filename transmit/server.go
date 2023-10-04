package transmit

import (
	"context"

	"github.com/TremblingV5/TinyCache/base"
	"github.com/TremblingV5/TinyCache/pb"
)

type CacheServer struct {
	pb.UnimplementedCacheServerServer
}

func (c *CacheServer) Get(ctx context.Context, req *pb.GetKeyRequest) (*pb.GetKeyResponse, error) {
	bucketName := req.GetBucket()
	key := req.GetKey()

	resp := &pb.GetKeyResponse{}

	if bucket := base.GetBucket(bucketName); bucket != nil {
		if value, err := bucket.GetLocally(key); err != nil {
			resp.ErrCode = base.ErrKeyNotFound.Code()
			return resp, err
		} else {
			resp.Value = value.B
			resp.ErrCode = base.Success.Code()
			return resp, nil
		}
	} else {
		resp.ErrCode = base.ErrBucketNotFound.Code()
		return resp, nil
	}
}

func (c *CacheServer) Set(ctx context.Context, req *pb.SetKeyRequest) (*pb.SetKeyResponse, error) {
	bucketName := req.GetBucket()
	key := req.GetKey()
	value := req.GetValue()

	if bucket := base.GetBucket(bucketName); bucket == nil {
		base.AddBucketLocally(bucketName, 20480)
	}

	bucket := base.GetBucket(bucketName)
	bucket.SetLocally(key, base.ByteView{
		B: []byte(value),
	})

	resp := &pb.SetKeyResponse{
		ErrCode: base.Success.Code(),
	}

	return resp, nil
}

func (c *CacheServer) Delete(ctx context.Context, req *pb.DeleteKeyRequest) (*pb.DeleteKeyResponse, error) {
	bucketName := req.GetBucket()
	key := req.GetKey()

	resp := &pb.DeleteKeyResponse{}

	if bucket := base.GetBucket(bucketName); bucket != nil {
		bucket.DelLocally(key)
		resp.ErrCode = base.Success.Code()
		return resp, nil
	} else {
		resp.ErrCode = base.ErrBucketNotFound.Code()
		return resp, nil
	}
}

func (c *CacheServer) DeleteBucket(ctx context.Context, req *pb.DeleteBucketRequest) (*pb.DeleteBucketResponse, error) {
	bucketName := req.GetBucket()

	base.RemoveBucketLocally(bucketName)

	resp := &pb.DeleteBucketResponse{
		ErrCode: base.Success.Code(),
	}
	return resp, nil
}
