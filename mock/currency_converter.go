package mock

import (
	"testing"

	"github.com/shopspring/decimal"
)

type CurrencyConverter struct {
	t *testing.T

	ConvertInvoked bool
	ConvertFn      func(amount decimal.Decimal, from, to string) (decimal.Decimal, error)
}

func NewCurrencyConverter(t *testing.T) *CurrencyConverter {
	return &CurrencyConverter{
		t: t,
	}
}

func (c *CurrencyConverter) Convert(amount decimal.Decimal, from, to string) (decimal.Decimal, error) {
	if c.ConvertFn == nil {
		c.t.Fatal("You need to set ConvertFn before use this mock")
	}

	c.ConvertInvoked = true
	return c.ConvertFn(amount, from, to)
}
