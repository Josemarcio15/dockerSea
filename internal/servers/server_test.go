package servers

import (
	"go-walis/internal/core/db"
	"testing"
)

func TestServerValidation(t *testing.T) {
	if ValidateServer(db.VpsServer{Name: "x", ConnectionType: "ssh"}) == "" {
		t.Fatal("SSH host should be required")
	}
	if ValidateServer(db.VpsServer{Name: "x", ConnectionType: "local"}) != "" {
		t.Fatal("local server should be valid")
	}
}
