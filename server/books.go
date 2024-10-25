package server

import (
	"net/http"

	"github.com/caleb-noodahl/bet-depot/server/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"gorm.io/gorm"
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
		FirstOrCreate(&req)
	if createtx.Error != nil {
		return c.JSON(http.StatusInternalServerError, createtx.Error)
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
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
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
	if len(topIDs) == 0 {
		return c.JSON(http.StatusOK, top)
	}
	books := []models.Book{}
	if len(topIDs) == 0 {
		return c.JSON(http.StatusOK, books)
	}

	if err := s.db.Client.WithContext(ctx).
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
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

	book := models.Book{}
	booktx := s.db.Client.WithContext(ctx).
		First(&book, &models.Book{StorageBase: models.StorageBase{ID: req.BookID}})
	if booktx.Error != nil {
		return c.JSON(http.StatusInternalServerError, booktx.Error)
	}
	book.Closed = true
	booktx = s.db.Client.WithContext(ctx).Save(&book)
	if booktx.Error != nil {
		return c.JSON(http.StatusInternalServerError, booktx.Error)
	}

	closedbooktx := s.db.Client.WithContext(ctx).
		Save(&req)
	if closedbooktx.Error != nil {
		return c.JSON(http.StatusInternalServerError, closedbooktx.Error)
	}

	for _, payout := range req.Payouts {
		wallet := models.Wallet{
			UserID: payout.UserID,
			Txs:    []models.Transaction{},
		}
		walletctx := s.db.Client.WithContext(ctx).
			Where(models.Wallet{UserID: payout.UserID}).
			Preload("Txs").
			FirstOrCreate(&wallet)
		if walletctx.Error != nil {
			return c.JSON(http.StatusInternalServerError, walletctx.Error)
		}
		wallet.AddTx(models.Transaction{
			Amount:     payout.Amount,
			WalletID:   wallet.ID,
			SourceID:   req.BookID,
			SourceType: "book_close",
		})
		walletctx = s.db.Client.WithContext(ctx).Save(wallet)
		if walletctx.Error != nil {
			return c.JSON(http.StatusInternalServerError, walletctx.Error)
		}
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
		Preload("Outcome").
		Save(&req)
	if createtx.Error != nil {
		return c.JSON(http.StatusInternalServerError, createtx.Error)
	}
	wallet, query := models.Wallet{}, models.Wallet{
		UserID: req.OwnerID,
		User:   req.Owner,
	}

	wallettx := s.db.Client.WithContext(ctx).
		Preload("Txs").
		FirstOrCreate(&wallet, &query)
	if wallettx.Error != nil {
		return c.JSON(http.StatusInternalServerError, wallettx.Error)
	}

	wallet.AddTx(models.Transaction{
		WalletID:   wallet.ID,
		SourceID:   req.StorageBase.ID,
		SourceType: "bet",
		Amount:     req.Amount * -1,
	})
	wallettx = s.db.Client.WithContext(ctx).Save(&wallet)
	if wallettx.Error != nil {
		return c.JSON(http.StatusInternalServerError, wallettx.Error)
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
		Save(&req)
	if createtx.Error != nil {
		return c.JSON(http.StatusInternalServerError, createtx.Error)
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
