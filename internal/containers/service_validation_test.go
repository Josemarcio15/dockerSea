package containers

import (
	"go-walis/internal/core/db"
	"testing"
)

func TestExecuteActionRejectsUnknownActionBeforeConnection(t *testing.T) {
	result := NewContainerService(nil).ExecuteAction(db.VpsServer{}, "destroy", nil)
	if result.Success || result.Message == "" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
