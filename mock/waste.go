package mock

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
)

type MockWasteService struct {
	Env     *config.Env
	Wastes  []*goparent.Waste
	Stat    *goparent.WasteSummary
	GetErr  error
	StatErr error
}

func (m *MockWasteService) Save(*goparent.Waste) error {
	if m.GetErr != nil {
		return m.GetErr
	}
	return nil
}

func (m *MockWasteService) Waste(*goparent.Family, uint64) ([]*goparent.Waste, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.Wastes != nil {
		return m.Wastes, nil
	}
	return nil, nil
}

func (m *MockWasteService) Stats(*goparent.Child) (*goparent.WasteSummary, error) {
	if m.StatErr != nil {
		return nil, m.StatErr
	}
	if m.Stat != nil {
		return m.Stat, nil
	}

	return nil, nil
}
