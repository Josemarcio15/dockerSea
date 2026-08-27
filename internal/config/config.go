package config

import (
	"go-walis/internal/core/db"
	"go-walis/internal/servers"
)

func NormalizeServerName(name string) string { return servers.NormalizeServerName(name) }

func ValidateServerInput(name, connectionType, host string) string {
	return servers.ValidateServer(db.VpsServer{Name: name, ConnectionType: connectionType, Host: host})
}
