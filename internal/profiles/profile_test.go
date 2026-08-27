package profiles

import "testing"

func TestProfileNameRules(t *testing.T) {
	if NormalizeProfileName("  pessoal ") != "pessoal" {
		t.Fatal("profile name was not normalized")
	}
	if ValidateProfileName("") == "" || ValidateProfileName("pessoal") != "" {
		t.Fatal("unexpected profile validation result")
	}
}
