package volumes

import "testing"

func TestNormalizeVolumeCreateRequest(t *testing.T) {
	req := NormalizeVolumeCreateRequest(VolumeCreateRequest{Name: "  data ", Driver: " local "})
	if req.Name != "data" || req.Driver != "local" {
		t.Fatalf("request was not normalized: %#v", req)
	}
}

func TestValidateVolumeCreateRequest(t *testing.T) {
	if got := ValidateVolumeCreateRequest(VolumeCreateRequest{}); got == "" {
		t.Fatal("expected missing-name validation error")
	}
	if got := ValidateVolumeCreateRequest(VolumeCreateRequest{Name: "data"}); got != "" {
		t.Fatalf("unexpected validation error: %s", got)
	}
}
