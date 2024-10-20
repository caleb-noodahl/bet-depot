package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/caleb-noodahl/bet-depot/server/models"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func (server *Server) CreateBook(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Join(strings.Split(m.Content, " ")[1:], " ")
	response, err := server.gpt.Prompt(BetPrompt, content)
	if err != nil {
		return err
	}
	user, err := server.bd.GetUser(server.ctx, models.User{DiscordID: m.Author.ID})
	if err != nil || user.ID == uuid.Nil {
		user = models.User{
			Username:      m.Author.Username,
			DiscordHandle: m.Author.Username,
			DiscordID:     m.Author.ID,
		}
	}

	book := models.Book{
		OwnerID: user.ID,
		Owner:   user,
	}
	if err := json.Unmarshal([]byte(response), &book); err != nil {
		return err
	}
	upsert, err := server.bd.UpsertBook(server.ctx, book)
	if err != nil {
		return err
	}
	msg, err := ParseTemplateResponse("book-create", upsert)
	if err != nil {
		return err
	}
	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}

func (server *Server) ViewBook(s *discordgo.Session, m *discordgo.MessageCreate) error {
	split := strings.Split(m.Content, " ")
	books, err := server.bd.GetBooks(server.ctx, models.Book{
		ShortID: strings.Replace(split[1], "#", "", -1),
	})
	if err != nil || len(books) < 1 {
		return fmt.Errorf("error fetching book")
	}
	book := books[0]
	data := map[string]interface{}{
		"Book":        book,
		"ImpliedOdds": book.Bets.ImpliedOdds(),
	}
	msg, err := ParseTemplateResponse("book-view", data)
	if err != nil {
		return err
	}

	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}

func (server *Server) CloseBook(s *discordgo.Session, m *discordgo.MessageCreate) error {
	split := strings.Split(m.Content, " ")
	if len(split) < 3 {
		return errors.New("invalid number of params to close")
	}
	option, err := strconv.Atoi(split[2])
	if err != nil {
		return err
	}
	books, err := server.bd.GetBooks(server.ctx, models.Book{
		ShortID: strings.Replace(split[1], "#", "", -1),
	})
	if err != nil || len(books) < 1 {
		return fmt.Errorf("error fetching book")
	}
	closed := books[0].CloseBook(option)
	upsert, err := server.bd.CloseBook(server.ctx, closed)
	if err != nil {
		return err
	}

	books[0].Closed = true
	_, err = server.bd.UpsertBook(server.ctx, books[0])
	if err != nil {
		return err
	}

	msg, err := ParseTemplateResponse("book-close", upsert)
	if err != nil {
		return err
	}

	s.ChannelMessageSend(m.ChannelID, msg)
	return err
}

func (server *Server) CreateBet(s *discordgo.Session, m *discordgo.MessageCreate) error {
	split := strings.Split(m.Content, " ")
	//cmd $bet #book:short_id option_num amount
	if len(split) < 3 {
		return errors.New("invalid number of params to create bet")
	}

	option, err := strconv.Atoi(strings.Replace(split[2], "#", "", -1))
	if err != nil {
		return err
	}
	amount, err := strconv.Atoi(strings.Replace(split[3], "$", "", -1))
	if err != nil {
		return err
	}
	user, err := server.bd.GetUser(server.ctx, models.User{DiscordID: m.Author.ID})
	if err != nil || user.ID == uuid.Nil {
		user = models.User{
			Username:      m.Author.Username,
			DiscordHandle: m.Author.Username,
			DiscordID:     m.Author.ID,
		}
	}

	books, err := server.bd.GetBooks(server.ctx, models.Book{ShortID: strings.Replace(split[1], "#", "", -1)})
	if err != nil || len(books) < 1 {
		return fmt.Errorf("unable to fetch book by short code: %s", split[1])
	}
	bet := models.Bet{
		BookID:    books[0].ID,
		OutcomeID: books[0].Options[option].ID,
		Outcome:   books[0].Options[option],
		Amount:    float64(amount),
		OwnerID:   user.ID,
		Owner:     user,
	}
	upsert, err := server.bd.UpsertBet(server.ctx, bet)
	if err != nil {
		return err
	}
	msg, err := ParseTemplateResponse("bet-create", upsert)
	if err != nil {
		return err
	}
	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}

func (server *Server) UpdateOdds(s *discordgo.Session, m *discordgo.MessageCreate) error {
	split := strings.Split(m.Content, " ")
	books, err := server.bd.GetBooks(server.ctx, models.Book{
		ShortID: strings.Replace(split[1], "#", "", -1),
	})
	if err != nil || len(books) < 1 {
		return fmt.Errorf("error fetching book")
	}
	book := books[0]
	option, err := strconv.Atoi(split[2])
	if err != nil {
		return err
	}

	odds, err := strconv.ParseFloat(split[3], 64)
	if err != nil {
		return err
	}
	outcome := book.Options[option]
	outcome.Odds = odds

	data, err := server.bd.UpsertOutcome(server.ctx, outcome)
	if err != nil {
		return err
	}

	msg, err := ParseTemplateResponse("odds-upsert", data)
	if err != nil {
		return err
	}

	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}

func (server *Server) TopBooks(s *discordgo.Session, m *discordgo.MessageCreate) error {
	books, err := server.bd.GetTopBooks(server.ctx)
	if err != nil {
		return err
	}

	msg, err := ParseTemplateResponse("books-top", books)
	if err != nil {
		return err
	}

	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}

func (server *Server) UpsertOutcome(s *discordgo.Session, m *discordgo.MessageCreate) error {
	split := strings.Split(m.Content, " ")
	books, err := server.bd.GetBooks(server.ctx, models.Book{
		ShortID: strings.Replace(split[1], "#", "", -1),
	})
	if err != nil || len(books) < 1 {
		return fmt.Errorf("error fetching book")
	}
	book := books[0]
	upsert, err := server.bd.UpsertOutcome(server.ctx, models.Outcome{
		BookID:      book.ID,
		Description: strings.Join(split[2:], " "),
	})
	if err != nil || len(books) < 1 {
		return err
	}

	data := map[string]interface{}{
		"Book":    book,
		"Outcome": upsert,
	}

	msg, err := ParseTemplateResponse("outcome-create", data)
	if err != nil {
		return err
	}
	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}
