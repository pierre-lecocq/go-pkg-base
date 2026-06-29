# database

SQLite connection management (CGO-free via `modernc.org/sqlite`) with WAL mode, connection pooling, and typed error helpers.

```go
cfg := database.NewConfigWithDSN("file:/var/data/app.db?mode=rwc")
db, err := database.Open(cfg)
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// DBTx is satisfied by both *sql.DB and *sql.Tx, so repository
// functions can run inside or outside a transaction transparently.
func insertUser(ctx context.Context, q database.DBTx, name string) error {
    _, err := q.ExecContext(ctx, "INSERT INTO users (name) VALUES (?)", name)
    if database.IsUniqueViolation(err) {
        return ErrAlreadyExists
    }

    return err
}
```
