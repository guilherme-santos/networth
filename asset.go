package networth

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
)

type assetInfo struct {
	Name string
	Tags []string
}

// Asset ...
type Asset struct {
	assetInfo
	Currency currency.Unit
	Amount   decimal.Decimal
}

func NewAsset(name, amount string, cur currency.Unit) (Asset, error) {
	a := Asset{}
	a.Name = name
	a.Currency = cur

	var err error
	a.Amount, err = decimal.NewFromString(amount)
	if err != nil {
		return a, errors.Wrapf(err, "cannot convert %s to decimal.Decimal", amount)
	}

	return a, nil
}

func (a Asset) Total(cur currency.Unit) (decimal.Decimal, error) {
	if a.Currency.String() == "XXX" || cur.String() == "XXX" || cur == a.Currency {
		return a.Amount.RoundBank(2), nil
	}

	amount, err := CurrencyConvert.Convert(a.Amount, a.Currency.String(), cur.String())
	return amount.RoundBank(2), err
}

// AssetGroup ...
type AssetGroup struct {
	assetInfo
	Assets []Totaler
}

func (ag AssetGroup) Total(cur currency.Unit) (decimal.Decimal, error) {
	var groupTotal decimal.Decimal

	for _, a := range ag.Assets {
		t, err := a.Total(cur)
		if err != nil {
			return groupTotal, err
		}

		groupTotal = groupTotal.Add(t)
	}

	return groupTotal, nil
}

// // ByMonth ...
type ByMonth []Totaler

func (m ByMonth) Total(cur currency.Unit) (decimal.Decimal, error) {
	var monthTotal decimal.Decimal

	for _, a := range m {
		t, err := a.Total(cur)
		if err != nil {
			return monthTotal, err
		}

		monthTotal = monthTotal.Add(t)
	}

	return monthTotal, nil
}
