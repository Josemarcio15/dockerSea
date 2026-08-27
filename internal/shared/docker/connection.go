package docker

import (
	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
)

func NewClient(server db.VpsServer) (*Client, error) { return connection.NewClient(server) }
