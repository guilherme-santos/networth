package networth

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
)

// CurrencyConvert is a global instance of CurrencyConverter that
// will be used internally
var CurrencyConvert CurrencyConverter

type CurrencyConverter interface {
	Convert(amount decimal.Decimal, from, to string) (decimal.Decimal, error)
}

type Totaler interface {
	Total(currency.Unit) (decimal.Decimal, error)
}

type Storage interface {
	AssetsPeriod(period string) (map[int]map[string]ByMonth, error)
	AssetsIn(year int, month time.Month) (map[int]map[string]ByMonth, error)
	Insert(Asset) (id string, err error)
	Get(id string) (Asset, error)
	Update(id string, a Asset) error
	Delete(id string) error
}

var InvalidPeriodErr = errors.New("period is not valid")
