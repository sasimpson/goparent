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

func (mws *MockWasteService) Save(*goparent.Waste) error {
	panic("not implemented")
}

func (mws *MockWasteService) Waste(*goparent.Family) ([]*goparent.Waste, error) {
	panic("not implemented")
}

func (mws *MockWasteService) Stats(*goparent.Child) (*goparent.WasteSummary, error) {
	if mws.StatErr != nil {
		return nil, mws.StatErr
	}
	if mws.Stat != nil {
		return mws.Stat, nil
	}

	return nil, nil
}
