package main

import (
	"log"

	"github.com/ravenocx/hospital-mgt/config"
	"github.com/ravenocx/hospital-mgt/db"
	"github.com/ravenocx/hospital-mgt/server"
)

func main() {
	config, err := config.LoadConfig()
	
	if err != nil {
		log.Fatalf("failed to load config : %+v", err)
	}

	db, err := db.OpenConnection(config)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	s := server.NewServer(db, config)

	s.RegisterRoute()

	s.StarApp(config)
}