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