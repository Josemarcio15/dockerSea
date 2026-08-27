package containers

import "strings"

func NormalizeContainerName(name string) string {
	return strings.TrimSpace(name)
}

func ValidateContainerAction(action string) string {
	switch action {
	case "start", "stop", "restart", "rm":
		return ""
	default:
		return "ação de container desconhecida"
	}
}
