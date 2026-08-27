package profiles

import "strings"

func NormalizeProfileName(name string) string { return strings.TrimSpace(name) }

func ValidateProfileName(name string) string {
	if NormalizeProfileName(name) == "" {
		return "o nome do perfil é obrigatório"
	}
	return ""
}
