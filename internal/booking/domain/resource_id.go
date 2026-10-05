package domain

import (
	"errors"

	"uuid"
)

var ErrEmptyResourceID = errors.New("resource Id can not be empty")

type ResourceID struct {
	value string
}

func NewResourceID() ResourceID {
	return ResourceID{value: uuid.New().String()}
}

func ResourceIDFromString(value string) (ResourceID, error) {
	if value == "" {
		return ResourceID{}, ErrEmptyResourceID
	}
	return ResourceID{value: value}, nil
}

func (r ResourceID) String() string {
	return r.value
}
