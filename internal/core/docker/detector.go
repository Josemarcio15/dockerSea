package docker

import (
	"fmt"
	"strings"

	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
	sharedDocker "go-walis/internal/shared/docker"
)

type DetectResult struct {
	Success           bool     `json:"success"`
	Message           string   `json:"message"`
	DockerSocketPath  string   `json:"dockerSocketPath"`
	DockerPath        string   `json:"dockerPath"`
	DockerComposePath string   `json:"dockerComposePath"`
	DockerVersion     string   `json:"dockerVersion"`
	AvailableSockets  []string `json:"availableSockets"`
	AvailableBins     []string `json:"availableBins"`
	AvailableComposes []string `json:"availableComposes"`
}

func parseUniqueLines(output string) []string {
	seen := make(map[string]bool)
	var list []string
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !seen[trimmed] {
			seen[trimmed] = true
			list = append(list, trimmed)
		}
	}
	return list
}

// AutoDetectEnvironment descobre dinamicamente os caminhos e versões do Docker no servidor
func AutoDetectEnvironment(server db.VpsServer) DetectResult {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return DetectResult{
			Success: false,
			Message: fmt.Sprintf("Falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	res := DetectResult{Success: true}

	// 1. Detectar Todos os Binários do Docker
	dockerPathOut, _ := client.ExecCommand(connection.CmdDiscoverAllDockerBins, false)
	res.AvailableBins = parseUniqueLines(dockerPathOut)
	if len(res.AvailableBins) > 0 {
		res.DockerPath = res.AvailableBins[0]
	} else {
		res.DockerPath = "/usr/bin/docker"
		res.AvailableBins = []string{"/usr/bin/docker"}
	}

	// 2. Detectar Todos os Sockets do Docker (Rootless + Rootfull)
	socketOut, _ := client.ExecCommand(connection.CmdDiscoverAllSockets, false)
	res.AvailableSockets = parseUniqueLines(socketOut)
	if len(res.AvailableSockets) > 0 {
		res.DockerSocketPath = res.AvailableSockets[0]
	} else {
		res.DockerSocketPath = "/var/run/docker.sock"
		res.AvailableSockets = []string{"/var/run/docker.sock"}
	}

	// 3. Detectar Todos os Docker Composes
	composeOut, _ := client.ExecCommand(connection.CmdDiscoverAllComposes, false)
	res.AvailableComposes = parseUniqueLines(composeOut)
	if len(res.AvailableComposes) > 0 {
		res.DockerComposePath = res.AvailableComposes[0]
	} else {
		res.DockerComposePath = "docker compose"
		res.AvailableComposes = []string{"docker compose"}
	}

	// 4. Capturar Versão do Docker
	versionOut, err := client.ExecCommand(connection.CmdCheckDockerVersion, false)
	if err == nil && strings.TrimSpace(versionOut) != "" {
		res.DockerVersion = strings.TrimSpace(versionOut)
	}

	res.Message = fmt.Sprintf("Ambiente autodetectado! Encontrados %d socket(s) e %d binário(s).", len(res.AvailableSockets), len(res.AvailableBins))
	return res
}
