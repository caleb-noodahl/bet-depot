package models

import "github.com/google/uuid"

type User struct {
	StorageBase
	DiscordID     string `json:"discord_id" gorm:"uniqueIndex"`
	DiscordHandle string `json:"discord_handle"`
	Username      string `json:"username" gorm:"not null;`
	Email         string `json:"email"`
}

type Transaction struct {
	StorageBase
	WalletID uuid.UUID `json:"wallet_id" gorm:"uniqueIndex"`
	Amount   float64   `json:"amount"`
}

type Wallet struct {
	StorageBase
	UserID  uuid.UUID     `json:"user_id" gorm:"uniqueIndex"`
	User    User          `json:"user" gorm:"foreignKey:UserID"`
	Balance float64       `json:"balance"`
	Txs     []Transaction `json:"txs" gorm:"foreignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (w *Wallet) AddTx(tx Transaction) {
	w.Balance += tx.Amount
	w.Txs = append(w.Txs, tx)
}
