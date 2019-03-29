package currencyconverter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var APIURL = "https://free.currencyconverterapi.com/api/v5/convert?q=%s_%s&compact=y"

type CurrencyConverter struct {
	httpClient *http.Client
	cache      map[string]map[string]decimal.Decimal
}

func New() *CurrencyConverter {
	return &CurrencyConverter{
		httpClient: &http.Client{},
		cache:      map[string]map[string]decimal.Decimal{},
	}
}

func (c *CurrencyConverter) Convert(amount decimal.Decimal, from, to string) (decimal.Decimal, error) {
	if strings.EqualFold(from, to) {
		return amount, nil
	}

	if from, ok := c.cache[from]; ok {
		if rate, ok := from[to]; ok {
			return amount.Mul(rate).RoundBank(2), nil
		}
	}

	resp, err := c.httpClient.Get(fmt.Sprintf(APIURL, from, to))
	if err != nil {
		return amount, errors.Wrap(err, "currencyconverter: http request failed")
	}

	defer resp.Body.Close()

	var body map[string]interface{}
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()

	err = dec.Decode(&body)
	if err != nil {
		return amount, errors.Wrap(err, "currencyconverter: cannot decode response")
	}

	var rate decimal.Decimal

	for _, values := range body {
		for _, v := range values.(map[string]interface{}) {
			n := v.(json.Number)

			rate, err = decimal.NewFromString(n.String())
			if err != nil {
				return amount, errors.Wrapf(err, "currencyconverter: cannot convert %v to decimal.Decimal", v)
			}
			break
		}
		break
	}

	if rate.Equal(decimal.NewFromFloat(0)) {
		return amount, errors.New("currencyconverter: exchange rate was zero")
	}

	if c.cache[from] == nil {
		c.cache[from] = map[string]decimal.Decimal{}
	}
	c.cache[from][to] = rate

	return amount.Mul(rate).RoundBank(2), nil
}
