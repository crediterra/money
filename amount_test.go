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
