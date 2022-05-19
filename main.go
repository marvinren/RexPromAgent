package main

import (
	"RexPromAgent/pkg/config"
	"RexPromAgent/pkg/log"
	"RexPromAgent/pkg/server"
)

func main() {
	config.Initialize()
	log.Initialize()
	newServer := server.NewServer()
	newServer.Start()

}
