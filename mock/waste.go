package mock

import (
	"context"

	"github.com/sasimpson/goparent"
)

//WasteService -
type WasteService struct {
	Env      *goparent.Env
	Wastes   []*goparent.Waste
	Stat     *goparent.WasteSummary
	Graph    *goparent.WasteChartData
	GetErr   error
	StatErr  error
	GraphErr error
}

//Save -
func (m *WasteService) Save(context.Context, *goparent.Waste) error {
	if m.GetErr != nil {
		return m.GetErr
	}
	return nil
}

//Waste -
func (m *WasteService) Waste(context.Context, *goparent.Family, uint64) ([]*goparent.Waste, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.Wastes != nil {
		return m.Wastes, nil
	}
	return nil, nil
}

//Stats -
func (m *WasteService) Stats(context.Context, *goparent.Child) (*goparent.WasteSummary, error) {
	if m.StatErr != nil {
		return nil, m.StatErr
	}
	if m.Stat != nil {
		return m.Stat, nil
	}

	return nil, nil
}

//GraphData -
func (m *WasteService) GraphData(context.Context, *goparent.Child) (*goparent.WasteChartData, error) {
	if m.GraphErr != nil {
		return nil, m.GraphErr
	}
	if m.Graph != nil {
		return m.Graph, nil
	}
	return nil, nil
}
