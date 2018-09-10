package datastore

import "github.com/sasimpson/goparent"

//ChildService -
type ChildService struct {
	Env *goparent.Env
}

//Save -
func (s *ChildService) Save(*goparent.Child) error {
	panic("not implemented")
}

//Child -
func (s *ChildService) Child(string) (*goparent.Child, error) {
	panic("not implemented")
}

//Delete -
func (s *ChildService) Delete(*goparent.Child) (int, error) {
	panic("not implemented")
}
