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

// GetClient obtém ou cria uma nova conexão para o servidor
func (m *Manager) GetClient(server db.VpsServer) (*Client, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client, exists := m.clients[server.ID]; exists {
		if client.IsAlive() {
			return client, nil
		}
		client.Close()
		delete(m.clients, server.ID)
	}

	client, err := NewClient(server)
	if err != nil {
		return nil, err
	}

	m.clients[server.ID] = client
	return client, nil
}
