package datastore

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sasimpson/goparent"
	"google.golang.org/appengine/datastore"
)

//ChildService -
type ChildService struct {
	Env *goparent.Env
}

//ChildKind is the constant for child kind in gcp datastore
const ChildKind = "Child"

//Save -
func (s *ChildService) Save(ctx context.Context, child *goparent.Child) error {
	var childKey *datastore.Key
	if child.ID == "" {
		u := uuid.New()
		childKey = datastore.NewKey(ctx, ChildKind, u.String(), 0, nil)
		child.CreatedAt = time.Now()
		child.ID = u.String()
	} else {
		childKey = datastore.NewKey(ctx, ChildKind, child.ID, 0, nil)
	}
	child.LastUpdated = time.Now()
	_, err := datastore.Put(ctx, childKey, child)
	if err != nil {
		return NewError("ChildService.Save", err)
	}

	return nil
}

//Child -
func (s *ChildService) Child(ctx context.Context, id string) (*goparent.Child, error) {
	var child goparent.Child
	childKey := datastore.NewKey(ctx, ChildKind, id, 0, nil)
	err := datastore.Get(ctx, childKey, &child)
	if err != nil {
		return nil, NewError("ChildService.Child", err)
	}
	return &child, nil
}

//Delete -
func (s *ChildService) Delete(ctx context.Context, child *goparent.Child) (int, error) {
	childKey := datastore.NewKey(ctx, ChildKind, child.ID, 0, nil)
	err := datastore.Delete(ctx, childKey)
	if err != nil {
		return 0, err
	}
	return 1, nil
}
