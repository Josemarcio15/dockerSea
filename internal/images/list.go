package images

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go-walis/internal/containers"
	"go-walis/internal/core/connection"
)

// ListImages lista as imagens do Docker Daemon e correlaciona com containers ativos para identificar quais estão em uso
func ListImages(client *connection.Client, all bool) ([]DockerImage, error) {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return nil, fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	// 1. Obter lista de containers para mapear quais imagens estão em uso
	containerList, _ := containers.ListContainers(client, true)
	imageUsageMap := make(map[string][]string) // imageID/name -> []containerNames
	for _, c := range containerList {
		cName := c.Name
		if cName == "" && len(c.Names) > 0 {
			cName = strings.TrimPrefix(c.Names[0], "/")
		}

		if c.Image != "" {
			imageUsageMap[c.Image] = append(imageUsageMap[c.Image], cName)
			// Se o container tem repo:tag ou repo
			if idx := strings.LastIndex(c.Image, ":"); idx != -1 {
				repoOnly := c.Image[:idx]
				imageUsageMap[repoOnly] = append(imageUsageMap[repoOnly], cName)
			}
		}
		if c.ImageID != "" {
			imageUsageMap[c.ImageID] = append(imageUsageMap[c.ImageID], cName)
			cleanID := strings.TrimPrefix(c.ImageID, "sha256:")
			imageUsageMap[cleanID] = append(imageUsageMap[cleanID], cName)
			if len(cleanID) >= 12 {
				imageUsageMap[cleanID[:12]] = append(imageUsageMap[cleanID[:12]], cName)
			}
		}
	}

	// 2. Consultar GET /images/json
	endpoint := fmt.Sprintf("http://localhost/images/json?all=%t", all)
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar Docker API /images/json: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Docker API /images/json retornou status %d: %s", resp.StatusCode, string(body))
	}

	var rawList []RawDockerImage
	if err := json.NewDecoder(resp.Body).Decode(&rawList); err != nil {
		return nil, fmt.Errorf("falha ao decodificar lista de imagens: %w", err)
	}

	var result []DockerImage
	for _, raw := range rawList {
		fullID := strings.TrimPrefix(raw.ID, "sha256:")
		shortID := fullID
		if len(shortID) > 12 {
			shortID = shortID[:12]
		}

		// Identificar quais containers utilizam esta imagem
		var containersUsing []string
		seen := make(map[string]bool)
		addContainers := func(list []string) {
			for _, name := range list {
				if !seen[name] {
					seen[name] = true
					containersUsing = append(containersUsing, name)
				}
			}
		}

		if list, ok := imageUsageMap[raw.ID]; ok {
			addContainers(list)
		}
		if list, ok := imageUsageMap[shortID]; ok {
			addContainers(list)
		}

		// Se a imagem não possui tags (dangling / <none>:<none>)
		if len(raw.RepoTags) == 0 {
			result = append(result, DockerImage{
				ID:              shortID,
				Repo:            "<none>",
				Tag:             "<none>",
				Size:            formatBytes(raw.Size),
				RawSizeBytes:    raw.Size,
				Created:         raw.Created,
				ContainersUsing: containersUsing,
				VirtualSize:     raw.VirtualSize,
				SharedSize:      raw.SharedSize,
			})
			continue
		}

		// Criar uma entrada para cada tag associada à imagem
		for _, repoTag := range raw.RepoTags {
			repo := repoTag
			tag := "latest"

			if list, ok := imageUsageMap[repoTag]; ok {
				addContainers(list)
			}

			if idx := strings.LastIndex(repoTag, ":"); idx != -1 {
				repo = repoTag[:idx]
				tag = repoTag[idx+1:]
			}

			if list, ok := imageUsageMap[repo]; ok {
				addContainers(list)
			}

			result = append(result, DockerImage{
				ID:              shortID,
				Repo:            repo,
				Tag:             tag,
				Size:            formatBytes(raw.Size),
				RawSizeBytes:    raw.Size,
				Created:         raw.Created,
				ContainersUsing: containersUsing,
				VirtualSize:     raw.VirtualSize,
				SharedSize:      raw.SharedSize,
			})
		}
	}

	return result, nil
}

// RemoveImage remove uma imagem pelo ID ou RepoTag via DELETE /images/{name}?force=...
func RemoveImage(client *connection.Client, nameOrID string, force bool) error {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	endpoint := fmt.Sprintf("http://localhost/images/%s?force=%t", url.PathEscape(nameOrID), force)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha na requisição para deletar imagem: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao deletar imagem (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
