package money

import (
	"bytes"
	"encoding/json"
	"sort"
	"time"

	"errors"
	"github.com/strongo/decimal"
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

// ffjson: skip
type Balanced struct {
	BalanceJson      string    `datastore:",noindex,omitempty" json:",omitempty"`
	LastTransferID   string    `datastore:",noindex,omitempty" json:",omitempty"`
	LastTransferAt   time.Time `datastore:",noindex,omitempty"`           // `json:",omitempty"` - http://stackoverflow.com/questions/32643815/golang-json-omitempty-with-time-time-field
	CountOfTransfers int       `datastore:",omitempty" json:",omitempty"` // Do not remove, need for hiding balance/history menu in Telegram
	BalanceCount     int       `datastore:",noindex,omitempty" json:"-"`
}

func (balanced *Balanced) Balance() (balance Balance) {
	if balanced.BalanceJson == "" || balanced.BalanceJson == "null" || balanced.BalanceJson == "nil" || balanced.BalanceJson == "{}" {
		balance = make(Balance, 1)
		return
	}
	balance = make(Balance, balanced.BalanceCount)
	if err := json.Unmarshal([]byte(balanced.BalanceJson), &balance); err != nil {
		panic(err)
	}
	return
}

func (balanced *Balanced) SetBalance(balance Balance) error {
	if len(balance) == 0 {
		balanced.BalanceJson = ""
		balanced.BalanceCount = 0
		return nil
	}
	for currency, val := range balance {
		if val == 0 {
			return errors.New("balance currency has 0 value: " + string(currency))
		}
	}
	if v, err := json.Marshal(balance); err != nil {
		return err
	} else {
		balanced.BalanceJson = string(v)
		balanced.BalanceCount = len(balance)
	}
	return nil
}

func (balanced *Balanced) AddToBalance(currency CurrencyCode, value decimal.Decimal64p2) (Balance, error) {
	oldBalance := balanced.Balance()
	newBalance := oldBalance.Add(Amount{Currency: currency, Value: value})
	//log.Debugf(c, "AddToBalance(): oldBalance: %v, newBalance: %v", oldBalance, newBalance)
	return newBalance, balanced.SetBalance(newBalance)
}
