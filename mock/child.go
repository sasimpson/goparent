package mock

import (
	"context"

	"github.com/sasimpson/goparent"
)

//ChildService -
type ChildService struct {
	Env       *goparent.Env
	Kid       *goparent.Child
	Deleted   int
	GetErr    error
	DeleteErr error
}

//Save -
func (mcs *ChildService) Save(context.Context, *goparent.Child) error {
	if mcs.GetErr != nil {
		return mcs.GetErr
	}
	return nil
}

//Child -
func (mcs *ChildService) Child(context.Context, string) (*goparent.Child, error) {
	if mcs.GetErr != nil {
		return nil, mcs.GetErr
	}
	return mcs.Kid, nil
}

//Delete -
func (mcs *ChildService) Delete(context.Context, *goparent.Child) (int, error) {
	if mcs.DeleteErr != nil {
		return 0, mcs.DeleteErr
	}
	return mcs.Deleted, nil
}
