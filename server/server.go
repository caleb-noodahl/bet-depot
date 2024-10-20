package server

import (
	"github.com/caleb-noodahl/bet-depot/config"
	"github.com/caleb-noodahl/bet-depot/database"
	"github.com/caleb-noodahl/bet-depot/server/models"
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type WebServer struct {
	ctx    context.Context
	conf   *config.APIConf
	logger *log.Logger
	client *echo.Echo
	db     database.PostgresDB
}

func NewWebServer(ctx context.Context, config *config.APIConf, db database.PostgresDB) *WebServer {
	e := echo.New()
	s := WebServer{
		ctx:    ctx,
		conf:   config,
		client: e,
		db:     db,
		logger: log.New(""),
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

	s.client.POST("/bets", s.UpsertBet, s.LogAPIRequest)

	s.client.POST("/outcomes", s.UpsertOutcome, s.LogAPIRequest)

	return &s
}

func (s *WebServer) Start() error {
	return s.client.Start(fmt.Sprintf(":%v", s.conf.Port))
}

func (s *WebServer) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func (s *WebServer) Migrate(c echo.Context) error {
	models := []interface{}{
		&models.Outcome{},
		&models.Bet{},
		&models.Book{},
		&models.Payout{},
		&models.ClosedBook{},
		&models.User{},
		&models.Transaction{},
		&models.Wallet{},
	}
	if err := s.db.MigrateDomainModels(models...); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusNoContent, "")

}

func (s *WebServer) LogAPIRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s.logger.Printf("%+v", c.Request())
		return next(c)
	}
}
