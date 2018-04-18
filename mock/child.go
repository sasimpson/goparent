package mock

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
)

type MockChildService struct {
	Env       *config.Env
	Kid       *goparent.Child
	Deleted   int
	GetErr    error
	DeleteErr error
}

func (mcs MockChildService) Save(*goparent.Child) error {
	if mcs.GetErr != nil {
		return mcs.GetErr
	}
	return nil
}

func (mcs MockChildService) Child(string) (*goparent.Child, error) {
	if mcs.GetErr != nil {
		return nil, mcs.GetErr
	}
	return mcs.Kid, nil
}

func (mcs MockChildService) Delete(*goparent.Child) (int, error) {
	if mcs.DeleteErr != nil {
		return 0, mcs.DeleteErr
	}
	return mcs.Deleted, nil
}
