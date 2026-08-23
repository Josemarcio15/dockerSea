package networks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go-walis/internal/core/connection"
)

// ListNetworks lista todas as redes Docker da VPS
func ListNetworks(client *connection.Client) ([]DockerNetwork, error) {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return nil, fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	resp, err := httpClient.Get("http://localhost/networks")
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar Docker API /networks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Docker API /networks retornou status %d: %s", resp.StatusCode, string(body))
	}

	var rawList []RawDockerNetwork
	if err := json.NewDecoder(resp.Body).Decode(&rawList); err != nil {
		return nil, fmt.Errorf("falha ao decodificar lista de redes: %w", err)
	}

	var result []DockerNetwork
	for _, raw := range rawList {
		var subnet, gateway string
		if len(raw.IPAM.Config) > 0 {
			subnet = raw.IPAM.Config[0].Subnet
			gateway = raw.IPAM.Config[0].Gateway
		}

		var connectedContainers []DockerNetworkContainer
		for _, c := range raw.Containers {
			ip := c.IPv4Address
			if idx := strings.Index(ip, "/"); idx != -1 {
				ip = ip[:idx]
			}
			connectedContainers = append(connectedContainers, DockerNetworkContainer{
				Name:        c.Name,
				EndpointID:  c.EndpointID,
				MacAddress:  c.MacAddress,
				IPv4Address: ip,
				IPv6Address: c.IPv6Address,
			})
		}

		result = append(result, DockerNetwork{
			ID:         raw.ID,
			Name:       raw.Name,
			Driver:     raw.Driver,
			Scope:      raw.Scope,
			Subnet:     subnet,
			Gateway:    gateway,
			Internal:   raw.Internal,
			Attachable: raw.Attachable,
			Containers: connectedContainers,
			Labels:     raw.Labels,
		})
	}

	return result, nil
}

// CreateNetwork cria uma nova rede Docker via POST /networks/create
func CreateNetwork(client *connection.Client, req NetworkCreateRequest) error {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	payload := map[string]interface{}{
		"Name":           req.Name,
		"CheckDuplicate": true,
	}
	if req.Driver != "" {
		payload["Driver"] = req.Driver
	} else {
		payload["Driver"] = "bridge"
	}
	if len(req.Labels) > 0 {
		payload["Labels"] = req.Labels
	}

	if req.Subnet != "" || req.Gateway != "" {
		ipamConfig := map[string]string{}
		if req.Subnet != "" {
			ipamConfig["Subnet"] = req.Subnet
		}
		if req.Gateway != "" {
			ipamConfig["Gateway"] = req.Gateway
		}
		payload["IPAM"] = map[string]interface{}{
			"Config": []map[string]string{ipamConfig},
		}
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := httpClient.Post("http://localhost/networks/create", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("falha ao criar rede: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao criar rede (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// RemoveNetwork remove uma rede Docker via DELETE /networks/{id}
func RemoveNetwork(client *connection.Client, idOrName string) error {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	endpoint := fmt.Sprintf("http://localhost/networks/%s", url.PathEscape(idOrName))
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao remover rede: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao remover rede (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// PruneNetworks remove redes não utilizadas via POST /networks/prune
func PruneNetworks(client *connection.Client) (*RawNetworkPruneResponse, error) {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return nil, fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	resp, err := httpClient.Post("http://localhost/networks/prune", "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar prune de redes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro ao limpar redes (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var pruneResp RawNetworkPruneResponse
	if err := json.NewDecoder(resp.Body).Decode(&pruneResp); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta de prune: %w", err)
	}

	return &pruneResp, nil
}

// ConnectContainer conecta um container a uma rede via POST /networks/{id}/connect
func ConnectContainer(client *connection.Client, networkName, containerName string) error {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	payload := map[string]interface{}{
		"Container": containerName,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("http://localhost/networks/%s/connect", url.PathEscape(networkName))
	resp, err := httpClient.Post(endpoint, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("falha ao conectar container à rede: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao conectar container (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// DisconnectContainer desconecta um container de uma rede via POST /networks/{id}/disconnect
func DisconnectContainer(client *connection.Client, networkName, containerName string, force bool) error {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	payload := map[string]interface{}{
		"Container": containerName,
		"Force":     force,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("http://localhost/networks/%s/disconnect", url.PathEscape(networkName))
	resp, err := httpClient.Post(endpoint, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("falha ao desconectar container da rede: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao desconectar container (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
