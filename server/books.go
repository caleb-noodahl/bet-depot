package server

import (
	"github.com/caleb-noodahl/bet-depot/server/models"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"gorm.io/gorm/clause"
)

func (s *WebServer) UpsertBook(c echo.Context) error {
	ctx := c.Request().Context()
	req := models.Book{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	req.SetDefaults()

	createtx := s.db.Client.WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&req)
	if createtx.Error != nil {
		return c.JSON(http.StatusInternalServerError, createtx.Error)
	}
	gettx := s.db.Client.WithContext(ctx).
		Where(models.Book{StorageBase: models.StorageBase{ID: req.ID}}).
		First(&req)
	if gettx.Error != nil {
		return c.JSON(http.StatusInternalServerError, gettx.Error)
	}
	return c.JSON(http.StatusOK, req)
}

func (s *WebServer) GetBooks(c echo.Context) error {
	ctx := c.Request().Context()

	books := []models.Book{}
	req := models.Book{
		StorageBase: models.StorageBase{ID: parseUUIDParam(c.QueryParam("id"))},
		ShortID:     c.QueryParam("short_id"),
	}

	result := s.db.Client.WithContext(ctx).
		Preload("Options").
		Preload("Bets").
		Order("created_at").
		Preload("Owner").
		Find(&books, &req)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}
	return c.JSON(http.StatusOK, books)
}

type topBooks struct {
	ID       string `json:"id"`
	BetCount int    `json:"bet_count"`
}

func (s *WebServer) GetTopBooks(c echo.Context) error {
	ctx := c.Request().Context()
	top := []topBooks{}
	if err := s.db.Client.WithContext(ctx).
		Table("books").
		Select("books.id, COUNT(bets.id) as bet_count").
		Joins("JOIN bets ON bets.book_id = books.id").
		Where("books.closed = ?", false).
		Group("books.id").
		Limit(3).
		Order("bet_count DESC").
		Scan(&top).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	topIDs := lo.Map(top, func(t topBooks, _ int) uuid.UUID {
		return uuid.MustParse(t.ID)
	})
	books := []models.Book{}
	if err := s.db.Client.WithContext(ctx).
		Preload("Options").
		Preload("Bets").
		Find(&books, topIDs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, books)
}

func (s *WebServer) CloseBook(c echo.Context) error {
	ctx := c.Request().Context()
	req := models.ClosedBook{}
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
		Where(models.Book{StorageBase: models.StorageBase{ID: req.ID}}).
		First(&req)
	if gettx.Error != nil {
		return c.JSON(http.StatusInternalServerError, gettx.Error)
	}
	return c.JSON(http.StatusOK, req)
}

func (s *WebServer) UpsertBet(c echo.Context) error {
	ctx := c.Request().Context()
	req := models.Bet{}
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
		Where(models.Bet{StorageBase: models.StorageBase{ID: req.ID}}).
		Preload("Outcome").
		First(&req)
	if gettx.Error != nil {
		return c.JSON(http.StatusInternalServerError, gettx.Error)
	}
	return c.JSON(http.StatusOK, req)
}

func (s *WebServer) UpsertOutcome(c echo.Context) error {
	ctx := c.Request().Context()
	req := models.Outcome{}
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
		Where(models.Outcome{StorageBase: models.StorageBase{ID: req.ID}}).
		First(&req)
	if gettx.Error != nil {
		return c.JSON(http.StatusInternalServerError, gettx.Error)
	}
	return c.JSON(http.StatusOK, req)
}

func (s *WebServer) GetBets(c echo.Context) error {
	ctx := c.Request().Context()
	bets := []models.Bet{}
	req := models.Bet{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	result := s.db.Client.WithContext(ctx).
		Preload("Owner").
		Find(&bets, &req)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}
	return c.JSON(http.StatusOK, bets)
}

func parseUUIDParam(id string) uuid.UUID {
	out, _ := uuid.Parse(id)
	return out
}
