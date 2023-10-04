package transmit

import (
	"context"
	"github.com/TremblingV5/TinyCache/ds/cmap"
	"github.com/TremblingV5/TinyCache/pb"
	"google.golang.org/grpc"
)

var clientList = cmap.New[*RpcClient](32)

type RpcClient struct {
	client pb.CacheServerClient
}

func GetRpcClient(addr string) *RpcClient {
	if client, ok := clientList.Get(addr); ok {
		return client
	} else {
		conn, _ := grpc.Dial(addr, grpc.WithInsecure())
		client := &RpcClient{
			client: pb.NewCacheServerClient(conn),
		}
		clientList.Set(addr, client)
		return client
	}
}

func (c *RpcClient) Get(bucket, key string) ([]byte, error) {
	resp, err := c.client.Get(context.Background(), &pb.GetKeyRequest{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
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
