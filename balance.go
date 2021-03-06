package money

import (
	"bytes"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/app"
	"github.com/strongo/decimal"
)

type TestStruct struct {
	ID   int64
	Name string
}

type Balance map[Currency]decimal.Decimal64p2

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

func (b Balance) CommaSeparatedUnsignedWithSymbols(translator strongo.SingleLocaleTranslator) string {
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
		//log.Debugf(c, "Balance.Add(): currency found: [%v], current=%v, newVal=%v", amount.Currency, current, newVal)
		if newVal == 0 {
			delete(b, amount.Currency)
		} else {
			b[amount.Currency] = newVal
		}
	} else {
		//log.Debugf(c, "Balance.Add(): currency NOT found: [%v], setting to: %v", amount.Currency, amount.Value)
		b[amount.Currency] = amount.Value
	}
	return b
}

// ffjson: skip
type Balanced struct {
	BalanceJson      string    `datastore:",noindex,omitempty" json:",omitempty"`
	LastTransferID   int64     `datastore:",noindex,omitempty" json:",omitempty"`
	LastTransferAt   time.Time `datastore:",noindex,omitempty"`           // `json:",omitempty"` - http://stackoverflow.com/questions/32643815/golang-json-omitempty-with-time-time-field
	CountOfTransfers int       `datastore:",omitempty" json:",omitempty"` // Do not remove, need for hiding balance/history menu in Telegram
	BalanceCount     int       `datastore:",noindex,omitempty" json:"-"`
}

func (b *Balanced) Balance() (balance Balance) {
	if b.BalanceJson == "" || b.BalanceJson == "null" || b.BalanceJson == "nil" || b.BalanceJson == "{}" {
		balance = make(Balance, 1)
		return
	}
	balance = make(Balance, b.BalanceCount)
	if err := ffjson.Unmarshal([]byte(b.BalanceJson), &balance); err != nil {
		panic(err)
	}
	return
}

func (b *Balanced) SetBalance(balance Balance) error {
	if balance == nil || len(balance) == 0 {
		b.BalanceJson = ""
		b.BalanceCount = 0
		return nil
	}
	for currency, val := range balance {
		if val == 0 {
			return errors.WithStack(errors.New("balance currency has 0 value: " + string(currency)))
		}
	}
	if v, err := ffjson.Marshal(balance); err != nil {
		return err
	} else {
		b.BalanceJson = string(v)
		b.BalanceCount = len(balance)
	}
	return nil
}

func (balanced *Balanced) AddToBalance(currency Currency, value decimal.Decimal64p2) (Balance, error) {
	oldBalance := balanced.Balance()
	newBalance := oldBalance.Add(Amount{Currency: currency, Value: value})
	//log.Debugf(c, "AddToBalance(): oldBalance: %v, newBalance: %v", oldBalance, newBalance)
	return newBalance, balanced.SetBalance(newBalance)
}
