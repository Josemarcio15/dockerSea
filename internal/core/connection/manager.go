package connection

import (
	"sync"

	"go-walis/internal/core/db"
)

// Manager gerencia instâncias de conexões ativas reutilizáveis
type Manager struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

var globalManager = &Manager{
	clients: make(map[string]*Client),
}

// GetManager retorna a instância global do gerenciador de conexões
func GetManager() *Manager {
	return globalManager
}

// GetClient obtém ou cria uma nova conexão para o servidor reaproveitando a sessão existente
func (m *Manager) GetClient(server db.VpsServer) (*Client, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client, exists := m.clients[server.ID]; exists {
		// Se os dados essenciais de conexão mudaram, fecha e recria
		if client.server.Host == server.Host &&
			client.server.Port == server.Port &&
			client.server.Username == server.Username &&
			client.server.SshPassword == server.SshPassword &&
			client.server.SshKeyPath == server.SshKeyPath &&
			client.server.DockerSocketPath == server.DockerSocketPath &&
			client.server.SshKeyPassphrase == server.SshKeyPassphrase &&
			client.server.ConnectionType == server.ConnectionType &&
			client.IsAlive() {
			return client, nil
		}
		client.ForceClose()
		delete(m.clients, server.ID)
	}

	client, err := NewClient(server)
	if err != nil {
		return nil, err
	}

	m.clients[server.ID] = client
	return client, nil
}

// CloseClient remove e fecha a conexão de um servidor específico
func (m *Manager) CloseClient(serverID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client, exists := m.clients[serverID]; exists {
		client.ForceClose()
		delete(m.clients, serverID)
	}
}
