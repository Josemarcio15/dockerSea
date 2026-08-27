package profiles

import "go-walis/internal/core/db"

type Repository interface {
	ListProfiles() ([]db.Profile, error)
	GetActiveProfile() (*db.Profile, error)
	SaveProfile(db.Profile) error
	DeleteProfile(string) error
	SetActiveProfile(string) error
	UpdateProfileLocale(string, string) error
}

type databaseRepository struct {
	database *db.DB
}

func newDatabaseRepository(database *db.DB) Repository {
	return &databaseRepository{database: database}
}

func (r *databaseRepository) ListProfiles() ([]db.Profile, error) {
	return r.database.ListProfiles()
}

func (r *databaseRepository) GetActiveProfile() (*db.Profile, error) {
	return r.database.GetActiveProfile()
}

func (r *databaseRepository) SaveProfile(profile db.Profile) error {
	return r.database.SaveProfile(profile)
}

func (r *databaseRepository) DeleteProfile(id string) error {
	return r.database.DeleteProfile(id)
}

func (r *databaseRepository) SetActiveProfile(id string) error {
	return r.database.SetActiveProfile(id)
}

func (r *databaseRepository) UpdateProfileLocale(id, locale string) error {
	return r.database.UpdateProfileLocale(id, locale)
}
