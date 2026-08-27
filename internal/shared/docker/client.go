package docker

import "go-walis/internal/core/connection"

// Client is the shared Docker transport used by domain adapters.
type Client = connection.Client
