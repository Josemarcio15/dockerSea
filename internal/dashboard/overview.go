package dashboard

import "go-walis/internal/core/connection"

type Overview struct {
	ServerName string                  `json:"serverName"`
	Usage      *connection.SystemUsage `json:"usage"`
}

func BuildOverview(serverName string, usage *connection.SystemUsage) Overview {
	return Overview{ServerName: serverName, Usage: usage}
}
