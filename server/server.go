package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/caleb-noodahl/bet-depot/clients"
	"github.com/caleb-noodahl/bet-depot/config"
	"github.com/caleb-noodahl/bet-depot/database"
	"github.com/caleb-noodahl/bet-depot/server/models"
	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type WebServer struct {
	ctx       context.Context
	conf      *config.APIConf
	gamesconf *config.GamesConfig
	logger    *log.Logger
	client    *echo.Echo
	db        database.PostgresDB
	fin       *clients.FinnhubClient
}

func NewWebServer(ctx context.Context, config *config.APIConf, gamesconf *config.GamesConfig, fin *clients.FinnhubClient, db database.PostgresDB) *WebServer {
	e := echo.New()
	s := WebServer{
		ctx:       ctx,
		conf:      config,
		gamesconf: gamesconf,
		client:    e,
		fin:       fin,
		db:        db,
		logger:    log.New(""),
	}

	e.Logger.SetLevel(log.DEBUG)
	s.client.GET("/health", s.Health, s.LogAPIRequest)
	s.client.POST("/migrate", s.Migrate, s.LogAPIRequest)

	s.client.POST("/books", s.UpsertBook, s.LogAPIRequest)
	s.client.GET("/books", s.GetBooks, s.LogAPIRequest)
	s.client.GET("/books/top", s.GetTopBooks, s.LogAPIRequest)
	s.client.POST("/books/close", s.CloseBook, s.LogAPIRequest)

	s.client.GET("/user", s.GetUser, s.LogAPIRequest)
	s.client.POST("/users", s.UpsertUser, s.LogAPIRequest)
	s.client.GET("/wallet", s.GetWallet, s.LogAPIRequest)
	s.client.POST("/tx", s.CreateTransaction, s.LogAPIRequest)

	s.client.POST("/bets", s.UpsertBet, s.LogAPIRequest)
	s.client.POST("/outcomes", s.UpsertOutcome, s.LogAPIRequest)

	s.client.POST("/games/stockpicker", s.StockPickerGame, s.LogAPIRequest)
	return &s
}

func (s *WebServer) Start() error {
	return s.client.Start(fmt.Sprintf(":%v", s.conf.Port))
}

func (s *WebServer) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func (s *WebServer) Migrate(c echo.Context) error {
	mdls := []interface{}{
		&models.Outcome{},
		&models.Bet{},
		&models.Book{},
		&models.Payout{},
		&models.ClosedBook{},
		&models.User{},
		&models.Transaction{},
		&models.Game{},
		&models.StockPickerGame{},
		&models.Wallet{},
	}
	if err := s.db.MigrateDomainModels(mdls...); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//default data
	games := []*models.Game{
		{
			StorageBase: models.StorageBase{
				ID: uuid.MustParse("4735a6ee-3518-4bdc-a522-db2c588cd112"),
			},
			OwnerID:   uuid.MustParse(s.gamesconf.AdminID),
			Name:      "stock_daily",
			StartCron: "30 14 * * 1-5",
			EndCron:   "0 21 * * 1-5",
		}, {
			StorageBase: models.StorageBase{
				ID: uuid.MustParse("828714a9-4e5d-40cb-bf25-1eb6183f0514"),
			},
			OwnerID:   uuid.MustParse(s.gamesconf.AdminID),
			Name:      "stock_weekly",
			StartCron: "30 14 * * 1",
			EndCron:   "0 21 * * 5",
		},
	}
	gametx := s.db.Client.Save(games)
	if gametx.Error != nil {
		log.Error(gametx.Error)
	}
	return c.JSON(http.StatusNoContent, "")

}

func (s *WebServer) LogAPIRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s.logger.Printf("%+v", c.Request())
		return next(c)
	}
}
