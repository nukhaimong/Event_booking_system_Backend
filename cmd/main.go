package main

import (
	"gotickets/internal/config"
	"gotickets/internal/server"
)

func main() {
	// load environment variables
	cfg := config.LoadEnv()
	db := config.ConnectDB(cfg)

	// start the server
	server.Start(db, cfg)
}
