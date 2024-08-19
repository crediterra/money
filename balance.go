package money

import (
	"bytes"
	"github.com/strongo/decimal"
	"sort"
)

type TestStruct struct {
	ID   int64
	Name string
}

type Balance map[CurrencyCode]decimal.Decimal64p2

func (b Balance) IsZero() bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

func (b Balance) Reversed() (reversed Balance) {
	reversed = make(Balance, len(b))
	for currency, value := range b {
		reversed[currency] = -value
	}
	return
}

func (b Balance) Equal(b2 Balance) bool {
	if len(b) != len(b2) {
		return false
	}
	for c, v := range b {
		if v != b2[c] {
			return false
		}
	}
	for c, v := range b2 {
		if v != b[c] {
			return false
		}
	}
	return true
}

func (b Balance) OnlyPositive() Balance {
	result := make(Balance, len(b))
	for c, v := range b {
		if v > 0 {
			result[c] = v
		}
	}
	return result
}

func (b Balance) OnlyNegative() Balance {
	result := make(Balance, len(b))
	for c, v := range b {
		if v < 0 {
			result[c] = v
		}
	}
	return result
}

func (b Balance) CommaSeparatedUnsignedWithSymbols(translator interface{ Translate(s string) string }) string {
	lastIndex := len(b) - 1
	if lastIndex == 0 {
		for currency, value := range b {
			return Amount{Currency: currency, Value: value.Abs()}.String()
		}
	}
	var buffer bytes.Buffer
	i := 0
	sorter := &amountSorter{amounts: make([]Amount, len(b))}
	for currency, value := range b {
		amount := Amount{Currency: currency, Value: value.Abs()}
		sorter.amounts[i] = amount
		i += 1
	}

	sort.Sort(sorter)
	for i, amount := range sorter.amounts {
		buffer.WriteString(amount.String())
		switch {
		case i < lastIndex-1:
			buffer.WriteString(", ")
		case i == lastIndex-1:
			buffer.WriteString(translator.Translate(" and "))
		}
	}
	//log.Infof(c, "amounts: %v", buffer.String())
	return buffer.String()
}

func (b Balance) Add(amount Amount) Balance {
	//log.Debugf(c, "Balance.Add(amount=%v)", amount)
	if current, ok := b[amount.Currency]; ok {
		newVal := current + amount.Value
		//log.Debugf(c, "Balance.Add(): currency found: [%v], current=%v, newVal=%v", amount.CurrencyCode, current, newVal)
		if newVal == 0 {
			delete(b, amount.Currency)
		} else {
			b[amount.Currency] = newVal
		}
	} else {
		//log.Debugf(c, "Balance.Add(): currency NOT found: [%v], setting to: %v", amount.CurrencyCode, amount.Value)
		b[amount.Currency] = amount.Value
	}
	return b
}
