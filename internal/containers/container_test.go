package containers

import "testing"

func TestContainerRules(t *testing.T) {
	if NormalizeContainerName("  web ") != "web" {
		t.Fatal("container name was not normalized")
	}
	if ValidateContainerAction("start") != "" {
		t.Fatal("start should be valid")
	}
	if ValidateContainerAction("destroy") == "" {
		t.Fatal("unknown action should be rejected")
	}
}
