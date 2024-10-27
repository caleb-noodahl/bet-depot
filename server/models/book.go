package models

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Bet struct {
	StorageBase
	OwnerID   uuid.UUID `json:"owner_id" gorm:"type:uuid;not null;"`
	Owner     User      `json:"owner" gorm:"foreignKey:OwnerID;"`
	BookID    uuid.UUID `json:"book_id" gorm:"type:uuid;not null;"`
	OutcomeID uuid.UUID `json:"outcome_id"`
	Outcome   Outcome   `json:"outcome" gorm:"foreignKey:OutcomeID"`
	Amount    float64   `json:"amount"`
}

type Outcome struct {
	StorageBase
	BookID      uuid.UUID `json:"book_id" gorm:"type:uuid;not null"`
	Description string    `json:"description"`
	Odds        float64   `json:"odds"`
	RefVal      float64   `json:"ref_val"`
}

type Payout struct {
	StorageBase
	BookID uuid.UUID `json:"book_id" gorm:"type:uuid;not null"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	BetID  uuid.UUID `json:"bet_id" gorm:"type:uuid;not null"`
	Amount float64   `json:"amount"`
}

type ClosedBook struct {
	StorageBase
	BookID      uuid.UUID `json:"book_id" gorm:"type:uuid;not null"`
	Book        Book      `json:"book" gorm:"foreignKey:BookID"`
	OutcomeID   uuid.UUID `json:"outcome_id" gorm:"type:uuid;not null"`
	Outcome     Outcome   `json:"outcome" gorm:"foreignKey:OutcomeID"`
	TotalPayout float64   `json:"total_payout"`
	Payouts     []Payout  `json:"payouts" gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Probability struct {
	Outcome uuid.UUID `json:"outcome_id"`
	Odds    float64   `json:"odds"`
	Count   int       `json:"count"`
}

type Book struct {
	StorageBase
	Name        string         `json:"name"`
	ShortID     string         `json:"short_id" gorm:"uniqueIndex"`
	Description string         `json:"description"`
	GameID      uuid.UUID      `json:"game_id"`
	Game        Game           `json:"game" gorm:"foreignKey:GameID"`
	WagerType   WagerType      `json:"wager_type"`
	Bets        Bets           `json:"bets" gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Options     []Outcome      `json:"options" gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OwnerID     uuid.UUID      `json:"owner_id"`
	Owner       User           `json:"owner" gorm:"foreignKey:OwnerID"`
	Closed      bool           `json:"closed"`
	Open        time.Time      `json:"open"`
	LastCall    time.Time      `json:"last_call"`
	MaxBet      float64        `json:"max_bet"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
}

func (b *Book) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	b.ShortID = b.ID.String()[:4]
	return
}

func (b *Book) CloseBook(outcomeIndex int) ClosedBook {
	closed := ClosedBook{
		BookID:    b.ID,
		Book:      *b,
		OutcomeID: b.Options[outcomeIndex].ID,
		Outcome:   b.Options[outcomeIndex],
		Payouts:   []Payout{},
	}

	total := 0.0
	for _, bet := range b.Bets {
		if bet.OutcomeID == closed.Outcome.ID {
			amount := bet.Amount * closed.Outcome.Odds
			closed.Payouts = append(closed.Payouts, Payout{
				BookID: b.ID,
				UserID: bet.OwnerID,
				BetID:  bet.ID,
				Amount: amount,
			})
			total += amount
		}
	}
	closed.TotalPayout = total
	return closed
}

type Bets []Bet

func (b Bets) ImpliedOdds() map[uuid.UUID]Probability {
	results := map[uuid.UUID]Probability{}
	total := 0.0
	grouped := lo.GroupBy(b, func(i Bet) uuid.UUID {
		total += i.Amount
		return i.OutcomeID
	})
	for key, val := range grouped {
		count := 0

		betTotal := lo.SumBy(val, func(i Bet) float64 {
			count++
			return i.Amount
		})
		// directly calculate implied probability as a percentage
		percent := math.Round((betTotal / total) * 100)
		// store the probability (percentage) directly
		results[key] = Probability{
			Outcome: key,
			Odds:    percent,
			Count:   count,
		}
	}
	return results
}
