package volumes

import (
	"go-walis/internal/core/db"
	"testing"
)

func TestCreateVolumeRejectsEmptyNameBeforeConnection(t *testing.T) {
	result := NewVolumeService(nil).CreateVolume(db.VpsServer{}, VolumeCreateRequest{})
	if result.Success || result.Message == "" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
