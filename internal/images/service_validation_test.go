package images

import (
	"go-walis/internal/core/db"
	"testing"
)

func TestPullImageRejectsEmptyNameBeforeConnection(t *testing.T) {
	result := NewImageService(nil).PullImage(db.VpsServer{}, "", "default")
	if result.Success || result.Message == "" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
