package main

import (
	"context"
	_ "embed"
	"log"

	"github.com/caleb-noodahl/bet-depot/bot/server"
	"github.com/caleb-noodahl/bet-depot/clients"
	"github.com/caleb-noodahl/bet-depot/config"
)

var s *server.Server

func init() {
	ctx := context.Background()
	c, err := config.ParseAPIConf()
	if err != nil {
		log.Panic(err)
	}

	disc, err := clients.NewDiscordClient(c)
	if err != nil {
		log.Panic(err)
	}

	bd := clients.NewBetDepotClient(c)

	gpt, _ := clients.NewGPTClient(c)
	s = server.NewServer(ctx, c, disc, gpt, bd)
}

func main() {
	if err := s.Start(); err != nil {
		log.Panic(err)
	}
}
