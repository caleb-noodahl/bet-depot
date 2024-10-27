package server

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/caleb-noodahl/bet-depot/server/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *WebServer) StockPickerGame(c echo.Context) error {
	ctx := c.Request().Context()
	req := models.StockPickerGame{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	game := models.Game{}
	gametx := s.db.Client.WithContext(ctx).
		First(&game, &models.Game{StorageBase: models.StorageBase{ID: req.GameID}})

	if gametx.Error != nil {
		return c.JSON(http.StatusInternalServerError, gametx.Error)
	}
	// setup the schedule game run times
	if err := game.Next(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	bookID := uuid.New()
	outcomes := []models.Outcome{}
	bookDesc := "A time based stock game pitting random stock symbols between each other"
	for range 3 {
		symbol := s.gamesconf.Symbols[Rand(0, len(s.gamesconf.Symbols))]
		quote, err := s.fin.Price(symbol)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		req.Symbols = append(req.Symbols, symbol)
		outcomes = append(outcomes, models.Outcome{
			Description: fmt.Sprintf("symbol %s performs best", symbol),
			Odds:        float64(Rand(2, 14)),
			RefVal:      quote,
		})
		outcomes = append(outcomes, models.Outcome{
			Description: fmt.Sprintf("symbol %s performs worst", symbol),
			Odds:        float64(Rand(2, 14)),
			RefVal:      quote,
		})
		outcomes = append(outcomes, models.Outcome{
			Description: fmt.Sprintf("symbol %s achieves top half", symbol),
			Odds:        float64(Rand(1, 5)),
		})
		bookDesc += fmt.Sprintf("\n[symbol](https://finance.yahoo.com/quote/%s/) %s : price - $%v", symbol, symbol, quote)
	}
	req.GameID = game.ID
	req.Game = game
	req.Book = models.Book{
		StorageBase: models.StorageBase{
			ID: bookID,
		},
		Name:        req.Name,
		Description: bookDesc,
		OwnerID:     uuid.MustParse(s.gamesconf.AdminID),
		Open:        game.NextStart,
		LastCall:    game.NextEnd.Add(-4 * time.Hour),
		WagerType:   "standard",
		Options:     outcomes,
	}

	reqtx := s.db.Client.WithContext(ctx).
		Create(&req)
	if reqtx.Error != nil {
		return c.JSON(http.StatusInternalServerError, reqtx.Error)
	}

	return c.JSON(http.StatusOK, req)
}

func (s *WebServer) CloseStockPickerGame(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}

func Rand(lowest, max int) int {
	out, _ := rand.Int(rand.Reader, big.NewInt(int64(max)-int64(lowest)))
	out = big.NewInt(out.Int64() + int64(lowest))
	return int(out.Int64())
}
