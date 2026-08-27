package profiles

import (
	"fmt"
	"go-walis/internal/core/db"
)

// Service coordinates profile persistence and activation.
type Service struct {
	repository Repository
}

func NewService(database *db.DB) *Service {
	return NewServiceWithRepository(newDatabaseRepository(database))
}

func NewServiceWithRepository(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) ListProfiles() ([]db.Profile, error) { return s.repository.ListProfiles() }

func (s *Service) GetActiveProfile() (*db.Profile, error) {
	return s.repository.GetActiveProfile()
}

func (s *Service) SaveProfile(profile db.Profile) error {
	profile.Name = NormalizeProfileName(profile.Name)
	if message := ValidateProfileName(profile.Name); message != "" {
		return fmt.Errorf("%s", message)
	}
	return s.repository.SaveProfile(profile)
}

func (s *Service) DeleteProfile(id string) error {
	if id == "" {
		return fmt.Errorf("id do perfil inválido")
	}
	return s.repository.DeleteProfile(id)
}

func (s *Service) SetActiveProfile(id string) error {
	if id == "" {
		return fmt.Errorf("id do perfil inválido")
	}
	return s.repository.SetActiveProfile(id)
}

func (s *Service) SetProfileLocale(id string, locale string) error {
	if id == "" || locale == "" {
		return fmt.Errorf("parâmetros inválidos")
	}
	return s.repository.UpdateProfileLocale(id, locale)
}
