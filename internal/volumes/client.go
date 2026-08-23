package volumes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go-walis/internal/containers"
	"go-walis/internal/core/connection"
)

// ListVolumes lista os volumes do Docker Daemon e correlaciona com containers em execução para verificar quais estão em uso
func ListVolumes(client *connection.Client) ([]DockerVolume, error) {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return nil, fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	// 1. Obter containers para mapear uso dos volumes
	containerList, _ := containers.ListContainers(client, true)
	volumeUsageMap := make(map[string][][]string) // volumeName -> [ [containerName, "ro" | "rw"] ]

	for _, c := range containerList {
		cName := c.Name
		if cName == "" && len(c.Names) > 0 {
			cName = strings.TrimPrefix(c.Names[0], "/")
		}

		for _, m := range c.Mounts {
			if m.Type == "volume" && m.Name != "" {
				roMode := "rw"
				if m.ReadOnly {
					roMode = "ro"
				}
				volumeUsageMap[m.Name] = append(volumeUsageMap[m.Name], []string{cName, roMode})
			}
		}
	}

	// 2. Consultar GET /volumes
	resp, err := httpClient.Get("http://localhost/volumes")
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar Docker API /volumes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Docker API /volumes retornou status %d: %s", resp.StatusCode, string(body))
	}

	var rawList RawDockerVolumeList
	if err := json.NewDecoder(resp.Body).Decode(&rawList); err != nil {
		return nil, fmt.Errorf("falha ao decodificar lista de volumes: %w", err)
	}

	var result []DockerVolume
	for _, raw := range rawList.Volumes {
		attachedContainers := volumeUsageMap[raw.Name]
		inUse := len(attachedContainers) > 0

		var sizeStr string
		if raw.UsageData != nil && raw.UsageData.Size > 0 {
			sizeStr = formatBytes(raw.UsageData.Size)
		}

		result = append(result, DockerVolume{
			Name:        raw.Name,
			Driver:      raw.Driver,
			Mountpoint:  raw.Mountpoint,
			CreatedAt:   raw.CreatedAt,
			Labels:      raw.Labels,
			Scope:       raw.Scope,
			InUse:       inUse,
			Containers:  attachedContainers,
			Size:        sizeStr,
		})
	}

	return result, nil
}

// CreateVolume cria um novo volume Docker via POST /volumes/create
func CreateVolume(client *connection.Client, req VolumeCreateRequest) error {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	payload := map[string]interface{}{
		"Name": req.Name,
	}
	if req.Driver != "" {
		payload["Driver"] = req.Driver
	}
	if len(req.Labels) > 0 {
		payload["Labels"] = req.Labels
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := httpClient.Post("http://localhost/volumes/create", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("falha ao criar volume: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao criar volume (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// RemoveVolume remove um volume Docker via DELETE /volumes/{name}?force=true
func RemoveVolume(client *connection.Client, name string, force bool) error {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	endpoint := fmt.Sprintf("http://localhost/volumes/%s?force=%t", url.PathEscape(name), force)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao deletar volume: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao deletar volume (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// PruneVolumes remove todos os volumes não utilizados via POST /volumes/prune
func PruneVolumes(client *connection.Client) (*RawVolumePruneResponse, error) {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return nil, fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	resp, err := httpClient.Post("http://localhost/volumes/prune", "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar prune de volumes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro ao limpar volumes (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var pruneResp RawVolumePruneResponse
	if err := json.NewDecoder(resp.Body).Decode(&pruneResp); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta de prune: %w", err)
	}

	return &pruneResp, nil
}

// GetVolumeSize calcula ou obtém o tamanho em disco de um volume usando comandos SSH ou docker system df
func GetVolumeSize(client *connection.Client, name string) (string, error) {
	// Primeiro tenta via du -sh no mountpoint inspecionado ou /var/lib/docker/volumes/<name>/_data
	cmd := fmt.Sprintf("docker volume inspect --format '{{ .Mountpoint }}' %s 2>/dev/null", url.QueryEscape(name))
	mountpoint, err := client.ExecCommand(cmd, false)
	mountpoint = strings.TrimSpace(mountpoint)
	if err == nil && mountpoint != "" && !strings.Contains(mountpoint, "Error") {
		sizeCmd := fmt.Sprintf("du -sh '%s' 2>/dev/null | cut -f1", mountpoint)
		out, errSize := client.ExecCommand(sizeCmd, false)
		if errSize == nil && strings.TrimSpace(out) != "" {
			return strings.TrimSpace(out), nil
		}
	}

	// Fallback padrão para diretório comum de volumes do Docker
	fallbackCmd := fmt.Sprintf("du -sh /var/lib/docker/volumes/%s/_data 2>/dev/null | cut -f1", name)
	out, err := client.ExecCommand(fallbackCmd, false)
	if err == nil && strings.TrimSpace(out) != "" {
		return strings.TrimSpace(out), nil
	}

	return "—", nil
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
