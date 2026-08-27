package images

import "strings"

func NormalizeImageName(name string) string {
	return strings.TrimSpace(name)
}

func ValidateImageName(name string) string {
	if NormalizeImageName(name) == "" {
		return "Nome da imagem não informado"
	}
	return ""
}
