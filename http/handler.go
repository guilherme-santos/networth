package http

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/foodora/go-ranger/fdhttp"
	"github.com/guilherme-santos/networth"
	"github.com/pkg/errors"
)

type Handler struct {
	storage networth.Storage
	router  *fdhttp.Router
}

var NoAssetErr = &fdhttp.Error{
	Code:    "no_asset_found",
	Message: "No asset was found to this date",
}

func NewHandler(storage networth.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) Init(r *fdhttp.Router) {
	h.router = r

	r = r.SubRouter()
	r.Prefix = "/v1/assets"

	r.POST("/", h.New)
	// https://www.networth.com/v1/assets
	// https://www.networth.com/v1/assets?when=2018
	// https://www.networth.com/v1/assets?when=2018-06
	// https://www.networth.com/v1/assets?period=6m
	// https://www.networth.com/v1/assets?period=1y
	// https://www.networth.com/v1/assets?period=all
	r.GET("/", h.Get)
	r.GET("/:id", h.GetByID).Name("asset_by_id")
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

func splitYearMonth(when string) (int, time.Month) {
	now := time.Now()

	arr := strings.Split(when, "-")
	if len(arr) == 0 {
		return now.Year(), time.Month(0)
	}

	year, err := strconv.Atoi(arr[0])
	if err != nil {
		year = now.Year()
	}

	if len(arr) == 1 {
		return year, time.Month(0)
	}

	month, err := strconv.Atoi(arr[1])
	if err != nil {
		month = 0
	}

	return year, time.Month(month)
}

func (h *Handler) New(ctx context.Context) (int, interface{}) {
	var asset networth.Asset

	err := fdhttp.RequestBodyJSON(ctx, &asset)
	if err != nil {
		return http.StatusBadRequest, fdhttp.Error{
			Code:    "invalid_body",
			Message: err.Error(),
		}
	}

	id, err := h.storage.Insert(asset)
	if err != nil {
		return http.StatusInternalServerError, fdhttp.Error{
			Code:    "unknown",
			Message: err.Error(),
		}
	}

	url := h.router.URLParam("asset_by_id", map[string]string{
		"id": id,
	})
	fdhttp.SetResponseHeaderValue(ctx, "Location", url)

	return http.StatusCreated, nil
}

func (h *Handler) Get(ctx context.Context) (int, interface{}) {
	period := fdhttp.RouteParam(ctx, "period")
	if period != "" {
		assets, err := h.storage.AssetsPeriod(period)
		if errors.Cause(err) == networth.InvalidPeriodErr {
			return http.StatusBadRequest, fdhttp.Error{
				Code:    "invalid_period",
				Message: err.Error(),
			}
		}
		if err != nil {
			return http.StatusInternalServerError, fdhttp.Error{
				Code:    "unknown",
				Message: err.Error(),
			}
		}

		if len(assets) == 0 {
			return http.StatusNotFound, NoAssetErr
		}

		return http.StatusOK, assets
	}

	when := fdhttp.RouteParam(ctx, "when")
	year, month := splitYearMonth(when)

	assets, err := h.storage.AssetsIn(year, month)
	if err != nil {
		return http.StatusInternalServerError, fdhttp.Error{
			Code:    "unknown",
			Message: err.Error(),
		}
	}

	if len(assets) == 0 {
		return http.StatusNotFound, NoAssetErr
	}

	return http.StatusOK, assets
}

func (h *Handler) GetByID(ctx context.Context) (int, interface{}) {
	return http.StatusNotFound, nil
}

func (h *Handler) Update(ctx context.Context) (int, interface{}) {
	return http.StatusNotFound, nil
}

func (h *Handler) Delete(ctx context.Context) (int, interface{}) {
	return http.StatusNotFound, nil
}
