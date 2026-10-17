package memory

import (
	"context"

	"github.com/spayder/bookery/internal/booking/domain"
)

type BookingRepository struct {
	bookings map[domain.BookingID]*domain.Booking
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{
		bookings: make(map[domain.BookingID]*domain.Booking),
	}
}

func (r *BookingRepository) FindByID(_ context.Context, id domain.BookingID) (*domain.Booking, error) {
	if b, ok := r.bookings[id]; ok {
		c := *b
		return &c, nil
	}
	return nil, domain.ErrBookingNotFound
}

func (r *BookingRepository) Save(_ context.Context, booking *domain.Booking) error {
	c := *booking
	r.bookings[booking.ID()] = &c
	return nil
}

func (r *BookingRepository) FindOverlapping(_ context.Context, resourceID domain.ResourceID, slot domain.TimeSlot) ([]*domain.Booking, error) {
	var result []*domain.Booking

	for _, b := range r.bookings {
		if b.ResourceID() != resourceID || !b.Slot().Overlaps(slot) || !b.IsActive() {
			continue
		}
		c := *b
		result = append(result, &c)
	}

	return result, nil
}
