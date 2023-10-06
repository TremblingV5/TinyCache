package transmit

import (
	"context"

	"google.golang.org/grpc"

	"github.com/TremblingV5/TinyCache/ds/cmap"
	"github.com/TremblingV5/TinyCache/pb"
	"github.com/TremblingV5/TinyCache/singleflight"
)

var clientList = cmap.New[*RpcClient](32)

type RpcClient struct {
	client pb.CacheServerClient
	loader *singleflight.Group
}

func GetRpcClient(addr string) *RpcClient {
	if client, ok := clientList.Get(addr); ok {
		return client
	} else {
		conn, _ := grpc.Dial(addr, grpc.WithInsecure())
		client := &RpcClient{
			client: pb.NewCacheServerClient(conn),
			loader: &singleflight.Group{},
		}
		clientList.Set(addr, client)
		return client
	}
}

func (c *RpcClient) Get(bucket, key string) ([]byte, error) {
	v, err := c.loader.Do(key, func() (interface{}, error) {
		resp, err := c.client.Get(context.Background(), &pb.GetKeyRequest{
			Bucket: bucket,
			Key:    key,
		})
		if err != nil {
			return nil, err
		}
		return resp.Value, nil
	})
	return v.([]byte), err
}

func (c *RpcClient) Set(bucket, key string, value []byte) error {
	_, err := c.client.Set(context.Background(), &pb.SetKeyRequest{
		Bucket: bucket,
		Key:    key,
		Value:  value,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *RpcClient) Delete(bucket, key string) error {
	_, err := c.client.Delete(context.Background(), &pb.DeleteKeyRequest{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *RpcClient) DeleteBucket(bucket string) error {
	_, err := c.client.DeleteBucket(context.Background(), &pb.DeleteBucketRequest{
		Bucket: bucket,
	})
	if err != nil {
		return err
	}
	return nil
}
