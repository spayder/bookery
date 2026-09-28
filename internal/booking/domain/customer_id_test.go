package domain_test

import (
	"testing"

	"github.com/spayder/bookery/internal/booking/domain"
)

func TestNewCustomerID_GeneratesNonEmptyUniqueIDs(t *testing.T) {
	a := domain.NewCustomerID()
	b := domain.NewCustomerID()

	if a.String() == "" {
		t.Fatal("expected NewCustomerID to produce a non-empty ID")
	}
	if a == b {
		t.Fatalf("expected two calls to NewCustomerID to produce different IDs, both were %q", a)
	}
}

func TestCustomerIDFromString_RejectsEmpty(t *testing.T) {
	_, err := domain.CustomerIDFromString("")
	if err != domain.ErrEmptyCustomerID {
		t.Fatalf("expected ErrEmptyCustomerID, got %v", err)
	}
}

func TestCustomerIDFromString_RoundTrips(t *testing.T) {
	id, err := domain.CustomerIDFromString("cust-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := id.String(); got != "cust-1" {
		t.Fatalf("String() = %q, want %q", got, "cust-1")
	}
}

func TestCustomerID_Equality(t *testing.T) {
	a, _ := domain.CustomerIDFromString("cust-1")
	b, _ := domain.CustomerIDFromString("cust-1")
	c, _ := domain.CustomerIDFromString("cust-2")

	if a != b {
		t.Errorf("expected IDs built from the same string to be equal")
	}
	if a == c {
		t.Errorf("expected IDs built from different strings to be unequal")
	}
}
