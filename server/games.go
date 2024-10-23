package server

import (
	"net/http"

	"github.com/caleb-noodahl/bet-depot/server/models"
	"github.com/labstack/echo/v4"
)

type GameReq struct {
	ShortID string                 `json:"short_id"`
	Type    string                 `json:"type"`
	Args    map[string]interface{} `json:"args"`
}

func (s *WebServer) StockPickerGame(c echo.Context) error {
	ctx := c.Request().Context()
	req := GameReq{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	game := models.Game{}
	//err := getCreate(&game, models.Game{ShortID: req.ShortID}, s.db.Client)
	gettx := s.db.Client.WithContext(ctx).
		Where(models.Game{ShortID: req.ShortID}).
		First(&game)
	if gettx.Error != nil {
		return c.JSON(http.StatusInternalServerError, gettx.Error)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (s *WebServer) CloseStockPickerGame(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}
