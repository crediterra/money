package money

import (
	"strings"
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

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on empty currency")
		}
	}()
	NewAmount("", 10)
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

func TestValidateExplainsCaseMistakes(t *testing.T) {
	err := Amount{Currency: "usd", Value: 100}.Validate()
	if err == nil {
		t.Fatal("lowercase currency must be rejected")
	}
	for _, want := range []string{"uppercase ISO 4217", `"usd"`, `"USD"`, "adapter"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q missing %q", err.Error(), want)
		}
	}
	// Genuinely unknown codes keep the plain message — no bogus suggestion.
	if err := (Amount{Currency: "XYZ", Value: 100}).Validate(); err == nil || strings.Contains(err.Error(), "uppercase") {
		t.Errorf("unknown currency should not claim a case problem: %v", err)
	}
}

func TestAmount_IsZero(t *testing.T) {
	if !((Amount{}).IsZero()) {
		t.Error("expected zero amount to return true")
	}
	if (Amount{Currency: CurrencyUSD, Value: 100}).IsZero() {
		t.Error("expected non-zero amount to return false")
	}
}

func TestAmount_String(t *testing.T) {
	a := Amount{Currency: CurrencyUSD, Value: 100}
	if got := a.String(); got != "1 USD" {
		t.Errorf("got %q, want '1 USD'", got)
	}
}

func TestAmountSorter(t *testing.T) {
	sorter := &amountSorter{
		amounts: []Amount{
			{Value: 100},
			{Value: 200},
		},
	}
	if sorter.Len() != 2 {
		t.Errorf("expected len 2, got %d", sorter.Len())
	}
	if sorter.Less(0, 1) { // 100 > 200 is false, Less sorts descending
		t.Error("expected Less(0, 1) to be false")
	}
	sorter.Swap(0, 1)
	if sorter.amounts[0].Value != 200 {
		t.Errorf("expected 200 after swap, got %v", sorter.amounts[0].Value)
	}
}
