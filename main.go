package main

import (
	"context"
	_ "embed"
	"log"

	"github.com/caleb-noodahl/bet-depot/config"
	"github.com/caleb-noodahl/bet-depot/database"
	"github.com/caleb-noodahl/bet-depot/server"
)

//go:embed config/api-config.yaml
var apiConfBytes []byte

func main() {
	ctx := context.Background()
	conf, err := config.ParseAPIConf(apiConfBytes)
	if err != nil {
		log.Panic(err)
	}
	db, err := database.NewPostgresDB(conf)
	if err != nil {
		log.Panic(err)
	}

	server := server.NewWebServer(ctx, conf, db)
	if err := server.Start(); err != nil {
		log.Panic(err)
	}
}
