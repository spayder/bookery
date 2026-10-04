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
- A booking can't be cancelled once its slot's start time has arrived.
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

Requires Go 1.23 or newer.

```bash
go test ./...                 # run all tests
go test -v ./...              # verbose, one line per test
go test -race ./...           # with the race detector
gofmt -l . && go vet ./...    # formatting and static checks
```

## Roadmap

- [x] Value objects: `TimeSlot` and typed IDs (`BookingID`, `ResourceID`, `CustomerID`)
- [ ] `Booking` aggregate with its lifecycle and rules *(in progress)*
- [ ] Cancellation policies as interchangeable strategies
- [ ] Repository interface with an in-memory implementation
- [ ] Use cases: create and cancel a booking
- [ ] Prevent double booking under concurrent requests
- [ ] Runnable demo
- [ ] Asynchronous notifications and expiry of unconfirmed bookings
- [ ] PostgreSQL storage (Docker Compose)
- [ ] HTTP API with graceful shutdown
