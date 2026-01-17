# search-bff

A minimal **Search Backend-for-Frontend (BFF)** built in Go using `net/http`.

This project reflects my learning journey as a **frontend engineer** exploring backend fundamentals. The focus is on understanding how HTTP servers work, how requests move through a system, and how to write backend code that behaves predictably under load, cancellation, and partial failure.

The goal is **clarity and correctness**, not feature completeness.

---

## Why this project exists

Coming from frontend development, I wanted to go beyond “API consumer” knowledge and understand:

- What happens when an HTTP request reaches a server
- How request lifecycles and cancellation work
- How concurrency is handled safely in real services
- Why Go enforces strict package boundaries
- How small design mistakes lead to production bugs

This repository is the result of learning those fundamentals step by step.

---

## Architecture overview

Client  
↓  
HTTP handlers and middleware  
↓  
Service layer (fan-out / fan-in)  
↓  
Upstream clients  
↓  
Shared models (pure data)

### Package responsibilities

- **cmd/api** — server startup and graceful shutdown
- **internal/http** — handlers, middleware, request validation
- **internal/search** — business orchestration and concurrency
- **internal/clients** — upstream integrations (hotels, flights)
- **internal/model** — shared data structures, no logic
- **internal/infra** — cross-cutting infrastructure (logging)

**Design rule:** dependencies flow inward only. Import cycles are not allowed.

---

## API

### Search endpoint

```
GET /search?q=<query>
```

### Example request

```
/search?q=paris
```

### Successful response

Returns hotel and flight results. Either upstream may succeed independently.

### Error behavior

- Missing query parameter → `400 Bad Request`
- Unsupported HTTP method → `405 Method Not Allowed`
- All upstreams fail → `502 Bad Gateway`
- Request timeout → `504 Gateway Timeout`

Partial failures are allowed.

---

## Key concepts learned

### Boundary validation

All request validation happens at the HTTP layer. Invalid input fails fast.

### Request lifecycle

One request maps to one goroutine. Context is created once and propagated. Cancellation flows from client to downstream calls.

### Concurrency

Independent upstream calls run concurrently. Goroutines never mutate shared state. Results are combined only after synchronization.

### Middleware discipline

Headers are set before handlers run. Status codes are written once. Business logic does not live in middleware.

---

## What this project intentionally avoids

This is not a full production system. It avoids:

- Frameworks
- Databases
- Authentication
- Caching
- Background jobs
- Global state

The focus is on understanding core backend mechanics first.

---

## Running locally

```
go run cmd/api/main.go
```

Server starts at:

```
http://localhost:8080
```

---

## Verification

```
go build ./...
go test -race ./...
```

The project builds cleanly and is free of data races.

---

## Who this project is for

- Frontend engineers learning backend fundamentals
- Developers new to Go who want to understand `net/http`

---

## Reflection

This project helped me understand why backend issues are often about flow rather than syntax, how concurrency mistakes affect production systems, and why Go’s strictness is useful when building reliable services.
