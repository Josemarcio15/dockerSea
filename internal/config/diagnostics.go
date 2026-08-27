package config

import (
	"go-walis/internal/core/db"
	"go-walis/internal/core/docker"
)

// Diagnostics keeps environment checks outside configuration persistence.
func Diagnostics(server db.VpsServer) docker.DiagnosticResult { return docker.RunDiagnostic(server) }
