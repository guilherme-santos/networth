package networth_test

import (
	"testing"

	"github.com/guilherme-santos/networth"
	"github.com/guilherme-santos/networth/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/currency"
)

func TestNewAsset(t *testing.T) {
	a, err := networth.NewAsset("Stocks", "1234.456", currency.EUR)
	assert.NoError(t, err)
	assert.Equal(t, "Stocks", a.Name)
	assert.Equal(t, currency.EUR, a.Currency)
	assert.Equal(t, "1234.456", a.Amount.String())
}

func TestAssetTotal(t *testing.T) {
	a, err := networth.NewAsset("Stocks", "1234.456", currency.EUR)
	assert.NoError(t, err)

	amount, err := a.Total(currency.EUR)
	assert.NoError(t, err)
	assert.Equal(t, "1234.46", amount.String())
}

func TestAssetTotal_DifferentCurrency(t *testing.T) {
	cc := mock.NewCurrencyConverter(t)
	cc.ConvertFn = func(amount decimal.Decimal, from, to string) (decimal.Decimal, error) {
		return amount.Mul(decimal.NewFromFloat(2)), nil
	}
	networth.CurrencyConvert = cc

	a, err := networth.NewAsset("Stocks", "1234.456", currency.EUR)
	assert.NoError(t, err)

	amount, err := a.Total(currency.BRL)
	assert.NoError(t, err)
	assert.Equal(t, "2468.91", amount.String())
}

func TestAssetGroupTotal(t *testing.T) {
	ag := networth.AssetGroup{}
	ag.Name = "Stocks"

	a, _ := networth.NewAsset("Fund", "1234.56", currency.EUR)
	ag.Assets = append(ag.Assets, a)

	a, _ = networth.NewAsset("Dividend", "2345.67", currency.EUR)
	ag.Assets = append(ag.Assets, a)

	amount, err := ag.Total(currency.EUR)
	assert.NoError(t, err)
	assert.Equal(t, "3580.23", amount.String())
}

func TestMonthTotal(t *testing.T) {
	ag := networth.AssetGroup{}
	ag.Name = "Stocks"

	a, _ := networth.NewAsset("Fund", "1234.56", currency.EUR)
	ag.Assets = append(ag.Assets, a)

	a, _ = networth.NewAsset("Dividend", "2345.67", currency.EUR)
	ag.Assets = append(ag.Assets, a)

	m := networth.ByMonth{}
	m = append(m, ag)

	a, _ = networth.NewAsset("Governmental Funds", "3456.78", currency.EUR)
	m = append(m, a)

	amount, err := m.Total(currency.EUR)
	assert.NoError(t, err)
	assert.Equal(t, "7037.01", amount.String())
}
