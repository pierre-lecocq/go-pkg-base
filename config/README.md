# config

Helpers for loading and reading environment variables.

```go
// Load a .env file (path comes from a flag or env var; no-op if empty)
if err := config.LoadEnvFileIfSet("/path/to/.env"); err != nil {
    log.Fatal(err)
}

// Assert required variables are present before starting
if err := config.ValidatePresenceOf("DATABASE_DSN", "PORT"); err != nil {
    log.Fatal(err)
}

// Read typed values
dsn, err := config.StringVal("DATABASE_DSN")
port, err := config.IntVal("PORT")
```
