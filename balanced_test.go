package money

import (
	"reflect"
	"testing"
)

func TestBalancedBalanceFirestoreTag(t *testing.T) {
	field, ok := reflect.TypeOf(Balanced{}).FieldByName("Balance")
	if !ok {
		t.Fatal("Balanced.Balance field is missing")
	}
	if got, want := field.Tag.Get("firestore"), "balance,omitempty"; got != want {
		t.Errorf("Balanced.Balance firestore tag = %q, want %q", got, want)
	}
}

func TestBalanced_AddToBalance(t *testing.T) {
	b := &Balanced{}
	b.AddToBalance(CurrencyUSD, 100)
	if b.Balance[CurrencyUSD] != 100 {
		t.Errorf("expected 100, got %v", b.Balance[CurrencyUSD])
	}
}
