package config

import (
	"fmt"
	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
	"go-walis/internal/core/docker"
	sharedDocker "go-walis/internal/shared/docker"
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
	server.Name = NormalizeServerName(server.Name)
	if message := ValidateServerInput(server.Name, server.ConnectionType, server.Host); message != "" {
		return fmt.Errorf("%s", message)
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
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return client.FetchSystemUsage()
}

func (s *ConfigService) ListProfiles() ([]db.Profile, error) {
	return s.database.ListProfiles()
}

func (s *ConfigService) GetActiveProfile() (*db.Profile, error) {
	return s.database.GetActiveProfile()
}

func (s *ConfigService) SaveProfile(profile db.Profile) error {
	if profile.Name == "" {
		return fmt.Errorf("o nome do perfil é obrigatório")
	}
	return s.database.SaveProfile(profile)
}

func (s *ConfigService) DeleteProfile(id string) error {
	if id == "" {
		return fmt.Errorf("id do perfil inválido")
	}
	return s.database.DeleteProfile(id)
}

func (s *ConfigService) SetActiveProfile(id string) error {
	if id == "" {
		return fmt.Errorf("id do perfil inválido")
	}
	return s.database.SetActiveProfile(id)
}

func (s *ConfigService) SetProfileLocale(id string, locale string) error {
	if id == "" || locale == "" {
		return fmt.Errorf("parâmetros inválidos")
	}
	return s.database.UpdateProfileLocale(id, locale)
}

func (s *ConfigService) ListContainerConfigs(profileID string) ([]db.ContainerConfig, error) {
	return s.database.ListContainerConfigs(profileID)
}

func (s *ConfigService) SaveContainerConfig(config db.ContainerConfig) error {
	return s.database.SaveContainerConfig(config)
}

func (s *ConfigService) DeleteContainerConfig(id string) error {
	return s.database.DeleteContainerConfig(id)
}

type DatabaseInfo struct {
	Path string `json:"path"`
}

func (s *ConfigService) GetDatabaseInfo() DatabaseInfo {
	return DatabaseInfo{
		Path: s.database.GetDBPath(),
	}
}

func (s *ConfigService) ExportDatabaseBackup(destinationPath string) error {
	if destinationPath == "" {
		return fmt.Errorf("caminho de destino é obrigatório")
	}
	return s.database.Backup(destinationPath)
}

func (s *ConfigService) RestoreDatabaseBackup(sourcePath string) error {
	if sourcePath == "" {
		return fmt.Errorf("caminho de origem é obrigatório")
	}
	return s.database.Restore(sourcePath)
}

func (s *ConfigService) ResetDatabase() error {
	return s.database.Reset()
}

