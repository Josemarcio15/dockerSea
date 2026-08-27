package builder

import "strings"

func NormalizeProjectName(name string) string { return strings.TrimSpace(name) }

func ValidateProjectName(name string) string {
	if NormalizeProjectName(name) == "" {
		return "nome/tag da imagem não informado"
	}
	return ""
}
