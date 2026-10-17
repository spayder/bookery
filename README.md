# Bookery

A booking system for a coworking space: members book rooms, desks and equipment for a time slot.

This is a **learning project**. I'm using it to get fluent in Go again while practicing Domain-Driven Design, SOLID, clean code, interfaces and concurrency on a domain with real business rules. It's built step by step with test-driven development, so the history is meant to be readable.

## Domain

The core concept is a **Booking**: a customer reserves a resource for a time slot.

```
          Confirm()
Pending ───────────► Confirmed
   │                     │
   │ Cancel(now)         │ Cancel(now)
   ▼                     ▼
       Cancelled (final)
```

Business rules enforced so far:

- A booking can't be made for a slot that has already started. Walk-ins are allowed: a slot that starts exactly now can be booked.
- Only a pending booking can be confirmed.
- Only an active (pending or confirmed) booking can be cancelled.
- Whether a booking can still be cancelled depends on its **cancellation policy**, which is fixed when the booking is made, so later rule changes don't affect existing bookings:
  - **Notice period:** cancellation is allowed until a set time before the start, e.g. until the start itself for desks, or 24 hours before for meeting rooms.
  - **Non-cancellable:** the booking can't be cancelled at all, e.g. event spaces.
- A failed operation never changes a booking's state.

"Resource" is deliberately generic, because the space rents several kinds of things, not only rooms.

## Architecture

Hexagonal (ports and adapters) inside a single bounded context:

```
internal/booking/
  domain/           aggregates, value objects, domain errors, repository interfaces
  application/      use cases that orchestrate the domain
  infrastructure/   adapters: in-memory storage now, Postgres and HTTP later
cmd/                entry points that wire everything together
```

The domain layer is plain Go with no database, HTTP or framework dependencies, so business rules are tested in isolation. `domain` and `application` never import `infrastructure`; they define interfaces that the infrastructure implements.

## Getting started

Requires Go 1.27 or newer (it uses the standard library `uuid` package). There are no external dependencies.

```bash
go test ./...                 # run all tests
go test -v ./...              # verbose, one line per test
go test -race ./...           # with the race detector
gofmt -l . && go vet ./...    # formatting and static checks
```

## Roadmap

- [x] Value objects: `TimeSlot` and typed IDs (`BookingID`, `ResourceID`, `CustomerID`)
- [x] `Booking` aggregate with its lifecycle and rules 
- [x] Cancellation policies as interchangeable strategies 
- [x] Repository interface with an in-memory implementation 
- [ ] Use cases: create and cancel a booking *(in progress)*
- [ ] Prevent double booking under concurrent requests
- [ ] Runnable demo
- [ ] Asynchronous notifications and expiry of unconfirmed bookings
- [ ] PostgreSQL storage (Docker Compose)
- [ ] HTTP API with graceful shutdown
