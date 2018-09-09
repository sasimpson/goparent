package mock

import "github.com/sasimpson/goparent"

//ChildService -
type ChildService struct {
	Env       *goparent.Env
	Kid       *goparent.Child
	Deleted   int
	GetErr    error
	DeleteErr error
}

//Save -
func (mcs *ChildService) Save(*goparent.Child) error {
	if mcs.GetErr != nil {
		return mcs.GetErr
	}
	return nil
}

//Child -
func (mcs *ChildService) Child(string) (*goparent.Child, error) {
	if mcs.GetErr != nil {
		return nil, mcs.GetErr
	}
	return mcs.Kid, nil
}

//Delete -
func (mcs *ChildService) Delete(*goparent.Child) (int, error) {
	if mcs.DeleteErr != nil {
		return 0, mcs.DeleteErr
	}
	return mcs.Deleted, nil
}
