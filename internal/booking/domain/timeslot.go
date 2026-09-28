package domain

import (
	"errors"
	"time"
)

var ErrInvalidTimeSlot = errors.New("time slot end must be after start")

// TimeSlot is a value object: two TimeSlots with the same Start/End
// are interchangeable, and once created it never changes.
type TimeSlot struct {
	start time.Time
	end   time.Time
}

func NewTimeSlot(start, end time.Time) (TimeSlot, error) {
	if !end.After(start) {
		return TimeSlot{}, ErrInvalidTimeSlot
	}
	return TimeSlot{start: start, end: end}, nil
}

func (t TimeSlot) Start() time.Time {
	return t.start
}

func (t TimeSlot) End() time.Time {
	return t.end
}

func (t TimeSlot) Overlaps(other TimeSlot) bool {
	return t.start.Before(other.end) && other.start.Before(t.end)
}
