package models

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type User struct {
	StorageBase
	DiscordID     string `json:"discord_id" gorm:"uniqueIndex"`
	DiscordHandle string `json:"discord_handle"`
	Username      string `json:"username" gorm:"not null;"`
	Email         string `json:"email"`
}

type Transaction struct {
	StorageBase
	WalletID   uuid.UUID `json:"wallet_id"`
	SourceID   uuid.UUID `json:"source_id"`
	SourceType string    `json:"source_type"`
	Amount     float64   `json:"amount"`
}

type Wallet struct {
	StorageBase
	UserID  uuid.UUID     `json:"user_id" gorm:"uniqueIndex"`
	User    User          `json:"user" gorm:"foreignKey:UserID"`
	Balance float64       `json:"balance"`
	Txs     []Transaction `json:"txs" gorm:"foreignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (w *Wallet) AddTx(tx Transaction) {
	w.Txs = append(w.Txs, tx)
	w.Balance = lo.SumBy(w.Txs, func(tx Transaction) float64 {
		return tx.Amount
	})
}
