# go-pkg-base

Reusable Go packages for HTTP service development.

## Packages

- [`config`](config/README.md) — load `.env` files and read typed environment variables
- [`database`](database/README.md) — SQLite connection management with WAL mode, connection pooling, and typed error helpers
- [`logger`](logger/README.md) — configure the global `slog` logger with colorized output
- [`response`](response/README.md) — HTTP response writers for JSON APIs and iCalendar payloads
- [`server`](server/README.md) — start an HTTP server with sensible timeouts and graceful shutdown

## Publishing a package

```sh
git commit -m "go-pkg-base: publish v1.0.0"
git tag v1.0.0
git push origin v1.0.0

GOPROXY=proxy.golang.org go list -m github.com/pierre-lecocq/go-pkg-base@v1.0.0
```
