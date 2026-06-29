package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func isInMemoryDSN(dsn string) bool {
	dsn = strings.TrimSpace(dsn)

	return dsn == ":memory:" || strings.HasPrefix(dsn, "file::memory:")
}

func resolveDSN(dsn string) (string, error) {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return "", fmt.Errorf("missing DSN for database connection")
	}

	if isInMemoryDSN(dsn) {
		return dsn, nil
	}

	// Remove the prefix for filepath.Abs
	prefix := "file:"
	file := dsn
	if strings.HasPrefix(dsn, prefix) {
		file = dsn[len(prefix):]
	}

	// Separate query params, otherwise filepath.Abs will use them as part of the file name
	query := ""
	if i := strings.IndexByte(file, '?'); i != -1 {
		query = file[i:]
		file = file[:i]
	}

	if filepath.IsAbs(file) {
		return prefix + file + query, nil
	}

	abs, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	return prefix + abs + query, nil
}

func dsnWithPragmas(dsn string, pragmas []string) string {
	if len(pragmas) == 0 {
		return dsn
	}

	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}

	parts := make([]string, 0, len(pragmas))
	for _, p := range pragmas {
		parts = append(parts, "_pragma="+p)
	}

	return dsn + sep + strings.Join(parts, "&")
}

func Open(cfg *Config) (*sql.DB, error) {
	// Set DSN
	dsn, err := resolveDSN(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("resolve DSN: %w", err)
	}

	// Validate config
	if err := cfg.IsValid(); err != nil {
		return nil, err
	}

	// Setup pragmas for each connection in the pool
	pragmas := []string{
		fmt.Sprintf("busy_timeout(%d)", cfg.MaxBusy),
		"foreign_keys(on)",
	}

	if !isInMemoryDSN(dsn) {
		pragmas = append(pragmas, "journal_mode(WAL)", "synchronous(NORMAL)")
	}

	connDSN := dsnWithPragmas(dsn, pragmas)

	// Open connection
	db, err := sql.Open("sqlite", connDSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetConnMaxLifetime(cfg.MaxLifeTime)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	// Verify connection is actually working
	ctx, cancel := context.WithTimeout(context.Background(), cfg.QueryTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}
