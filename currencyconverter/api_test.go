package currencyconverter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guilherme-santos/networth/currencyconverter"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "EUR", r.URL.Query().Get("from"))
		assert.Equal(t, "BRL", r.URL.Query().Get("to"))
		fmt.Fprint(w, `{"EUR_BRL":{"val":4.32}}`)
	}))
	defer ts.Close()
	currencyconverter.APIURL = ts.URL + "?from=%s&to=%s"

	value, _ := decimal.NewFromString("1234.56")

	c := currencyconverter.New()
	res, err := c.Convert(value, "EUR", "BRL")
	assert.NoError(t, err)
	assert.Equal(t, "5333.3", res.String())
}

func TestConvert_SameCurrency(t *testing.T) {
	var srvCalled bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled = true
	}))
	defer ts.Close()

	value, _ := decimal.NewFromString("1234.56")

	c := currencyconverter.New()
	res, err := c.Convert(value, "EUR", "EUR")
	assert.NoError(t, err)
	assert.Equal(t, value, res)
	assert.False(t, srvCalled)
}

func TestConvert_UseCache(t *testing.T) {
	var srvCalled int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled++
		fmt.Fprint(w, `{"EUR_BRL":{"val":4.32}}`)
	}))
	defer ts.Close()
	currencyconverter.APIURL = ts.URL + "?from=%s&to=%s"

	c := currencyconverter.New()

	value, _ := decimal.NewFromString("1234.56")
	res, err := c.Convert(value, "EUR", "BRL")
	assert.NoError(t, err)
	assert.Equal(t, "5333.3", res.String())

	value, _ = decimal.NewFromString("6543.21")
	res, err = c.Convert(value, "EUR", "BRL")
	assert.NoError(t, err)
	assert.Equal(t, "28266.67", res.String())

	assert.Equal(t, 1, srvCalled)
}
