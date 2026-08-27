package images

import "testing"

func TestImageNameRules(t *testing.T) {
	if got := NormalizeImageName("  nginx:latest "); got != "nginx:latest" {
		t.Fatalf("got %q", got)
	}
	if ValidateImageName("") == "" || ValidateImageName("nginx") != "" {
		t.Fatal("unexpected image-name validation result")
	}
}
