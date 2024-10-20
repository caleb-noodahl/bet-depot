package server

import (
	"github.com/caleb-noodahl/bet-depot/clients"
	"github.com/caleb-noodahl/bet-depot/config"
	"context"
)

type Server struct {
	ctx  context.Context
	conf *config.APIConf
	disc *clients.DiscordClient
	gpt  *clients.GPTClient
	bd   *clients.BetDepotClient
}

func NewServer(ctx context.Context, conf *config.APIConf, disc *clients.DiscordClient, gpt *clients.GPTClient, bd *clients.BetDepotClient) *Server {
	server := Server{
		ctx:  ctx,
		conf: conf,
		disc: disc,
		gpt:  gpt,
		bd:   bd,
	}
	server.disc.Cmds["create"] = server.CreateBook
	server.disc.Cmds["view"] = server.ViewBook
	server.disc.Cmds["top"] = server.TopBooks
	server.disc.Cmds["close"] = server.CloseBook
	server.disc.Cmds["bet"] = server.CreateBet
	server.disc.Cmds["odds"] = server.UpdateOdds
	server.disc.Cmds["option"] = server.UpsertOutcome

	return &server
}

func (server *Server) Start() error {
	return server.disc.Start(server.ctx)
}
