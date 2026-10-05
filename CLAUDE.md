# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this project is

Bookery is a **learning project**: the user is relearning Go (day job is PHP) by building a booking system for a **coworking space** that rents rooms, desks and equipment. The goal is fluency in Go plus deliberate practice of **interfaces, concurrency, SOLID, clean code and DDD**. The code is the byproduct; the user's learning is the point.

## How to work with the user (teaching mode)

- **Fading scaffold.** For a new kind of concept, explain *why* it exists and where it fits in the roadmap before handing out any task. For a pattern the user has already practiced once, don't write it again: describe the requirements, let the user write the code **and the tests**, then review.
- Don't write production code or tests for the user unless they ask. If they seem stuck, offer to help first rather than taking over.
- Reviews should explain the Go/DDD idiom behind each problem, not just give the fix.
- The user's editor may not autosave. Before running tests on their code, `cat` the file to confirm what's on disk.
- Design each milestone so it exercises one of the focus areas, instead of treating them as side topics.

## Git

- Do **not** add a `Co-Authored-By: Claude` trailer (or any Claude attribution) to commit messages or PR descriptions. The user doesn't want Claude listed as a contributor.
- `main` is protected: never commit or push to it directly, and never force-push it. Work on a branch, push the branch, and open a PR. The CI `test` job (gofmt, vet, `go test -race`, in `.github/workflows/ci.yml`) must pass before merging.

## Commands

```bash
go test ./...                                   # all tests
go test -v -run TestBookingID ./internal/booking/domain/...   # single test / pattern
go test -race ./...                             # required once concurrency code exists
go test -count=1 ./...                          # bypass the test cache
gofmt -l . && go vet ./...                      # formatting + static checks (should print nothing)
```

## Architecture

Hexagonal (ports and adapters) inside one bounded context, `internal/booking/`:

- `domain/`: aggregates, value objects, domain errors, and repository **interfaces**. Pure Go with no DB, HTTP or framework imports.
- `application/`: use cases (one struct per use case) that orchestrate domain objects. Rules that span several aggregates, such as "no double booking", are enforced here, because a single `Booking` can't see the others.
- `infrastructure/`: adapters that implement the domain's interfaces (`memory/` first, Postgres later).
- `cmd/`: entry points that wire dependencies by hand.

**Dependency rule:** `domain` and `application` never import `infrastructure`. The consumer defines the interface.

**Inside `domain`, group by meaning, not by kind:** all domain types live in the one `domain` package. Don't create `entity/` or `valueobject/` subpackages (import cycles, forced exporting).

## Domain conventions

- **Ubiquitous language:** the aggregate is `Booking`, not "Reservation". The thing being booked is a generic `Resource`, on purpose, because the space rents several kinds of things.
- **Value objects** (`TimeSlot`, `BookingID`, `ResourceID`, `CustomerID`) have unexported fields and value receivers, and are built only through validating constructors, so an invalid instance can't exist.
- **ID pattern:** `NewXID()` takes no arguments and generates a UUID with the standard library `uuid` package (Go 1.27+; the project deliberately has no external dependencies, so don't add `github.com/google/uuid` back). `XIDFromString(s)` reconstructs an existing ID and validates it. Each ID is its own type so the compiler catches mixed-up arguments. The three ID types are intentionally *not* merged with generics.
- **Aggregates** change state only through business methods (`Confirm()`, `Cancel(now)`), never setters. They use pointer receivers and return `*Booking` from their constructor. A failed operation must leave the state unchanged.
- **Time:** domain methods take `now time.Time` as a parameter and never call `time.Now()`. A `Clock` interface belongs in the application layer.
- **Start-time boundary (deliberate business decision, the asymmetry is intended):** booking is allowed when `now` equals the slot start (walk-ins), but cancelling is *not* allowed from the slot start onward. Don't unify these two checks.
- **Errors** are package-level sentinel values (`var ErrX = errors.New(...)`), compared with `errors.Is` once wrapping is in use.
- **Tests** live in the external `domain_test` package (black-box tests), are table-driven where it fits, and share helpers such as `mustTime` in `timeslot_test.go`.

## Roadmap

Check the code for the current position. Each milestone lists its main practice area.

1. Value objects: `TimeSlot` and the ID types *(done)*
2. `Booking` aggregate: Pending → Confirmed → Cancelled lifecycle, and "no booking or cancelling once the slot has started"
3. `CancellationPolicy` strategies (open/closed principle)
4. Repository interface + in-memory implementation (interface segregation, dependency inversion)
5. `CreateBooking` / `CancelBooking` use cases, `Clock` interface, error wrapping
6. Reproduce and then fix the double-booking race condition (goroutines, mutex, `-race`, WaitGroup)
7. Runnable demo in `cmd/demo`
8. Async notifications + expiry of unconfirmed bookings (channels, worker pool, `select`, context)
9. Postgres via Docker Compose + a shared contract test suite for all repository implementations (Liskov substitution). Deliberately postponed until this milestone.
10. HTTP API (`net/http`) + graceful shutdown
