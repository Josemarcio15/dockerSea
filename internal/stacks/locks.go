package stacks

import (
	"sync"
)

type lockEntry struct {
	mu       sync.Mutex
	refCount int
}

// StackLockManager gerencia exclusão mútua em memória por chave (ex: profileID + ":" + stackID)
// com contagem de referências para evitar vazamento de memória com o tempo.
type StackLockManager struct {
	mu    sync.Mutex
	locks map[string]*lockEntry
}

// NewStackLockManager cria uma nova instância do gerenciador de locks
func NewStackLockManager() *StackLockManager {
	return &StackLockManager{
		locks: make(map[string]*lockEntry),
	}
}

// Lock adquire um lock para a chave especificada e retorna a função unlock correspondente.
// A execução é thread-safe e serializa operações na mesma chave, enquanto permite
// concorrência total entre chaves diferentes.
func (m *StackLockManager) Lock(key string) func() {
	m.mu.Lock()
	entry, exists := m.locks[key]
	if !exists {
		entry = &lockEntry{}
		m.locks[key] = entry
	}
	entry.refCount++
	m.mu.Unlock()

	entry.mu.Lock()

	return func() {
		entry.mu.Unlock()

		m.mu.Lock()
		entry.refCount--
		if entry.refCount <= 0 {
			delete(m.locks, key)
		}
		m.mu.Unlock()
	}
}
