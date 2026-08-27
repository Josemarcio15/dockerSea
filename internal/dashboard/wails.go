package dashboard

import "go-walis/internal/core/db"

// WailsFacade documents the public dashboard surface exposed by Service.
type WailsFacade interface {
	ListServers() ([]db.VpsServer, error)
}

var _ WailsFacade = (*Service)(nil)
