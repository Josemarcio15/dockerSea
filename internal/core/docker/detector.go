package docker

import (
	"fmt"
	"strings"

	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
)

type DetectResult struct {
	Success           bool   `json:"success"`
	Message           string `json:"message"`
	DockerSocketPath  string `json:"dockerSocketPath"`
	DockerPath        string `json:"dockerPath"`
	DockerComposePath string `json:"dockerComposePath"`
	DockerVersion     string `json:"dockerVersion"`
}

// AutoDetectEnvironment descobre dinamicamente os caminhos e versões do Docker no servidor
func AutoDetectEnvironment(server db.VpsServer) DetectResult {
	client, err := connection.NewClient(server)
	if err != nil {
		return DetectResult{
			Success: false,
			Message: fmt.Sprintf("Falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	res := DetectResult{Success: true}

	// 1. Detectar Binário do Docker (Systemd + PATH + Diretórios Padrão)
	dockerPathOut, err := client.ExecCommand(connection.CmdDiscoverDockerBin, false)
	if err == nil && strings.TrimSpace(dockerPathOut) != "" {
		res.DockerPath = strings.TrimSpace(dockerPathOut)
	} else {
		res.DockerPath = "/usr/bin/docker"
	}

	// 2. Detectar Socket do Docker (Systemd ListenStream + DOCKER_HOST + Rootless + Sockets Padrão)
	socketOut, err := client.ExecCommand(connection.CmdDiscoverSocket, false)
	if err == nil && strings.TrimSpace(socketOut) != "" {
		res.DockerSocketPath = strings.TrimSpace(socketOut)
	} else {
		res.DockerSocketPath = "/var/run/docker.sock"
	}

	// 3. Detectar Docker Compose (Plugin V2 ou Standalone V1)
	composeOut, err := client.ExecCommand(connection.CmdDiscoverCompose, false)
	if err == nil && strings.TrimSpace(composeOut) != "" {
		res.DockerComposePath = strings.TrimSpace(composeOut)
	} else {
		res.DockerComposePath = "docker compose"
	}

	// 4. Capturar Versão do Docker
	versionOut, err := client.ExecCommand(connection.CmdCheckDockerVersion, false)
	if err == nil && strings.TrimSpace(versionOut) != "" {
		res.DockerVersion = strings.TrimSpace(versionOut)
	}

	res.Message = "Ambiente e caminhos do Docker autodetectados com sucesso!"
	return res
}
