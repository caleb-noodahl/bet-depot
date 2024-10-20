package server

import (
	"github.com/caleb-noodahl/bet-depot/server/models"
	"net/http"

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
