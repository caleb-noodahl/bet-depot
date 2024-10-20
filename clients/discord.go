package clients

import (
	"github.com/caleb-noodahl/bet-depot/config"
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type DiscordClient struct {
	prefix  string
	conf    *config.APIConf
	Cmds    map[string]func(*discordgo.Session, *discordgo.MessageCreate) error
	Session *discordgo.Session
}

func (d *DiscordClient) Start(ctx context.Context) error {
	if err := d.Session.Open(); err != nil {
		return err
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc // Block until a signal is received
	return d.Session.Close()
}

func NewDiscordClient(conf *config.APIConf) (*DiscordClient, error) {
	session, err := discordgo.New("Bot " + conf.DiscordBotToken)
	if err != nil {
		return nil, err
	}

	client := &DiscordClient{
		conf:    conf,
		prefix:  "$",
		Session: session,
		Cmds:    map[string]func(*discordgo.Session, *discordgo.MessageCreate) error{},
	}
	session.AddHandler(client.MessageRouter)
	return client, nil
}

func (d *DiscordClient) MessageRouter(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, "$") {
		return
	}
	input := strings.Split(m.Content[1:], " ")
	log.Printf("-= cmd : %v", input)

	if val, ok := d.Cmds[input[0]]; ok {
		if err := val(s, m); err != nil {
			log.Printf(" != error: %s", err)
		}
	}
}
