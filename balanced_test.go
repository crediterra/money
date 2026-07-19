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
