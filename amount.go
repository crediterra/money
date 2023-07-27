package money

import (
	"fmt"

	"github.com/strongo/decimal"
)

// Amount holds amount of money with currency
type Amount struct {
	Currency CurrencyCode
	Value    decimal.Decimal64p2
}

// NewAmount creates new amount
func NewAmount(currency CurrencyCode, value decimal.Decimal64p2) Amount {
	if currency == "" {
		panic("CurrencyCode not provided")
	}
	return Amount{
		Currency: currency,
		Value:    value,
	}
}

// Validate returns error if amount is invalid
func (a Amount) Validate() error {
	if a.Currency == "" {
		return fmt.Errorf("currency not provided")
	}
	if !IsKnownCurrency(a.Currency) {
		return fmt.Errorf("unknown currency: %v", a.Currency)
	}
	return nil
}

// IsZero returns true if amount is zero
func (a Amount) IsZero() bool {
	return a.Value == 0 && a.Currency == ""
}

// String returns string representation of amount
func (a Amount) String() string {
	//if currencySign, ok := currencySigns[v.CurrencyCode]; ok {
	//	return fmt.Sprintf("%v%v", currencySign, v.Value)
	//}
	return fmt.Sprintf("%v %v", a.Value, a.Currency)
}

type amountSorter struct {
	amounts []Amount
}

// Len is part of sort.Interface.
func (s *amountSorter) Len() int {
	return len(s.amounts)
}

// Swap is part of sort.Interface.
func (s *amountSorter) Swap(i, j int) {
	s.amounts[i], s.amounts[j] = s.amounts[j], s.amounts[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *amountSorter) Less(i, j int) bool {
	return s.amounts[i].Value > s.amounts[j].Value // Reverse sort - large amounts first
}
