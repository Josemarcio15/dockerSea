package connection

import (
	"fmt"
	"time"

	"go-walis/internal/core/db"
)

// LogConnectionStarted registra no terminal o início de uma nova conexão
func LogConnectionStarted(server db.VpsServer, startedAt time.Time) {
	connType := "SSH Remoto"
	target := fmt.Sprintf("%s@%s:%d", server.Username, server.Host, server.Port)
	if isLocal(server) {
		connType = "Local (Unix Socket)"
		target = server.DockerSocketPath
		if target == "" {
			target = "/var/run/docker.sock"
		}
	}

	fmt.Printf("[CONN LOG] [%s] Conexão INICIADA | Tipo: %s | Servidor: %s (ID: %s) | Alvo: %s\n",
		startedAt.Format("2006-01-02 15:04:05.000"),
		connType,
		server.Name,
		server.ID,
		target,
	)
}

// LogConnectionClosed registra no terminal o encerramento da conexão
func LogConnectionClosed(server db.VpsServer, startedAt time.Time, closedAt time.Time) {
	connType := "SSH Remoto"
	target := fmt.Sprintf("%s@%s:%d", server.Username, server.Host, server.Port)
	if isLocal(server) {
		connType = "Local (Unix Socket)"
		target = server.DockerSocketPath
		if target == "" {
			target = "/var/run/docker.sock"
		}
	}

	duration := closedAt.Sub(startedAt).Round(time.Millisecond)

	fmt.Printf("[CONN LOG] [%s] Conexão FINALIZADA | Tipo: %s | Servidor: %s (ID: %s) | Alvo: %s | Duração Ativa: %s\n",
		closedAt.Format("2006-01-02 15:04:05.000"),
		connType,
		server.Name,
		server.ID,
		target,
		duration,
	)
}
