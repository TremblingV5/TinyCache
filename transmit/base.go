package transmit

import (
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/TremblingV5/TinyCache/base"
	"github.com/TremblingV5/TinyCache/pb"
)

var config *base.Config
var pool *ServerPool

func Init(cfg *base.Config) {
	config = cfg

	InitCachePool()
	StartCacheServer()
}

func StartCacheServer() {
	listen, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterCacheServerServer(server, &CacheServer{})

	reflection.Register(server)
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}

func InitCachePool() {
	pool = NewServerPool(config.SecondaryNum)

	serverList := strings.Split(config.SecondaryList, ",")
	pool.Add(serverList...)
}

func Get(bucket, key string) ([]byte, error) {
	addr := pool.Pick(key)
	client := GetRpcClient(addr)
	return client.Get(bucket, key)
}

func Set(bucket, key string, value []byte) error {
	addr := pool.Pick(key)
	client := GetRpcClient(addr)
	return client.Set(bucket, key, value)
}

func Delete(bucket, key string) error {
	addr := pool.Pick(key)
	client := GetRpcClient(addr)
	return client.Delete(bucket, key)
}
