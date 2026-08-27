package volumes

import (
	"fmt"
	"strings"

	"go-walis/internal/core/db"
	sharedDocker "go-walis/internal/shared/docker"
)

type VolumeService struct {
	database *db.DB
}

func NewVolumeService(database *db.DB) *VolumeService {
	return &VolumeService{
		database: database,
	}
}

// ListVolumes obtém a lista de volumes da VPS e correlaciona os containers em uso
func (s *VolumeService) ListVolumes(server db.VpsServer) ([]DockerVolume, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return ListVolumes(client)
}

// CreateVolume cria um novo volume Docker na VPS
func (s *VolumeService) CreateVolume(server db.VpsServer, req VolumeCreateRequest) VolumeActionResult {
	req = NormalizeVolumeCreateRequest(req)
	if message := ValidateVolumeCreateRequest(req); message != "" {
		return VolumeActionResult{
			Success: false,
			Message: message,
		}
	}

	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return VolumeActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	if err := CreateVolume(client, req); err != nil {
		return VolumeActionResult{
			Success: false,
			Message: fmt.Sprintf("Erro ao criar volume: %v", err),
		}
	}

	return VolumeActionResult{
		Success: true,
		Message: fmt.Sprintf("Volume '%s' criado com sucesso!", req.Name),
	}
}

// DeleteVolumes remove um ou mais volumes pelo nome
func (s *VolumeService) DeleteVolumes(server db.VpsServer, names []string) VolumeActionResult {
	if len(names) == 0 {
		return VolumeActionResult{
			Success: false,
			Message: "Nenhum volume selecionado",
		}
	}

	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return VolumeActionResult{
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

		if err := RemoveVolume(client, trimmed, true); err != nil {
			errList = append(errList, fmt.Sprintf("%s: %v", trimmed, err))
		} else {
			deletedCount++
		}
	}

	if len(errList) > 0 {
		return VolumeActionResult{
			Success: deletedCount > 0,
			Message: strings.Join(errList, "\n"),
			Count:   deletedCount,
			Errors:  errList,
		}
	}

	return VolumeActionResult{
		Success: true,
		Message: fmt.Sprintf("%d volume(s) removido(s) com sucesso!", deletedCount),
		Count:   deletedCount,
	}
}

// PruneVolumes remove volumes não utilizados na VPS
func (s *VolumeService) PruneVolumes(server db.VpsServer) VolumeActionResult {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return VolumeActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	res, err := PruneVolumes(client)
	if err != nil {
		return VolumeActionResult{
			Success: false,
			Message: fmt.Sprintf("Erro ao limpar volumes: %v", err),
		}
	}

	count := len(res.VolumesDeleted)
	return VolumeActionResult{
		Success: true,
		Message: fmt.Sprintf("%d volume(s) não utilizados removidos!", count),
		Count:   count,
	}
}

// GetVolumeSize calcula o tamanho ocupado pelo volume
func (s *VolumeService) GetVolumeSize(server db.VpsServer, name string) (string, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return "—", fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return GetVolumeSize(client, name)
}
