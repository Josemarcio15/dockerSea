package servers

import "go-walis/internal/core/db"

type Repository interface {
	ListServers() ([]db.VpsServer, error)
	SaveServer(db.VpsServer) error
	DeleteServer(string) error
	SetActiveServer(string) error
}

type databaseRepository struct {
	database *db.DB
}

func newDatabaseRepository(database *db.DB) Repository {
	return &databaseRepository{database: database}
}

func (r *databaseRepository) ListServers() ([]db.VpsServer, error) {
	return r.database.ListVpsServers()
}

func (r *databaseRepository) SaveServer(server db.VpsServer) error {
	return r.database.SaveVpsServer(server)
}

func (r *databaseRepository) DeleteServer(id string) error {
	return r.database.DeleteVpsServer(id)
}

func (r *databaseRepository) SetActiveServer(id string) error {
	return r.database.SetActiveVpsServer(id)
}
