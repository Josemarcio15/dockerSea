package dashboard

import (
	"fmt"
	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
	sharedDocker "go-walis/internal/shared/docker"
)

func CollectSystemUsage(server db.VpsServer) (*connection.SystemUsage, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()
	return client.FetchSystemUsage()
}
