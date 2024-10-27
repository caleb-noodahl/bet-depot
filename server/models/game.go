package models

import (
	"time"

	"github.com/adhocore/gronx"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type StockPickerGame struct {
	GameID  uuid.UUID      `json:"game_id" gorm:"type:uuid;not null"`
	Game    Game           `json:"game" gorm:"foreignKey:GameID"`
	BookID  uuid.UUID      `json:"book_id" gorm:"type:uuid;not null"`
	Book    Book           `json:"book" gorm:"foreignKey:BookID;"`
	Name    string         `json:"name"`
	Symbols pq.StringArray `json:"tags" gorm:"type:text[]"`
}

type Game struct {
	StorageBase
	ShortID   string    `json:"short_id"`
	OwnerID   uuid.UUID `json:"owner_id" gorm:"type:uuid;not null"`
	Owner     User      `json:"owner" gorm:"foreignKey:OwnerID;"`
	Name      string    `json:"name" gorm:"uniqueIndex"`
	StartCron string    `json:"start_cron"`
	EndCron   string    `json:"stop_cron"`
	NextStart time.Time `json:"next_start"`
	NextEnd   time.Time `json:"next_end"`
	Error     string    `json:"error"`
}

func (g *Game) BeforeCreate(tx *gorm.DB) (err error) {
	g.ID = lo.Ternary(g.ID == uuid.Nil, uuid.New(), g.ID)
	g.ShortID = g.ID.String()[:4]
	return
}

func (g *Game) Next() error {
	var err error
	g.NextStart, err = gronx.NextTick(g.StartCron, true)
	if err != nil {
		return err
	}
	g.NextEnd, err = gronx.NextTickAfter(g.EndCron, g.NextStart, true)
	if err != nil {
		return err
	}
	return nil
}
