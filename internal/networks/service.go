package networks

import (
	"fmt"
	"strings"

	"go-walis/internal/containers"
	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
)

type NetworkService struct {
	database *db.DB
}

func NewNetworkService(database *db.DB) *NetworkService {
	return &NetworkService{
		database: database,
	}
}

// ListNetworks obtém a lista de redes da VPS
func (s *NetworkService) ListNetworks(server db.VpsServer) ([]DockerNetwork, error) {
	client, err := connection.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return ListNetworks(client)
}

// ListNetworkContainers lista containers da VPS para seleção nos modais de conexão
func (s *NetworkService) ListNetworkContainers(server db.VpsServer) ([]containers.Container, error) {
	client, err := connection.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return containers.ListContainers(client, false)
}

// CreateNetwork cria uma nova rede na VPS
func (s *NetworkService) CreateNetwork(server db.VpsServer, req NetworkCreateRequest) NetworkActionResult {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return NetworkActionResult{
			Success: false,
			Message: "Nome da rede não informado",
		}
	}

	client, err := connection.NewClient(server)
	if err != nil {
		return NetworkActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	if err := CreateNetwork(client, req); err != nil {
		return NetworkActionResult{
			Success: false,
			Message: fmt.Sprintf("Erro ao criar rede: %v", err),
		}
	}

	return NetworkActionResult{
		Success: true,
		Message: fmt.Sprintf("Rede '%s' criada com sucesso!", req.Name),
	}
}

// DeleteNetworks remove uma ou mais redes pelo nome ou ID
func (s *NetworkService) DeleteNetworks(server db.VpsServer, names []string) NetworkActionResult {
	if len(names) == 0 {
		return NetworkActionResult{
			Success: false,
			Message: "Nenhuma rede informada",
		}
	}

	client, err := connection.NewClient(server)
	if err != nil {
		return NetworkActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	var errList []string
	deletedCount := 0

	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}

		if err := RemoveNetwork(client, trimmed); err != nil {
			errList = append(errList, fmt.Sprintf("%s: %v", trimmed, err))
		} else {
			deletedCount++
		}
	}

	if len(errList) > 0 {
		return NetworkActionResult{
			Success: deletedCount > 0,
			Message: strings.Join(errList, "\n"),
			Count:   deletedCount,
			Errors:  errList,
		}
	}

	return NetworkActionResult{
		Success: true,
		Message: fmt.Sprintf("%d rede(s) removida(s) com sucesso!", deletedCount),
		Count:   deletedCount,
	}
}

// PruneNetworks remove redes não utilizadas
func (s *NetworkService) PruneNetworks(server db.VpsServer) NetworkActionResult {
	client, err := connection.NewClient(server)
	if err != nil {
		return NetworkActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	res, err := PruneNetworks(client)
	if err != nil {
		return NetworkActionResult{
			Success: false,
			Message: fmt.Sprintf("Erro ao limpar redes: %v", err),
		}
	}

	count := len(res.NetworksDeleted)
	return NetworkActionResult{
		Success: true,
		Message: fmt.Sprintf("%d rede(s) não utilizadas removidas!", count),
		Count:   count,
	}
}

// ConnectContainer conecta um container a uma rede
func (s *NetworkService) ConnectContainer(server db.VpsServer, networkName, containerName string) NetworkActionResult {
	client, err := connection.NewClient(server)
	if err != nil {
		return NetworkActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	if err := ConnectContainer(client, networkName, containerName); err != nil {
		return NetworkActionResult{
			Success: false,
			Message: fmt.Sprintf("Erro ao conectar container: %v", err),
		}
	}

	return NetworkActionResult{
		Success: true,
		Message: fmt.Sprintf("Container '%s' conectado à rede '%s' com sucesso!", containerName, networkName),
	}
}

// DisconnectContainer desconecta um container de uma rede
func (s *NetworkService) DisconnectContainer(server db.VpsServer, networkName, containerName string) NetworkActionResult {
	client, err := connection.NewClient(server)
	if err != nil {
		return NetworkActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	if err := DisconnectContainer(client, networkName, containerName, false); err != nil {
		return NetworkActionResult{
			Success: false,
			Message: fmt.Sprintf("Erro ao desconectar container: %v", err),
		}
	}

	return NetworkActionResult{
		Success: true,
		Message: fmt.Sprintf("Container '%s' desconectado da rede '%s'!", containerName, networkName),
	}
}
