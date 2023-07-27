package money

import "testing"

func TestIsKnownCurrency(t *testing.T) {
	if IsKnownCurrency(CurrencyCode("")) {
		t.Error("Empty currency should not be known")
	}
	for _, c := range currencies {
		if !IsKnownCurrency(c) {
			t.Errorf("%s should be known currency", c)
		}
	}
}
