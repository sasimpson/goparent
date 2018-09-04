package datastore

import "github.com/sasimpson/goparent"

//FamilyService -
type FamilyService struct {
}

//Save -
func (s *FamilyService) Save(*goparent.Family) error {
	panic("not implemented")
}

//Family -
func (s *FamilyService) Family(string) (*goparent.Family, error) {
	panic("not implemented")
}

//Children -
func (s *FamilyService) Children(*goparent.Family) ([]*goparent.Child, error) {
	panic("not implemented")
}

//AddMember -
func (s *FamilyService) AddMember(*goparent.Family, *goparent.User) error {
	panic("not implemented")
}

//GetAdminFamily -
func (s *FamilyService) GetAdminFamily(*goparent.User) (*goparent.Family, error) {
	panic("not implemented")
}
