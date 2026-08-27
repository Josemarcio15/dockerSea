package profiles

import (
	"go-walis/internal/core/db"
	"testing"
)

type profileRepoFake struct {
	saved           db.Profile
	deleted, active string
}

func (f *profileRepoFake) ListProfiles() ([]db.Profile, error)      { return nil, nil }
func (f *profileRepoFake) GetActiveProfile() (*db.Profile, error)   { return nil, nil }
func (f *profileRepoFake) SaveProfile(p db.Profile) error           { f.saved = p; return nil }
func (f *profileRepoFake) DeleteProfile(id string) error            { f.deleted = id; return nil }
func (f *profileRepoFake) SetActiveProfile(id string) error         { f.active = id; return nil }
func (f *profileRepoFake) UpdateProfileLocale(string, string) error { return nil }

func TestServiceSaveProfileNormalizesBeforeRepository(t *testing.T) {
	repo := &profileRepoFake{}
	if err := NewServiceWithRepository(repo).SaveProfile(db.Profile{Name: "  demo  "}); err != nil {
		t.Fatal(err)
	}
	if repo.saved.Name != "demo" {
		t.Fatalf("saved name = %q", repo.saved.Name)
	}
}

func TestServiceDelegatesProfileLifecycle(t *testing.T) {
	repo := &profileRepoFake{}
	service := NewServiceWithRepository(repo)
	if err := service.SetActiveProfile("active"); err != nil {
		t.Fatal(err)
	}
	if repo.active != "active" {
		t.Fatalf("active id = %q", repo.active)
	}
	if err := service.DeleteProfile("removed"); err != nil {
		t.Fatal(err)
	}
	if repo.deleted != "removed" {
		t.Fatalf("deleted id = %q", repo.deleted)
	}
}
