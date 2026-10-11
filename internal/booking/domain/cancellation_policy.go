package domain

import (
	"errors"
	"time"
)

var ErrNegativeNotice = errors.New("negative notice period is not allowed")
var ErrTooLateToCancel = errors.New("too late to cancel")

type CancellationPolicy interface {
	CheckCancellation(slot TimeSlot, now time.Time) error
}

type NoticePeriodPolicy struct {
	notice time.Duration
}

func NewNoticePeriodPolicy(notice time.Duration) (NoticePeriodPolicy, error) {
	if notice < 0 {
		return NoticePeriodPolicy{}, ErrNegativeNotice
	}
	return NoticePeriodPolicy{notice: notice}, nil
}

func (p NoticePeriodPolicy) CheckCancellation(slot TimeSlot, now time.Time) error {
	deadline := slot.Start().Add(-p.notice)
	if !now.Before(deadline) {
		return ErrTooLateToCancel
	}
	return nil
}
