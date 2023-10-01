package transmit

import (
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
