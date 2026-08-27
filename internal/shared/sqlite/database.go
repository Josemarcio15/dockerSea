package sqlite

import (
	"database/sql"
	"fmt"
)

// Database is the shared boundary used by repositories.
type Database struct{ DB *sql.DB }

func Wrap(db *sql.DB) Database { return Database{DB: db} }

func (d Database) Close() error {
	if d.DB == nil {
		return nil
	}
	return d.DB.Close()
}

func (d Database) Exec(query string, args ...any) (sql.Result, error) {
	if d.DB == nil {
		return nil, fmt.Errorf("banco SQLite não inicializado")
	}
	return d.DB.Exec(query, args...)
}
