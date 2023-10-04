package main

import (
	"github.com/TremblingV5/TinyCache/base"
	"github.com/TremblingV5/TinyCache/server"
	"github.com/TremblingV5/TinyCache/transmit"
)

var config *base.Config

func startAPIServer() {
	if config.StartApi {
		server.Init(config)
		server.RunApiServer(config.ApiPort)
	}
}

func startCacheServer() {
	transmit.Init(config)
}

func main() {
	config = base.LoadConfig()

	go startAPIServer()
	startCacheServer()
}
