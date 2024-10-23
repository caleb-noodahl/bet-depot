package models

import (
	"time"

	"github.com/adhocore/gronx"
	"github.com/google/uuid"
)

type GameInstance struct {
	OwnerID   uuid.UUID `json:"owner_id" gorm:"type:uuid;not null"`
	Owner     User      `json:"owner" gorm:"foreignKey:OwnerID;"`
	GameID    uuid.UUID `json:"game_id" gorm:"type:uuid;not null"`
	Game      Game      `json:"game" gorm:"foreignKey:GameID;"`
	OutcomeID uuid.UUID `json:"outcome_id" gorm:"type:uuid;not null"`
}

type Game struct {
	StorageBase
	ShortID   string    `json:"short_id"`
	OwnerID   uuid.UUID `json:"owner_id" gorm:"type:uuid;not null"`
	Owner     User      `json:"owner" gorm:"foreignKey:OwnerID;"`
	Name      string    `json:"name"`
	StartCron string    `json:"start_cron"`
	EndCron   string    `json:"stop_cron"`
	NextStart time.Time `json:"next_start"`
	NextEnd   time.Time `json:"next_end"`
	Error     string    `json:"error"`
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
