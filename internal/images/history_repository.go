package images

import "go-walis/internal/core/db"

// HistoryRepository isolates image-history persistence from Docker operations.
type HistoryRepository struct{ database *db.DB }

func NewHistoryRepository(database *db.DB) *HistoryRepository {
	return &HistoryRepository{database: database}
}

func (r *HistoryRepository) List(profileID string) ([]db.ImageHistoryItem, error) {
	return r.database.ListImageHistory(profileID)
}

func (r *HistoryRepository) Delete(ids []string) error { return r.database.DeleteImageHistory(ids) }

func (r *HistoryRepository) Clear(profileID string) error {
	return r.database.ClearImageHistory(profileID)
}
