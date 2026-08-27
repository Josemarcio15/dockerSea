package servers

import (
	"go-walis/internal/core/db"
	"testing"
)

type serverRepoFake struct {
	saved   db.VpsServer
	deleted string
	active  string
}

func (f *serverRepoFake) ListServers() ([]db.VpsServer, error) { return nil, nil }
func (f *serverRepoFake) SaveServer(s db.VpsServer) error      { f.saved = s; return nil }
func (f *serverRepoFake) DeleteServer(id string) error         { f.deleted = id; return nil }
func (f *serverRepoFake) SetActiveServer(id string) error      { f.active = id; return nil }

func TestServiceSaveServerNormalizesBeforeRepository(t *testing.T) {
	repo := &serverRepoFake{}
	if err := NewServiceWithRepository(repo).SaveServer(db.VpsServer{Name: "  local  ", ConnectionType: "local"}); err != nil {
		t.Fatal(err)
	}
	if repo.saved.Name != "local" {
		t.Fatalf("saved name = %q", repo.saved.Name)
	}
}

func TestServiceDelegatesServerLifecycle(t *testing.T) {
	repo := &serverRepoFake{}
	service := NewServiceWithRepository(repo)
	if err := service.SetActiveServer("active"); err != nil {
		t.Fatal(err)
	}
	if repo.active != "active" {
		t.Fatalf("active id = %q", repo.active)
	}
	if err := service.DeleteServer("removed"); err != nil {
		t.Fatal(err)
	}
	if repo.deleted != "removed" {
		t.Fatalf("deleted id = %q", repo.deleted)
	}
}
