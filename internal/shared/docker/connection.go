package docker

import (
	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
)

// NewClient obtém ou reutiliza uma conexão ativa no pool de conexões
func NewClient(server db.VpsServer) (*Client, error) {
	return connection.GetManager().GetClient(server)
}

// CloseClient encerra e remove do pool uma conexão de servidor
func CloseClient(serverID string) {
	connection.GetManager().CloseClient(serverID)
}
