package database

import (
	"errors"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func IsForeignKeyViolation(err error) bool {
	var sqliteErr *sqlite.Error
	return errors.As(err, &sqliteErr) &&
		sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY
}

func IsUniqueViolation(err error) bool {
	var sqliteErr *sqlite.Error
	return errors.As(err, &sqliteErr) &&
		sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}

func IsNoRows(err error) bool {
	var sqliteErr *sqlite.Error
	return errors.As(err, &sqliteErr) &&
		sqliteErr.Code() == sqlite3.SQLITE_NOTFOUND
}
