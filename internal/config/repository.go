package config

import "go-walis/internal/core/db"

// Repository isolates configuration persistence from diagnostics and services.
type Repository struct{ database *db.DB }

func NewRepository(database *db.DB) *Repository            { return &Repository{database: database} }
func (r *Repository) Servers() ([]db.VpsServer, error)     { return r.database.ListVpsServers() }
func (r *Repository) SaveServer(server db.VpsServer) error { return r.database.SaveVpsServer(server) }
func (r *Repository) DeleteServer(id string) error         { return r.database.DeleteVpsServer(id) }
func (r *Repository) Profiles() ([]db.Profile, error)      { return r.database.ListProfiles() }
