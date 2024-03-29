package money

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/strongo/decimal"
)

func TestBalance_Add(t *testing.T) {
	balance := make(Balance)
	rub := CurrencyCode(CurrencyRUB)
	balance2 := balance.Add(NewAmount(rub, decimal.NewDecimal64p2FromFloat64(123.45)))
	if balance.IsZero() {
		t.Error("balance.IsZero()")
	}
	if balance2.IsZero() {
		t.Error("balance2.IsZero()")
	}
	if len(balance2) != 1 {
		t.Error("len(balance2) != 1")
	}
	if v, ok := balance2[rub]; !ok {
		t.Error("balance2[rub] => !ok")
	} else if v != decimal.NewDecimal64p2FromFloat64(123.45) {
		t.Error("balance2[rub] != 123.45")
	}
	balance2.Add(NewAmount(rub, decimal.NewDecimal64p2FromFloat64(0.67)))
	if len(balance2) != 1 {
		t.Error("len(balance2) != 1")
	}
	if v, ok := balance2[rub]; !ok {
		t.Error("balance2[rub] => !ok")
	} else if v != decimal.NewDecimal64p2FromFloat64(124.12) {
		t.Error("balance2[rub] != 124.12")
	}
}

func TestBalance_ffjson(t *testing.T) {
	balance1 := Balance{
		CurrencyEUR: decimal.NewDecimal64p2(10, 2),
		CurrencyRUB: decimal.NewDecimal64p2(100, 0),
	}

	serialized, err := json.Marshal(balance1)
	if err != nil {
		t.Errorf("Failed to marshal: %v", err)
	}
	s := string(serialized)
	if !strings.Contains(s, `"EUR":10.02`) {
		t.Errorf("Missing correct EUR value, got: %v", s)
	}
}
