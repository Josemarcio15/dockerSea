package sqlite

import (
	"database/sql"
	"fmt"
)

func Apply(db *sql.DB, migrations ...string) error {
	if db == nil {
		return fmt.Errorf("banco SQLite não inicializado")
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar transação de migration: %w", err)
	}
	for _, migration := range migrations {
		if _, err := tx.Exec(migration); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("aplicar migration: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar migrations: %w", err)
	}
	return nil
}
