package stacks

import "go-walis/internal/core/db"

// Repository isolates stack persistence from deploy and compose workflows.
type Repository struct{ database *db.DB }

func NewRepository(database *db.DB) *Repository { return &Repository{database: database} }
func (r *Repository) List(profileID string) ([]db.Stack, error) {
	return r.database.ListStacks(profileID)
}
func (r *Repository) Get(id string) (*db.Stack, error) { return r.database.GetStack(id) }
func (r *Repository) Save(stack db.Stack) error        { return r.database.SaveStack(stack) }
func (r *Repository) Delete(id string) error           { return r.database.DeleteStack(id) }
