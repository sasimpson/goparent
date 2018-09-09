package mock

import "github.com/sasimpson/goparent"

//FamilyService -
type FamilyService struct {
	Env    *goparent.Env
	Kids   []*goparent.Child
	GetErr error
}

//Save -
func (mfs *FamilyService) Save(*goparent.Family) error {
	panic("not implemented")
}

//Family -
func (mfs *FamilyService) Family(string) (*goparent.Family, error) {
	panic("not implemented")
}

//Children -
func (mfs *FamilyService) Children(*goparent.Family) ([]*goparent.Child, error) {
	if mfs.GetErr != nil {
		return nil, mfs.GetErr
	}
	return mfs.Kids, nil
}

//AddMember -
func (mfs *FamilyService) AddMember(*goparent.Family, *goparent.User) error {
	panic("not implemented")
}

//GetAdminFamily -
func (mfs *FamilyService) GetAdminFamily(*goparent.User) (*goparent.Family, error) {
	panic("not implemented")
}
