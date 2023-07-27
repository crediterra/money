package money

import (
	"testing"

	"github.com/strongo/decimal"
)

func TestNewAmount(t *testing.T) {
	rub := CurrencyRUB
	amount := NewAmount(rub, decimal.NewDecimal64p2FromFloat64(123.45))
	if amount.Currency != rub {
		t.Error("amount.CurrencyCode != rub")
	}
	if amount.Value != decimal.NewDecimal64p2FromFloat64(123.45) {
		t.Error("amount.Value != 123.45")
	}
}

func TestAmount_Validate(t *testing.T) {
	if err := new(Amount).Validate(); err == nil {
		t.Error("new(Amount).Validate() should return error")
	}

	amount := NewAmount(CurrencyRUB, decimal.NewDecimal64p2FromFloat64(123.45))
	if err := amount.Validate(); err != nil {
		t.Error("amount.Validate() should not return error")
	}
}
