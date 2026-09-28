package domain_test

import (
	"testing"

	"github.com/spayder/bookery/internal/booking/domain"
)

func TestNewReservationID_GeneratesNonEmptyUniqueIDs(t *testing.T) {
	a := domain.NewReservationID()
	b := domain.NewReservationID()

	if a.String() == "" {
		t.Fatal("expected NewReservationID to produce a non-empty ID")
	}
	if a.String() == b.String() {
		t.Fatalf("expected two calls to NewReservationID to produce different IDs, both were %q", a.String())
	}
}

func TestReservationIDFromString_RejectsEmpty(t *testing.T) {
	_, err := domain.ReservationIDFromString("")
	if err != domain.ErrEmptyReservationID {
		t.Fatalf("expected ErrEmptyReservationID, got %v", err)
	}
}

func TestReservationIDFromString_RoundTrips(t *testing.T) {
	id, err := domain.ReservationIDFromString("res-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := id.String(); got != "res-123" {
		t.Fatalf("String() = %q, want %q", got, "res-123")
	}
}

func TestReservationID_Equality(t *testing.T) {
	a, err := domain.ReservationIDFromString("res-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := domain.ReservationIDFromString("res-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c, err := domain.ReservationIDFromString("res-456")
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
