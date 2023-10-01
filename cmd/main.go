package main

import (
	"github.com/TremblingV5/TinyCache/base"
	"github.com/TremblingV5/TinyCache/server/api"
	"github.com/TremblingV5/TinyCache/server/cache"
)

var config *base.Config

func startAPIServer() {
	api.Init(config)
	api.RunApiServer(config.ApiPort)
}

func startCacheServer() {
	cache.Init(config)
	cache.RunCacheServer(config.Port)
}

func main() {
	config = base.LoadConfig()

	go startAPIServer()

	startCacheServer()
}
