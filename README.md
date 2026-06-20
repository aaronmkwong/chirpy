# Chirpy

A RESTful HTTP API server built in Go, backed by PostgreSQL. Chirpy mimics the
core backend of a microblogging platform (think Twitter/X) — letting users
register and post short messages called "chirps."

## Project Description

Chirpy exposes a JSON API over HTTP that supports user management and chirp
(short message) creation. It uses raw SQL via `sqlc` for type-safe database
access, and `goose` for schema migrations. The server is structured around
Go's standard `net/http` package with no external web framework.

## Coding Concepts

- **HTTP routing & handlers** — registering method-specific routes and
  writing handler functions that read request bodies and write JSON responses.
- **Database migrations** — versioned SQL schema changes (up/down) that
  evolve the database structure safely over time.
- **Foreign keys & referential integrity** — linking chirps to users via a
  `user_id` foreign key with `ON DELETE CASCADE`, so deleting a user also
  removes their chirps automatically.
- **UUIDs as primary keys** — using universally unique identifiers instead of
  auto-incrementing integers for distributed-friendly, opaque resource IDs.
- **Input validation** — enforcing business rules (e.g., max chirp length,
  banned words) before persisting data.
- **Middleware** — wrapping handlers to track metrics like fileserver hit counts.
- **Environment-based configuration** — reading secrets and settings from
  environment variables rather than hardcoding them.

## Applications

| Concept | Application |
|---|---|
| REST + JSON API | GitHub API, Stripe, Twilio |
| PostgreSQL + migrations | Nearly every production SaaS backend |
| Foreign keys with CASCADE | E-commerce (orders deleted with account) |
| UUID primary keys | AWS resource IDs, Stripe object IDs |
| Input validation at the API layer | Any form submission, payment processing |
| Middleware chains | Logging, auth, rate-limiting in Express/Django/Rails |

## Appendix: Go & HTTP Concepts Cheatsheet

A plain-language reference for the core concepts used in this project.

### Go Language Basics

| Term | Plain language |
|------|----------------|
| **struct** | A custom data type that groups related fields together. Like a labeled box with named compartments. |
| **struct tag** | The backtick text after a field (`` `json:"user_id"` ``). Instructions telling a serializer what to name the field in JSON. |
| **method** | A function attached to a type. `func (cfg *apiConfig) handlerX(...)` — `handlerX` belongs to `apiConfig`. |
| **receiver** | The `(cfg *apiConfig)` part. The specific value the method runs against, like `self` in other languages. |
| **pointer** (`*`, `&`) | `*T` means "address of a T." `&x` gets x's address. Lets functions share/modify the same data instead of copying it. |
| **interface** | A contract: "any type with these methods qualifies." `http.Handler` is one. |
| **slice** (`[]T`) | A growable list. `[]string` is a list of strings. |
| **map** | Key→value lookup table. |
| **error** | A returned value (not a thrown exception) signaling something failed. Checked with `if err != nil`. |
| **goroutine** | A lightweight concurrent task. The server runs each request in its own one. |
| **atomic** | A counter safe to update from many goroutines at once without corruption. |

### HTTP / Server

| Term | Plain language |
|------|----------------|
| **handler** | A function that receives a request and writes a response. |
| **ServeMux** | The router. Matches incoming URLs/methods to the right handler. |
| **middleware** | A wrapper around handlers that runs before/after them. |
| **ResponseWriter** | Your outbox — write status code, headers, and body to it. |
| **Request** | The incoming inbox — method, URL, headers, body. |
| **status code** | A number signaling outcome: 200 OK, 201 Created, 400 client error, 500 server error. |

### JSON

| Term | Plain language |
|------|----------------|
| **serialize / Marshal** | Go struct → JSON bytes. Done when leaving your program. |
| **deserialize / Unmarshal / Decode** | JSON bytes → Go struct. Done when data arrives. |

### Storage

| Term | Plain language |
|------|----------------|
| **migration** | A versioned script that changes DB structure. "Up" applies it, "down" reverses it. |
| **Goose** | The tool that runs your migrations in order. |
| **sqlc** | Reads your SQL files and generates type-safe Go functions from them. |
| **query** | A reusable SQL operation (insert, select) sqlc turns into a Go function. |
| **foreign key** | A column that points at another table's row (chirp's `user_id` → user's `id`). |
| **ON DELETE CASCADE** | "If the parent row is deleted, delete its children too." |
| **context** | A request-scoped carrier for deadlines/cancellation, passed as the first arg to DB calls. |