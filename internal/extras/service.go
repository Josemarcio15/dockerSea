package extras

import (
	"fmt"

	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
)

type ExtraService struct {
	database *db.DB
}

func NewExtraService(database *db.DB) *ExtraService {
	return &ExtraService{database: database}
}

// ListNginxSites lista os arquivos de configuração do Nginx em sites-available e sites-enabled
func (s *ExtraService) ListNginxSites(server db.VpsServer) (NginxSitesResult, error) {
	client, err := connection.NewClient(server)
	if err != nil {
		return NginxSitesResult{}, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return listNginxSites(client)
}

// ReadNginxSite lê o conteúdo de um arquivo de configuração do Nginx
func (s *ExtraService) ReadNginxSite(server db.VpsServer, filename string, directory string) (string, error) {
	client, err := connection.NewClient(server)
	if err != nil {
		return "", fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return readNginxSite(client, filename, directory)
}

// SaveNginxSite salva o conteúdo de um arquivo em /etc/nginx/sites-available
func (s *ExtraService) SaveNginxSite(server db.VpsServer, filename string, content string) ExtraActionResult {
	client, err := connection.NewClient(server)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	return saveNginxSite(client, filename, content)
}

// EnableNginxSite cria link simbólico de sites-available para sites-enabled
func (s *ExtraService) EnableNginxSite(server db.VpsServer, filename string) ExtraActionResult {
	client, err := connection.NewClient(server)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	return enableNginxSite(client, filename)
}

// DeleteNginxSite apaga o arquivo de sites-available ou o link de sites-enabled
func (s *ExtraService) DeleteNginxSite(server db.VpsServer, filename string, directory string) ExtraActionResult {
	client, err := connection.NewClient(server)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	return deleteNginxSite(client, filename, directory)
}

// TestNginxConfig executa `nginx -t` para validar sintaxe das configurações
func (s *ExtraService) TestNginxConfig(server db.VpsServer) ExtraActionResult {
	client, err := connection.NewClient(server)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	return testNginxConfig(client)
}

// RestartNginx reinicia ou recarrega o serviço Nginx
func (s *ExtraService) RestartNginx(server db.VpsServer) ExtraActionResult {
	client, err := connection.NewClient(server)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	return restartNginx(client)
}

// GetNginxLogs obtém os logs mais recentes de access.log e error.log do Nginx
func (s *ExtraService) GetNginxLogs(server db.VpsServer, lines int) (string, error) {
	client, err := connection.NewClient(server)
	if err != nil {
		return "", fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return getNginxLogs(client, lines)
}

// ListListeningPorts consulta portas TCP e UDP em escuta no servidor usando ss / netstat / lsof
func (s *ExtraService) ListListeningPorts(server db.VpsServer) ([]PortEntry, error) {
	client, err := connection.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return listListeningPorts(client)
}
