package money

import (
	"github.com/strongo/decimal"
	"time"
)

// Balanced is a struct to store balance and summary of transfers
// ffjson: skip
type Balanced struct {
	Balance        Balance   `firestore:"balance,noindex,omitempty,noindex" json:"balance,omitempty"`
	LastTransferID string    `datastore:"lastTransferID,omitempty,noindex" json:"lastTransferID,omitempty"`
	LastTransferAt time.Time `datastore:"lastTransferAt,omitempty,noindex" json:"lastTransferAt,omitempty"`

	// Do not remove, needed to hide a balance/history menu in Telegram
	CountOfTransfers int `datastore:"countOfTransfers,omitempty" json:"countOfTransfers,omitempty"`
}

func (v *Balanced) AddToBalance(currency CurrencyCode, value decimal.Decimal64p2) {
	newBalance := v.Balance.Add(Amount{Currency: currency, Value: value})
	v.Balance = newBalance
}
