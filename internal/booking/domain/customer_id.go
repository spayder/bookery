package domain

import (
	"errors"

	"uuid"
)

var ErrEmptyCustomerID = errors.New("customer Id can not be empty")

type CustomerID struct {
	value string
}

func NewCustomerID() CustomerID {
	return CustomerID{value: uuid.New().String()}
}

func CustomerIDFromString(s string) (CustomerID, error) {
	if s == "" {
		return CustomerID{}, ErrEmptyCustomerID
	}
	return CustomerID{value: s}, nil
}

func (c CustomerID) String() string {
	return c.value
}
