package networks

import (
	"go-walis/internal/core/db"
	"testing"
)

func TestCreateNetworkRejectsEmptyNameBeforeConnection(t *testing.T) {
	result := NewNetworkService(nil).CreateNetwork(db.VpsServer{}, NetworkCreateRequest{})
	if result.Success || result.Message == "" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
