package money

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/strongo/decimal"
)

type dummyTranslator struct{}

func (dummyTranslator) Translate(s string) string {
	return s
}

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

	// Adding exact opposite removes currency from balance
	balance2.Add(NewAmount(rub, decimal.NewDecimal64p2FromFloat64(-124.12)))
	if len(balance2) != 0 {
		t.Errorf("expected empty balance after zeroing out, got %v", balance2)
	}
}

func TestBalance_IsZero(t *testing.T) {
	if !(Balance{}).IsZero() {
		t.Error("expected empty balance to be zero")
	}
	if !(Balance{CurrencyUSD: 0}).IsZero() {
		t.Error("expected zero value balance to be zero")
	}
	if (Balance{CurrencyUSD: 100}).IsZero() {
		t.Error("expected non-zero balance not to be zero")
	}
}

func TestBalance_Reversed(t *testing.T) {
	b := Balance{CurrencyUSD: 100, CurrencyEUR: -50}.Reversed()
	if b[CurrencyUSD] != -100 || b[CurrencyEUR] != 50 {
		t.Errorf("unexpected reversed balance: %v", b)
	}
}

func TestBalance_Equal(t *testing.T) {
	b1 := Balance{CurrencyUSD: 100, CurrencyEUR: 50}
	b2 := Balance{CurrencyUSD: 100, CurrencyEUR: 50}
	if !b1.Equal(b2) {
		t.Error("expected equal balances to be equal")
	}
	if b1.Equal(Balance{CurrencyUSD: 100}) {
		t.Error("expected different length balances not to be equal")
	}
	if b1.Equal(Balance{CurrencyUSD: 100, CurrencyEUR: 60}) {
		t.Error("expected different values not to be equal")
	}
	if b1.Equal(Balance{CurrencyUSD: 100, CurrencyGBP: 50}) {
		t.Error("expected different currencies not to be equal")
	}
}

func TestBalance_OnlyPositive_OnlyNegative(t *testing.T) {
	b := Balance{CurrencyUSD: 100, CurrencyEUR: -50}
	pos := b.OnlyPositive()
	if len(pos) != 1 || pos[CurrencyUSD] != 100 {
		t.Errorf("unexpected OnlyPositive: %v", pos)
	}
	neg := b.OnlyNegative()
	if len(neg) != 1 || neg[CurrencyEUR] != -50 {
		t.Errorf("unexpected OnlyNegative: %v", neg)
	}
}

func TestBalance_CommaSeparatedUnsignedWithSymbols(t *testing.T) {
	tr := dummyTranslator{}

	// Single currency
	single := Balance{CurrencyUSD: 100}
	if got := single.CommaSeparatedUnsignedWithSymbols(tr); got != "1 USD" {
		t.Errorf("single currency got %q, want '1 USD'", got)
	}

	// Multiple currencies (3 items)
	multi := Balance{
		CurrencyUSD: 300,
		CurrencyEUR: 200,
		CurrencyGBP: 100,
	}
	gotMulti := multi.CommaSeparatedUnsignedWithSymbols(tr)
	if !strings.Contains(gotMulti, " and ") || !strings.Contains(gotMulti, ", ") {
		t.Errorf("expected multi to contain ', ' and ' and ', got %q", gotMulti)
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
