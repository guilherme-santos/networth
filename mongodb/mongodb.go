package mongodb

import (
	"time"

	"github.com/guilherme-santos/networth"
)

type Storage struct {
}

func NewStorage(url string) *Storage {
	return &Storage{}
}

func (s *Storage) AssetsPeriod(period string) (map[int]map[string]networth.ByMonth, error) {
	return nil, nil
}

func (s *Storage) AssetsIn(year int, month time.Month) (map[int]map[string]networth.ByMonth, error) {
	return nil, nil
}
