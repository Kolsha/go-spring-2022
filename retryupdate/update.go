//go:build !solution
// +build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
)

import "gitlab.com/slon/shad-go/retryupdate/kvapi"

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var oldValue *string = nil
	var oldVersion uuid.UUID
	var authErr *kvapi.AuthError
	for acquired := false; !acquired; {
		getRes, err := c.Get(&kvapi.GetRequest{Key: key})

		switch true {
		case errors.Is(err, kvapi.ErrKeyNotFound):
			acquired = true
		case err == nil:
			acquired = true
			oldValue = &getRes.Value
			oldVersion = getRes.Version
		case errors.As(err, &authErr):
			return err
		}
	}
	updated, err := updateFn(oldValue)
	if err != nil {
		return err
	}
	var conflictErr *kvapi.ConflictError
	newVer := uuid.Must(uuid.NewV4())
	for written := false; !written; {
		_, err := c.Set(&kvapi.SetRequest{Key: key, Value: updated, OldVersion: oldVersion, NewVersion: newVer})
		switch true {
		case errors.Is(err, kvapi.ErrKeyNotFound):
			oldVersion = uuid.UUID{}
			updated, err = updateFn(nil)
			if err != nil {
				return err
			}
		case err == nil || errors.As(err, &authErr):
			return err
		case errors.As(err, &conflictErr):
			if conflictErr.ExpectedVersion == newVer {
				return nil
			}
			return UpdateValue(c, key, updateFn)
		}
	}

	return nil
}
