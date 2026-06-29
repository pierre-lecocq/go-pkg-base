# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./...

# Test
go test ./...
go test ./... -run TestName   # single test

# Vet / lint
go vet ./...
```

No external build tooling (Make, Task, etc.) is used — standard `go` commands are sufficient.

## Architecture

This is a Go shared-library module (`github.com/pierre-lecocq/go-pkg-base`) with no `main` package. It provides reusable building blocks for HTTP service development. Each top-level directory is an independent package:

- **`config/`** — environment variable helpers. `LoadEnvFileIfSet` loads a `.env` file via `godotenv`; `ValidatePresenceOf` asserts required vars exist; `StringVal`/`IntVal` fetch typed env values.
- **`database/`** — SQLite connection management via `modernc.org/sqlite` (pure-Go, CGO-free). `Open(cfg)` returns a `*sql.DB` with pragmas applied per connection (WAL + `NORMAL` sync for file DBs; skipped for `:memory:`). `DBTx` is an interface satisfied by both `*sql.DB` and `*sql.Tx`, enabling repository functions to work transparently inside or outside a transaction. `errors.go` exposes typed error helpers (`IsForeignKeyViolation`, `IsUniqueViolation`, `IsNoRows`).
- **`logger/`** — `SetGlobalLogger()` replaces the default `slog` handler with a colorized `tint` handler writing to stderr. Log level is read from the `LOGGER_LEVEL` env var (integer matching `slog.Level`; defaults to `Info`).
- **`response/`** — HTTP response writers. `JSONResponse`/`JSONError` for JSON APIs; `ICSResponse` for iCalendar payloads.
- **`server/`** — `ServeWithGracefulShutdown(addr, handler)` starts an `http.Server` with hardcoded timeouts and drains connections on `SIGINT`/`SIGTERM` with a 15-second grace period.

The `database` defaults set `MaxIdle == MaxOpen` (both 10) intentionally — WAL mode requires this to prevent extra connections from being torn down when load drops.
