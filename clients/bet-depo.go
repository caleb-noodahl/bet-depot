package clients

import (
	"github.com/caleb-noodahl/bet-depot/config"
	"github.com/caleb-noodahl/bet-depot/server/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type BetDepotClient struct {
	conf   *config.APIConf
	client *resty.Client
}

func NewBetDepotClient(conf *config.APIConf) *BetDepotClient {
	return &BetDepotClient{
		conf:   conf,
		client: resty.New(),
	}
}

func (b *BetDepotClient) UpsertBook(ctx context.Context, book models.Book) (models.Book, error) {
	out := models.Book{}
	_, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(book).
		SetResult(&out).
		Post(fmt.Sprintf("%s/books", b.conf.BaseUrl))
	return out, err
}

func (b *BetDepotClient) CloseBook(ctx context.Context, closed models.ClosedBook) (models.ClosedBook, error) {
	out := models.ClosedBook{}
	bytes, _ := json.Marshal(closed)
	_, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytes).
		SetResult(&out).
		Post(fmt.Sprintf("%s/books/close", b.conf.BaseUrl))
	return out, err
}

func (b *BetDepotClient) GetBooks(ctx context.Context, book models.Book) ([]models.Book, error) {
	out := []models.Book{}
	_, err := b.client.R().
		SetQueryParams(map[string]string{
			"id":       book.ID.String(),
			"short_id": book.ShortID,
			"owner_id": book.OwnerID.String(),
		}).
		SetResult(&out).
		Get(fmt.Sprintf("%s/books", b.conf.BaseUrl))
	return out, err
}

func (b *BetDepotClient) GetTopBooks(ctx context.Context) ([]models.Book, error) {
	out := []models.Book{}
	_, err := b.client.R().
		SetResult(&out).
		Get(fmt.Sprintf("%s/books/top", b.conf.BaseUrl))
	return out, err
}

func (b *BetDepotClient) GetUser(ctx context.Context, user models.User) (models.User, error) {
	out := models.User{}
	_, err := b.client.R().
		SetQueryParams(map[string]string{
			"id":         user.ID.String(),
			"discord_id": user.DiscordID,
		}).
		SetResult(&out).
		Get(fmt.Sprintf("%s/user", b.conf.BaseUrl))
	return out, err
}

func (b *BetDepotClient) UpsertBet(ctx context.Context, bet models.Bet) (models.Bet, error) {
	out := models.Bet{}
	_, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bet).
		SetResult(&out).
		Post(fmt.Sprintf("%s/bets", b.conf.BaseUrl))
	return out, err
}

func (b *BetDepotClient) GetBets(ctx context.Context, bet models.Bet) ([]models.Bet, error) {
	out := []models.Bet{}
	_, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bet).
		SetResult(&out).
		Post(fmt.Sprintf("%s/bets", b.conf.BaseUrl))
	return out, err
}

func (b *BetDepotClient) UpsertOutcome(ctx context.Context, outcome models.Outcome) (models.Outcome, error) {
	out := models.Outcome{}
	_, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(outcome).
		SetResult(&out).
		Post(fmt.Sprintf("%s/outcomes", b.conf.BaseUrl))
	return out, err
}
