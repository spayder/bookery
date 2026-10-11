package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/spayder/bookery/internal/booking/domain"
)

var _ domain.CancellationPolicy = domain.NoticePeriodPolicy{}

func TestNewNoticePeriodPolicy_RejectsNegativeNotice(t *testing.T) {
	_, err := domain.NewNoticePeriodPolicy(-time.Minute)
	if !errors.Is(err, domain.ErrNegativeNotice) {
		t.Fatalf("error = %v, want %v", err, domain.ErrNegativeNotice)
	}
}

func TestNoticePeriodPolicy_CheckCancellation(t *testing.T) {
	slot := mustSlot(t, "2026-01-01T10:00:00Z", "2026-01-01T12:00:00Z")

	cases := []struct {
		name    string
		notice  time.Duration
		now     time.Time
		wantErr error
	}{
		{"no notice: 1 hour before start", 0, slot.Start().Add(-time.Hour), nil},
		{"no notice: 1 second before start", 0, slot.Start().Add(-time.Second), nil},
		{"no notice: exactly at start", 0, slot.Start(), domain.ErrTooLateToCancel},
		{"no notice: after start", 0, slot.Start().Add(time.Minute), domain.ErrTooLateToCancel},

		{"24h notice: 25 hours before start", 24 * time.Hour, slot.Start().Add(-25 * time.Hour), nil},
		{"24h notice: exactly 24 hours before start", 24 * time.Hour, slot.Start().Add(-24 * time.Hour), domain.ErrTooLateToCancel},
		{"24h notice: 23 hours before start", 24 * time.Hour, slot.Start().Add(-23 * time.Hour), domain.ErrTooLateToCancel},
		{"24h notice: after start", 24 * time.Hour, slot.Start().Add(time.Hour), domain.ErrTooLateToCancel},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			policy, err := domain.NewNoticePeriodPolicy(c.notice)
			if err != nil {
				t.Fatalf("unexpected error creating policy: %v", err)
			}

			err = policy.CheckCancellation(slot, c.now)
			if !errors.Is(err, c.wantErr) {
				t.Errorf("error = %v, want %v", err, c.wantErr)
			}
		})
	}
}
