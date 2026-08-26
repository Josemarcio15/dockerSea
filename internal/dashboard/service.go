package dashboard

import (
	"fmt"
	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
)

// Service exposes the operations required by the server overview.
type Service struct {
	database *db.DB
}

func NewService(database *db.DB) *Service {
	return &Service{database: database}
}

func (s *Service) ListServers() ([]db.VpsServer, error) {
	if s.database == nil {
		return nil, fmt.Errorf("banco de dados SQLite não inicializado")
	}
	return s.database.ListVpsServers()
}

func (s *Service) SetActiveServer(id string) error {
	if id == "" {
		return fmt.Errorf("id inválido")
	}
	if s.database == nil {
		return fmt.Errorf("banco de dados SQLite não inicializado")
	}
	return s.database.SetActiveVpsServer(id)
}

func (s *Service) GetSystemUsage(server db.VpsServer) (*connection.SystemUsage, error) {
	client, err := connection.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return client.FetchSystemUsage()
}