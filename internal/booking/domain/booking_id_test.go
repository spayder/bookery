package domain_test

import (
	"testing"

	"github.com/spayder/bookery/internal/booking/domain"
)

func TestNewBookingID_GeneratesNonEmptyUniqueIDs(t *testing.T) {
	a := domain.NewBookingID()
	b := domain.NewBookingID()

	if a.String() == "" {
		t.Fatal("expected NewBookingID to produce a non-empty ID")
	}
	if a.String() == b.String() {
		t.Fatalf("expected two calls to NewBookingID to produce different IDs, both were %q", a.String())
	}
}

func TestBookingIDFromString_RejectsEmpty(t *testing.T) {
	_, err := domain.BookingIDFromString("")
	if err != domain.ErrEmptyBookingID {
		t.Fatalf("expected ErrEmptyBookingID, got %v", err)
	}
}

func TestBookingIDFromString_RoundTrips(t *testing.T) {
	id, err := domain.BookingIDFromString("bk-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := id.String(); got != "bk-123" {
		t.Fatalf("String() = %q, want %q", got, "bk-123")
	}
}

func TestBookingID_Equality(t *testing.T) {
	a, err := domain.BookingIDFromString("bk-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := domain.BookingIDFromString("bk-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c, err := domain.BookingIDFromString("bk-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a != b {
		t.Errorf("expected IDs built from the same string to be equal")
	}
	if a == c {
		t.Errorf("expected IDs built from different strings to be unequal")
	}
}
