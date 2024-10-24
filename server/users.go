package server

import (
	"net/http"

	"github.com/caleb-noodahl/bet-depot/server/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

func (s *WebServer) UpsertUser(c echo.Context) error {
	ctx := c.Request().Context()
	req := models.User{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	createtx := s.db.Client.WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&req)
	if createtx.Error != nil {
		return c.JSON(http.StatusInternalServerError, createtx.Error)
	}

	gettx := s.db.Client.WithContext(ctx).
		Where(models.User{StorageBase: models.StorageBase{ID: req.ID}}).
		First(&req)
	if gettx.Error != nil {
		return c.JSON(http.StatusInternalServerError, gettx.Error)
	}

	return c.JSON(http.StatusOK, req)
}

func (s *WebServer) GetUser(c echo.Context) error {
	ctx := c.Request().Context()
	user, out := models.User{}, models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	result := s.db.Client.WithContext(ctx).
		First(&out, user)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}
	return c.JSON(http.StatusOK, out)
}

func (s *WebServer) GetWallet(c echo.Context) error {
	ctx := c.Request().Context()
	user, query := models.User{}, models.User{
		StorageBase: models.StorageBase{
			ID: parseUUIDParam(c.QueryParam("user_id")),
		},
		DiscordID: c.QueryParam("discord_id"),
		Username:  c.QueryParam("username"),
	}
	userctx := s.db.Client.WithContext(ctx).
		FirstOrCreate(&user, &query)
	if userctx.Error != nil {
		return c.JSON(http.StatusInternalServerError, userctx.Error)
	}
	wallet, wq := models.Wallet{}, models.Wallet{
		UserID: user.ID,
	}
	walletctx := s.db.Client.WithContext(ctx).
		Preload("User").
		FirstOrCreate(&wallet, &wq)
	if walletctx.Error != nil {
		return c.JSON(http.StatusInternalServerError, walletctx.Error)
	}
	return c.JSON(http.StatusOK, wallet)
}

func (s *WebServer) CreateTransaction(c echo.Context) error {
	ctx := c.Request().Context()
	tx := models.Transaction{}
	if err := c.Bind(&tx); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	wallet := models.Wallet{}
	walletctx := s.db.Client.WithContext(ctx).
		Preload("User").
		First(&wallet, &models.Wallet{StorageBase: models.StorageBase{ID: tx.WalletID}})
	if walletctx.Error != nil {
		return c.JSON(http.StatusInternalServerError, walletctx.Error)
	}
	wallet.AddTx(tx)

	walletctx = s.db.Client.WithContext(ctx).Save(wallet)
	if walletctx.Error != nil {
		return c.JSON(http.StatusInternalServerError, walletctx.Error)
	}

	return c.JSON(http.StatusOK, tx)
}
