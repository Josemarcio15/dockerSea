package sqlite

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"testing"
)

func TestApply(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := Apply(db, "CREATE TABLE test (id INTEGER)"); err != nil {
		t.Fatal(err)
	}
}
