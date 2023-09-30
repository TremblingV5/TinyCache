package main

import (
	"github.com/TremblingV5/TinyCache/base"
	"github.com/TremblingV5/TinyCache/server/api"
	"github.com/TremblingV5/TinyCache/server/cache"
)

var config *base.Config

func startAPIServer() {
	api.Init(config)
	api.RunApiServer("9999")
}

func startCacheServer() {
	cache.RunCacheServer("8001")
}

func main() {
	config = base.LoadConfig()

	go startAPIServer()

	startCacheServer()
}
