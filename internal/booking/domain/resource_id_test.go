package domain_test

import (
	"testing"

	"github.com/spayder/bookery/internal/booking/domain"
)

func TestNewResourceID_GeneratesNonEmptyUniqueIDs(t *testing.T) {
	a := domain.NewResourceID()
	b := domain.NewResourceID()

	if a.String() == "" {
		t.Fatal("expected NewResourceID to produce a non-empty ID")
	}
	if a == b {
		t.Fatalf("expected two calls to NewResourceID to produce different IDs, both were %q", a)
	}
}

func TestResourceIDFromString_RejectsEmpty(t *testing.T) {
	_, err := domain.ResourceIDFromString("")
	if err != domain.ErrEmptyResourceID {
		t.Fatalf("expected ErrEmptyResourceID, got %v", err)
	}
}

func TestResourceIDFromString_RoundTrips(t *testing.T) {
	id, err := domain.ResourceIDFromString("desk-42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := id.String(); got != "desk-42" {
		t.Fatalf("String() = %q, want %q", got, "desk-42")
	}
}

func TestResourceID_Equality(t *testing.T) {
	a, _ := domain.ResourceIDFromString("desk-42")
	b, _ := domain.ResourceIDFromString("desk-42")
	c, _ := domain.ResourceIDFromString("room-7")

	if a != b {
		t.Errorf("expected IDs built from the same string to be equal")
	}
	if a == c {
		t.Errorf("expected IDs built from different strings to be unequal")
	}
}
