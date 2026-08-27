package servers

import "strings"

func NormalizeServerName(name string) string { return strings.TrimSpace(name) }

func ValidateServer(server Server) string {
	if NormalizeServerName(server.Name) == "" {
		return "o nome do servidor é obrigatório"
	}
	if server.ConnectionType == "ssh" && strings.TrimSpace(server.Host) == "" {
		return "o host/IP é obrigatório para conexões SSH"
	}
	return ""
}
