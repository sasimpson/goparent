package datastore

import "github.com/sasimpson/goparent"

//WasteService -
type WasteService struct {
}

//Save -
func (s *WasteService) Save(*goparent.Waste) error {
	panic("not implemented")
}

//Waste -
func (s *WasteService) Waste(*goparent.Family, uint64) ([]*goparent.Waste, error) {
	panic("not implemented")
}

//Stats -
func (s *WasteService) Stats(*goparent.Child) (*goparent.WasteSummary, error) {
	panic("not implemented")
}

//GraphData -
func (s *WasteService) GraphData(*goparent.Child) (*goparent.WasteChartData, error) {
	panic("not implemented")
}
