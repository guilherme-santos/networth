package mock

import (
	"testing"
	"time"

	"github.com/guilherme-santos/networth"
)

type Storage struct {
	t *testing.T

	AssetsPeriodInvoked bool
	AssetsPeriodFn      func(period string) (map[int]map[string]networth.ByMonth, error)

	AssetsInInvoked bool
	AssetsInFn      func(year int, month time.Month) (map[int]map[string]networth.ByMonth, error)

	InsertInvoked bool
	InsertFn      func(networth.Asset) (id string, err error)

	GetInvoked bool
	GetFn      func(id string) (networth.Asset, error)

	UpdateInvoked bool
	UpdateFn      func(id string, a networth.Asset) error

	DeleteInvoked bool
	DeleteFn      func(id string) error
}

func NewStorage(t *testing.T) *Storage {
	return &Storage{
		t: t,
	}
}

func (s *Storage) AssetsPeriod(period string) (map[int]map[string]networth.ByMonth, error) {
	if s.AssetsPeriodFn == nil {
		s.t.Fatal("You need to set AssetsPeriodFn before use this mock")
	}

	s.AssetsPeriodInvoked = true
	return s.AssetsPeriodFn(period)
}

func (s *Storage) AssetsIn(year int, month time.Month) (map[int]map[string]networth.ByMonth, error) {
	if s.AssetsInFn == nil {
		s.t.Fatal("You need to set AssetsInFn before use this mock")
	}

	s.AssetsInInvoked = true
	return s.AssetsInFn(year, month)
}

func (s *Storage) Insert(a networth.Asset) (string, error) {
	if s.InsertFn == nil {
		s.t.Fatal("You need to set InsertFn before use this mock")
	}

	s.InsertInvoked = true
	return s.InsertFn(a)
}

func (s *Storage) Get(id string) (networth.Asset, error) {
	if s.GetFn == nil {
		s.t.Fatal("You need to set GetFn before use this mock")
	}

	s.GetInvoked = true
	return s.GetFn(id)
}

func (s *Storage) Update(id string, a networth.Asset) error {
	if s.UpdateFn == nil {
		s.t.Fatal("You need to set UpdateFn before use this mock")
	}

	s.UpdateInvoked = true
	return s.UpdateFn(id, a)
}

func (s *Storage) Delete(id string) error {
	if s.DeleteFn == nil {
		s.t.Fatal("You need to set DeleteFn before use this mock")
	}

	s.DeleteInvoked = true
	return s.DeleteFn(id)
}
