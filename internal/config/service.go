package config

import (
	"fmt"
	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
	"go-walis/internal/core/docker"
)

type ConfigService struct {
	database *db.DB
}

func NewConfigService(database *db.DB) *ConfigService {
	return &ConfigService{database: database}
}

func (s *ConfigService) ListServers() ([]db.VpsServer, error) {
	return s.database.ListVpsServers()
}

func (s *ConfigService) SaveServer(server db.VpsServer) error {
	if server.Name == "" {
		return fmt.Errorf("o nome do servidor é obrigatório")
	}
	if server.ConnectionType == "ssh" && server.Host == "" {
		return fmt.Errorf("o host/IP é obrigatório para conexões SSH")
	}
	return s.database.SaveVpsServer(server)
}

func (s *ConfigService) DeleteServer(id string) error {
	if id == "" {
		return fmt.Errorf("id inválido")
	}
	return s.database.DeleteVpsServer(id)
}

func (s *ConfigService) SetActiveServer(id string) error {
	if id == "" {
		return fmt.Errorf("id inválido")
	}
	return s.database.SetActiveVpsServer(id)
}

func (s *ConfigService) TestConnection(server db.VpsServer) docker.DiagnosticResult {
	return docker.RunDiagnostic(server)
}

func (s *ConfigService) AutoDetectDocker(server db.VpsServer) docker.DetectResult {
	return docker.AutoDetectEnvironment(server)
}

func (s *ConfigService) GetSystemUsage(server db.VpsServer) (*connection.SystemUsage, error) {
	client, err := connection.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return client.FetchSystemUsage()
}
