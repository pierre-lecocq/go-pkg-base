# logger

Configures the global `slog` logger with a colorized output handler writing to stderr.

Log level is controlled by the `LOGGER_LEVEL` environment variable (integer matching `slog.Level`: -4=Debug, 0=Info, 4=Warn, 8=Error). Defaults to Info.

```go
// Call once at startup, before any logging occurs.
logger.SetGlobalLogger()

slog.Info("server started", "addr", ":8080")
slog.Error("query failed", "error", err)
```
