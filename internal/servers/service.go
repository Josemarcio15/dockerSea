package servers

import (
	"fmt"

	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
	"go-walis/internal/core/docker"
	sharedDocker "go-walis/internal/shared/docker"
)

// Service coordinates server management and connection diagnostics.
type Service struct {
	repository Repository
}

func NewService(database *db.DB) *Service {
	return NewServiceWithRepository(newDatabaseRepository(database))
}

func NewServiceWithRepository(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) ListServers() ([]db.VpsServer, error) {
	if s.repository == nil {
		return nil, fmt.Errorf("repositório de servidores não inicializado")
	}
	return s.repository.ListServers()
}

func (s *Service) SaveServer(server db.VpsServer) error {
	server.Name = NormalizeServerName(server.Name)
	if message := ValidateServer(server); message != "" {
		return fmt.Errorf("%s", message)
	}
	sharedDocker.CloseClient(server.ID)
	return s.repository.SaveServer(server)
}

func (s *Service) DeleteServer(id string) error {
	if id == "" {
		return fmt.Errorf("id inválido")
	}
	sharedDocker.CloseClient(id)
	return s.repository.DeleteServer(id)
}

func (s *Service) SetActiveServer(id string) error {
	if id == "" {
		return fmt.Errorf("id inválido")
	}
	return s.repository.SetActiveServer(id)
}

func (s *Service) TestConnection(server db.VpsServer) docker.DiagnosticResult {
	return docker.RunDiagnostic(server)
}

func (s *Service) AutoDetectDocker(server db.VpsServer) docker.DetectResult {
	return docker.AutoDetectEnvironment(server)
}

func (s *Service) GetSystemUsage(server db.VpsServer) (*connection.SystemUsage, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()
	return client.FetchSystemUsage()
}
