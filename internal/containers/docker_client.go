package containers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go-walis/internal/core/connection"
)

// ListContainers lista todos os containers usando a API HTTP nativa do Docker Daemon
func ListContainers(client *connection.Client, all bool) ([]Container, error) {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return nil, fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	endpoint := fmt.Sprintf("http://localhost/containers/json?all=%t", all)
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar Docker API /containers/json: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Docker API retornou status %d: %s", resp.StatusCode, string(body))
	}

	var rawList []RawDockerContainer
	if err := json.NewDecoder(resp.Body).Decode(&rawList); err != nil {
		return nil, fmt.Errorf("falha ao decodificar lista de containers: %w", err)
	}

	result := make([]Container, 0, len(rawList))
	for _, raw := range rawList {
		c := normalizeContainer(raw)

		// Buscar variáveis de ambiente e restart policy detalhados via inspect
		if inspect, err := InspectContainer(client, raw.ID); err == nil && inspect != nil {
			c.Env = inspect.Config.Env
			if inspect.HostConfig.RestartPolicy.Name != "" {
				c.RestartPolicy = inspect.HostConfig.RestartPolicy.Name
			}
		}

		result = append(result, c)
	}

	return result, nil
}

// InspectContainer obtém detalhes adicionais de um container específico via GET /containers/{id}/json
func InspectContainer(client *connection.Client, id string) (*RawDockerInspect, error) {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return nil, fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	endpoint := fmt.Sprintf("http://localhost/containers/%s/json", url.PathEscape(id))
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status inspect: %d", resp.StatusCode)
	}

	var inspect RawDockerInspect
	if err := json.NewDecoder(resp.Body).Decode(&inspect); err != nil {
		return nil, err
	}

	return &inspect, nil
}

// StartContainer inicia um container via POST /containers/{id}/start
func StartContainer(client *connection.Client, id string) error {
	return postContainerAction(client, id, "start", nil)
}

// StopContainer para um container via POST /containers/{id}/stop?t=10
func StopContainer(client *connection.Client, id string) error {
	return postContainerAction(client, id, "stop?t=10", nil)
}

// RestartContainer reinicia um container via POST /containers/{id}/restart?t=10
func RestartContainer(client *connection.Client, id string) error {
	return postContainerAction(client, id, "restart?t=10", nil)
}

// RemoveContainer remove um container via DELETE /containers/{id}?v=1&force=true
func RemoveContainer(client *connection.Client, id string, force bool) error {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	endpoint := fmt.Sprintf("http://localhost/containers/%s?v=1&force=%t", url.PathEscape(id), force)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao remover container: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao remover container (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func CreateContainer(client *connection.Client, input CreateContainerInput) error {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}
	config := map[string]interface{}{"Image": input.Image}
	if input.Name != "" {
		config["Hostname"] = input.Name
	}
	if len(input.Env) > 0 {
		env := []string{}
		for _, item := range input.Env {
			if strings.TrimSpace(item.Name) != "" {
				env = append(env, item.Name+"="+item.Value)
			}
		}
		if len(env) > 0 {
			config["Env"] = env
		}
	}
	if input.Command != "" {
		config["Cmd"] = strings.Fields(input.Command)
	}
	hostConfig := map[string]interface{}{}
	if input.RestartPolicy != "" {
		hostConfig["RestartPolicy"] = map[string]interface{}{"Name": input.RestartPolicy}
	}
	binding := map[string][]map[string]string{}
	exposed := map[string]struct{}{}
	for _, p := range input.Ports {
		privatePort, err := strconv.Atoi(strings.TrimSpace(p.Internal))
		publicPort, publicErr := strconv.Atoi(strings.TrimSpace(p.External))
		if err != nil || privatePort <= 0 {
			continue
		}
		key := fmt.Sprintf("%d/tcp", privatePort)
		exposed[key] = struct{}{}
		if publicErr == nil && publicPort > 0 {
			binding[key] = []map[string]string{{"HostPort": fmt.Sprintf("%d", publicPort), "HostIp": ""}}
		}
	}
	if len(exposed) > 0 {
		config["ExposedPorts"] = exposed
		hostConfig["PortBindings"] = binding
	}
	if len(input.Volumes) > 0 {
		binds := []string{}
		for _, v := range input.Volumes {
			if strings.TrimSpace(v) != "" {
				binds = append(binds, v)
			}
		}
		if len(binds) > 0 {
			hostConfig["Binds"] = binds
		}
	}
	config["HostConfig"] = hostConfig
	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("falha ao montar configuração: %w", err)
	}
	endpoint := "http://localhost/containers/create?name=" + url.QueryEscape(input.Name)
	resp, err := httpClient.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("falha ao criar container: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Docker recusou a criação (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	var created struct {
		ID string `json:"Id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return fmt.Errorf("falha ao ler resposta do Docker: %w", err)
	}
	if input.Network != "" {
		networkBody, _ := json.Marshal(map[string]string{"Container": created.ID})
		networkURL := "http://localhost/networks/" + url.PathEscape(input.Network) + "/connect"
		networkResp, err := httpClient.Post(networkURL, "application/json", bytes.NewReader(networkBody))
		if err != nil {
			return fmt.Errorf("falha ao conectar container à rede: %w", err)
		}
		defer networkResp.Body.Close()
		if networkResp.StatusCode != http.StatusOK && networkResp.StatusCode != http.StatusCreated && networkResp.StatusCode != http.StatusNoContent {
			data, _ := io.ReadAll(networkResp.Body)
			return fmt.Errorf("falha ao conectar à rede (HTTP %d): %s", networkResp.StatusCode, strings.TrimSpace(string(data)))
		}
	}
	if err := postContainerAction(client, created.ID, "start", nil); err != nil {
		return err
	}
	return nil
}

// GetContainerLogs obtém os logs recentes do container via GET /containers/{id}/logs
func GetContainerLogs(client *connection.Client, id string, tail int) (string, error) {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return "", fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	if tail <= 0 {
		tail = 200
	}

	endpoint := fmt.Sprintf("http://localhost/containers/%s/logs?stdout=1&stderr=1&tail=%d&timestamps=0", url.PathEscape(id), tail)
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("falha ao consultar logs: %w", err)
	}
	defer resp.Body.Close()

	rawLogs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("falha ao ler logs: %w", err)
	}

	// O Docker multiplexa logs com header de 8 bytes (1 byte stream type + 3 bytes padding + 4 bytes tamanho)
	cleanedLogs := stripDockerLogHeaders(rawLogs)
	return cleanedLogs, nil
}

// StreamEvents escuta eventos de containers em tempo real via GET /events
func StreamEvents(ctx context.Context, client *connection.Client, eventType string, onEvent func(event DockerEvent)) error {
	httpClient := client.GetStreamHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	endpoint := "http://localhost/events"
	if eventType != "" {
		endpoint = fmt.Sprintf("http://localhost/events?filters=%%7B%%22type%%22:%%5B%%22%s%%22%%5D%%7D", url.QueryEscape(eventType))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao abrir stream de eventos do Docker: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Docker /events retornou status %d: %s", resp.StatusCode, string(body))
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var ev DockerEvent
			if err := decoder.Decode(&ev); err != nil {
				if err == io.EOF || ctx.Err() != nil {
					return nil
				}
				return fmt.Errorf("erro na leitura do stream de eventos: %w", err)
			}
			onEvent(ev)
		}
	}
}

func postContainerAction(client *connection.Client, id, action string, body io.Reader) error {
	httpClient := client.GetHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	endpoint := fmt.Sprintf("http://localhost/containers/%s/%s", url.PathEscape(id), action)
	resp, err := httpClient.Post(endpoint, "application/json", body)
	if err != nil {
		return fmt.Errorf("falha ao executar ação '%s' no container: %w", action, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ação '%s' falhou (HTTP %d): %s", action, resp.StatusCode, string(respBody))
	}

	return nil
}

func normalizeContainer(raw RawDockerContainer) Container {
	name := ""
	if len(raw.Names) > 0 {
		name = strings.TrimPrefix(raw.Names[0], "/")
	}

	// Normalizar Portas
	var portStrs []string
	var portList []PortMapping
	for _, p := range raw.Ports {
		portList = append(portList, PortMapping{
			IP:          p.IP,
			PrivatePort: p.PrivatePort,
			PublicPort:  p.PublicPort,
			Type:        p.Type,
		})

		if p.PublicPort > 0 {
			ip := p.IP
			if ip == "" {
				ip = "0.0.0.0"
			}
			portStrs = append(portStrs, fmt.Sprintf("%s:%d->%d/%s", ip, p.PublicPort, p.PrivatePort, p.Type))
		} else {
			portStrs = append(portStrs, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
		}
	}

	// Normalizar Redes
	networks := make(map[string]NetworkEndpoint)
	for netName, netData := range raw.NetworkSettings.Networks {
		networks[netName] = NetworkEndpoint{
			NetworkID:  netData.NetworkID,
			IPAddress:  netData.IPAddress,
			Gateway:    netData.Gateway,
			MacAddress: netData.MacAddress,
		}
	}

	// Normalizar Montagens
	mounts := make([]MountInfo, 0, len(raw.Mounts))
	for _, m := range raw.Mounts {
		mounts = append(mounts, MountInfo{
			Type:        m.Type,
			Name:        m.Name,
			Source:      m.Source,
			Destination: m.Destination,
			ReadOnly:    !m.RW,
		})
	}

	return Container{
		ID:            raw.ID,
		Names:         raw.Names,
		Name:          name,
		Image:         raw.Image,
		ImageID:       raw.ImageID,
		Command:       raw.Command,
		Created:       raw.Created,
		State:         raw.State,
		Status:        raw.Status,
		Ports:         strings.Join(portStrs, ", "),
		PortList:      portList,
		Networks:      networks,
		Mounts:        mounts,
		RestartPolicy: raw.HostConfig.RestartPolicy.Name,
		Labels:        raw.Labels,
	}
}

func stripDockerLogHeaders(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}

	var buf bytes.Buffer
	r := bytes.NewReader(raw)
	header := make([]byte, 8)

	for {
		n, err := io.ReadFull(r, header)
		if err != nil {
			if n > 0 {
				buf.Write(header[:n])
			}
			break
		}

		payloadSize := int(header[4])<<24 | int(header[5])<<16 | int(header[6])<<8 | int(header[7])
		if payloadSize <= 0 || payloadSize > len(raw) {
			buf.Write(header)
			_, _ = io.Copy(&buf, r)
			break
		}

		payload := make([]byte, payloadSize)
		_, err = io.ReadFull(r, payload)
		buf.Write(payload)
		if err != nil {
			break
		}
	}

	return buf.String()
}
