package config

import "testing"

func TestValidateServerInput(t *testing.T) {
	if ValidateServerInput("", "local", "") == "" {
		t.Fatal("name should be required")
	}
	if ValidateServerInput("remote", "ssh", "") == "" {
		t.Fatal("SSH host should be required")
	}
	if ValidateServerInput("remote", "ssh", "10.0.0.1") != "" {
		t.Fatal("valid server rejected")
	}
}
