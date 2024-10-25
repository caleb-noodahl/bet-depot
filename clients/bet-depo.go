package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/caleb-noodahl/bet-depot/config"
	"github.com/caleb-noodahl/bet-depot/server/models"

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
	resp, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(book).
		SetResult(&out).
		Post(fmt.Sprintf("%s/books", b.conf.BaseUrl))
	if resp.IsError() {
		return out, fmt.Errorf("upsert book failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return out, err
}

func (b *BetDepotClient) CloseBook(ctx context.Context, closed models.ClosedBook) (models.ClosedBook, error) {
	out := models.ClosedBook{}
	bytes, _ := json.Marshal(closed)
	resp, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytes).
		SetResult(&out).
		Post(fmt.Sprintf("%s/books/close", b.conf.BaseUrl))
	if resp.IsError() {
		return out, fmt.Errorf("close book failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return out, err
}

func (b *BetDepotClient) GetBooks(ctx context.Context, book models.Book) ([]models.Book, error) {
	out := []models.Book{}
	resp, err := b.client.R().
		SetQueryParams(map[string]string{
			"id":       book.ID.String(),
			"short_id": book.ShortID,
			"owner_id": book.OwnerID.String(),
		}).
		SetResult(&out).
		Get(fmt.Sprintf("%s/books", b.conf.BaseUrl))
	if resp.IsError() {
		return out, fmt.Errorf("get books failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return out, err
}

func (b *BetDepotClient) GetTopBooks(ctx context.Context) ([]models.Book, error) {
	out := []models.Book{}
	resp, err := b.client.R().
		SetResult(&out).
		Get(fmt.Sprintf("%s/books/top", b.conf.BaseUrl))
	if resp.IsError() {
		return out, fmt.Errorf("get top books failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return out, err
}

func (b *BetDepotClient) GetUser(ctx context.Context, user models.User) (models.User, error) {
	out := models.User{}
	resp, err := b.client.R().
		SetQueryParams(map[string]string{
			"id":         user.ID.String(),
			"discord_id": user.DiscordID,
		}).
		SetResult(&out).
		Get(fmt.Sprintf("%s/user", b.conf.BaseUrl))
	if resp.IsError() {
		return out, fmt.Errorf("get user failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return out, err
}

func (b *BetDepotClient) UpsertBet(ctx context.Context, bet models.Bet) (models.Bet, error) {
	out := models.Bet{}
	resp, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bet).
		SetResult(&out).
		Post(fmt.Sprintf("%s/bets", b.conf.BaseUrl))
	if resp.IsError() {
		return out, fmt.Errorf("upsert bet failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return out, err
}

func (b *BetDepotClient) GetBets(ctx context.Context, bet models.Bet) ([]models.Bet, error) {
	out := []models.Bet{}
	resp, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bet).
		SetResult(&out).
		Post(fmt.Sprintf("%s/bets", b.conf.BaseUrl))
	if resp.IsError() {
		return out, fmt.Errorf("get bets failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return out, err
}

func (b *BetDepotClient) UpsertOutcome(ctx context.Context, outcome models.Outcome) (models.Outcome, error) {
	out := models.Outcome{}
	resp, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(outcome).
		SetResult(&out).
		Post(fmt.Sprintf("%s/outcomes", b.conf.BaseUrl))
	if resp.IsError() {
		return out, fmt.Errorf("upsert outcome failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return out, err
}

func (b *BetDepotClient) GetWallet(ctx context.Context, wallet models.Wallet) (models.Wallet, error) {
	out := models.Wallet{}
	resp, err := b.client.R().
		SetQueryParams(map[string]string{
			"id":         wallet.ID.String(),
			"user_id":    wallet.User.ID.String(),
			"discord_id": wallet.User.DiscordID,
			"username":   wallet.User.Username,
		}).
		SetResult(&out).
		Get(fmt.Sprintf("%s/wallet", b.conf.BaseUrl))
	if resp.IsError() {
		return out, fmt.Errorf("get wallet failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return out, err
}

func (b *BetDepotClient) CreateTx(ctx context.Context, tx models.Transaction) (models.Transaction, error) {
	out := models.Transaction{}
	resp, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(tx).
		SetResult(&out).
		Post(fmt.Sprintf("%s/tx", b.conf.BaseUrl))
	if resp.IsError() {
		return out, fmt.Errorf("create tx failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return out, err
}
