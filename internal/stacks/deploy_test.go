package stacks

import (
	"strings"
	"sync"
	"testing"
)

type mockEventEmitter struct {
	mu     sync.Mutex
	events []struct {
		Event string
		Data  interface{}
	}
}

func (m *mockEventEmitter) Emit(event string, data interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, struct {
		Event string
		Data  interface{}
	}{Event: event, Data: data})
}

func (m *mockEventEmitter) GetEvents() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	var names []string
	for _, e := range m.events {
		names = append(names, e.Event)
	}
	return names
}

func TestGenerateDeployID(t *testing.T) {
	id1 := GenerateDeployID()
	id2 := GenerateDeployID()

	if !strings.HasPrefix(id1, "dep_") {
		t.Fatalf("esperado prefixo 'dep_', obtido %s", id1)
	}
	if id1 == id2 {
		t.Fatalf("deploy IDs sucessivos devem ser únicos: %s == %s", id1, id2)
	}
}

func TestBuildRemoteDeployDir(t *testing.T) {
	tests := []struct {
		stackID   string
		deployID  string
		expectErr bool
		expected  string
	}{
		{
			stackID:   "stk_123",
			deployID:  "dep_456",
			expectErr: false,
			expected:  "$HOME/.docksea/deploys/stk_123/dep_456",
		},
		{
			stackID:   "stk-valid_id",
			deployID:  "dep_999",
			expectErr: false,
			expected:  "$HOME/.docksea/deploys/stk-valid_id/dep_999",
		},
		{
			stackID:   "../malicious",
			deployID:  "dep_1",
			expectErr: true,
		},
		{
			stackID:   "stk_1",
			deployID:  "../../escape",
			expectErr: true,
		},
		{
			stackID:   "",
			deployID:  "dep_1",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		res, err := BuildRemoteDeployDir(tt.stackID, tt.deployID)
		if tt.expectErr {
			if err == nil {
				t.Errorf("esperado erro para stackID=%s, deployID=%s", tt.stackID, tt.deployID)
			}
		} else {
			if err != nil {
				t.Errorf("erro inesperado: %v", err)
			}
			if res != tt.expected {
				t.Errorf("esperado %s, obtido %s", tt.expected, res)
			}
		}
	}
}

func TestDeployEventBroadcaster(t *testing.T) {
	mock := &mockEventEmitter{}
	broadcaster := NewDeployEventBroadcaster("stk_test", "dep_test", mock)

	broadcaster.EmitStarted("meu-projeto")
	broadcaster.EmitProgress(PhasePreparing, "Preparando deploy...")
	broadcaster.EmitProgress(PhaseBuilding, "Compilando imagem...\r\n")
	broadcaster.EmitComplete("Deploy finalizado com sucesso!")

	events := mock.GetEvents()
	expectedEvents := []string{
		EventDeployStarted,
		EventDeployProgress,
		EventDeployProgress,
		EventDeployComplete,
	}

	if len(events) != len(expectedEvents) {
		t.Fatalf("esperado %d eventos, obtido %d: %v", len(expectedEvents), len(events), events)
	}

	for i, name := range expectedEvents {
		if events[i] != name {
			t.Errorf("evento[%d] esperado %s, obtido %s", i, name, events[i])
		}
	}
}

func TestDeployEventBroadcasterFailed(t *testing.T) {
	mock := &mockEventEmitter{}
	broadcaster := NewDeployEventBroadcaster("stk_test", "dep_test", mock)

	broadcaster.EmitStarted("meu-projeto")
	broadcaster.EmitFailed(PhaseValidating, "Arquivo docker-compose.yml inválido")

	events := mock.GetEvents()
	if len(events) != 2 {
		t.Fatalf("esperado 2 eventos, obtido %d", len(events))
	}
	if events[1] != EventDeployFailed {
		t.Errorf("esperado %s, obtido %s", EventDeployFailed, events[1])
	}
}
