package main

import (
	"context"
	"embed"
	"log"

	"github.com/caleb-noodahl/bet-depot/clients"
	"github.com/caleb-noodahl/bet-depot/config"
	"github.com/caleb-noodahl/bet-depot/database"
	"github.com/caleb-noodahl/bet-depot/server"
)

//go:embed config/data/*
var data embed.FS

func main() {
	ctx := context.Background()
	confBytes, err := data.ReadFile("config/data/api-config-dev.yaml")
	if err != nil {
		log.Panic(err)
	}
	conf, err := config.ParseAPIConf(confBytes)
	if err != nil {
		log.Panic(err)
	}
	gcb, err := data.ReadFile("config/data/games-config.yaml")
	if err != nil {
		log.Panic(err)
	}

	gamesconf, err := config.ParseGamesConf(gcb)
	if err != nil {
		log.Panic(err)
	}

	fin, err := clients.NewFinnhubClient(ctx, conf)
	if err != nil {
		log.Panic(err)
	}

	db, err := database.NewPostgresDB(conf)
	if err != nil {
		log.Panic(err)
	}

	server := server.NewWebServer(ctx, conf, gamesconf, fin, db)
	if err := server.Start(); err != nil {
		log.Panic(err)
	}
}
