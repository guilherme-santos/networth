package http_test

import (
	"context"
	gohttp "net/http"
	"testing"
	"time"

	"github.com/foodora/go-ranger/fdhttp"
	"github.com/guilherme-santos/networth"
	"github.com/guilherme-santos/networth/http"
	"github.com/guilherme-santos/networth/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetByPeriod(t *testing.T) {
	storage := mock.NewStorage(t)
	storage.AssetsPeriodFn = func(period string) (map[int]map[string]networth.ByMonth, error) {
		assert.Equal(t, "1m", period)
		return map[int]map[string]networth.ByMonth{
			2018: map[string]networth.ByMonth{
				time.August.String(): networth.ByMonth{
					networth.Asset{},
				},
			},
		}, nil
	}

	h := http.NewHandler(storage)

	ctx := context.Background()
	ctx = fdhttp.SetRouteParams(ctx, map[string]string{"period": "1m"})

	statusCode, _ := h.Get(ctx)
	assert.Equal(t, gohttp.StatusOK, statusCode)
}

func TestGetByPeriod_InvalidPeriod(t *testing.T) {
	storage := mock.NewStorage(t)
	storage.AssetsPeriodFn = func(period string) (map[int]map[string]networth.ByMonth, error) {
		return nil, networth.InvalidPeriodErr
	}

	h := http.NewHandler(storage)

	ctx := context.Background()
	ctx = fdhttp.SetRouteParams(ctx, map[string]string{"period": "invalid"})

	statusCode, _ := h.Get(ctx)
	assert.Equal(t, gohttp.StatusBadRequest, statusCode)
}

func TestGetByPeriod_NoAsset(t *testing.T) {
	storage := mock.NewStorage(t)
	storage.AssetsPeriodFn = func(period string) (map[int]map[string]networth.ByMonth, error) {
		return nil, nil
	}

	h := http.NewHandler(storage)

	ctx := context.Background()
	ctx = fdhttp.SetRouteParams(ctx, map[string]string{"period": "1m"})

	statusCode, _ := h.Get(ctx)
	assert.Equal(t, gohttp.StatusNotFound, statusCode)
}

func TestGetByWhen(t *testing.T) {
	storage := mock.NewStorage(t)
	storage.AssetsInFn = func(year int, month time.Month) (map[int]map[string]networth.ByMonth, error) {
		assert.Equal(t, 2018, year)
		assert.Equal(t, time.August, month)
		return map[int]map[string]networth.ByMonth{
			2018: map[string]networth.ByMonth{
				time.August.String(): networth.ByMonth{
					networth.Asset{},
				},
			},
		}, nil
	}

	h := http.NewHandler(storage)

	ctx := context.Background()
	ctx = fdhttp.SetRouteParams(ctx, map[string]string{"when": "2018-08"})

	statusCode, _ := h.Get(ctx)
	assert.Equal(t, gohttp.StatusOK, statusCode)
}

func TestGetByWhen_NoAsset(t *testing.T) {
	storage := mock.NewStorage(t)
	storage.AssetsInFn = func(year int, month time.Month) (map[int]map[string]networth.ByMonth, error) {
		return nil, nil
	}

	h := http.NewHandler(storage)

	ctx := context.Background()
	ctx = fdhttp.SetRouteParams(ctx, map[string]string{"when": "2018-08"})

	statusCode, _ := h.Get(ctx)
	assert.Equal(t, gohttp.StatusNotFound, statusCode)
}
