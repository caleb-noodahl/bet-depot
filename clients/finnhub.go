package clients

import (
	"context"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/caleb-noodahl/bet-depot/config"
)

type FinnhubClient struct {
	ctx    context.Context
	conf   *config.APIConf
	client *finnhub.DefaultApiService
}

func NewFinnhubClient(ctx context.Context, conf *config.APIConf) (*FinnhubClient, error) {
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", conf.FinnhubApiKey)

	return &FinnhubClient{
		ctx:    ctx,
		conf:   conf,
		client: finnhub.NewAPIClient(cfg).DefaultApi,
	}, nil
}

func (f *FinnhubClient) Price(symbol string) (float64, error) {
	quote, _, err := f.client.Quote(f.ctx, symbol)
	return float64(quote.C), err
}
