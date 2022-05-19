package main

import (
	"RexPromAgent/config"
	"RexPromAgent/log"
	"RexPromAgent/server"
)

func main() {
	config.Initialize()
	log.Initialize()
	newServer := server.NewServer()
	newServer.Start()

}
