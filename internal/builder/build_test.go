package builder

import "testing"

func TestProjectNameRules(t *testing.T) {
	if NormalizeProjectName("  app ") != "app" {
		t.Fatal("project name was not normalized")
	}
	if ValidateProjectName("") == "" || ValidateProjectName("app") != "" {
		t.Fatal("unexpected validation result")
	}
}

func TestValidateBuildInput(t *testing.T) {
	if got := ValidateBuildInput("", "app"); got == "" {
		t.Fatal("empty folder must be rejected")
	}
	if got := ValidateBuildInput("/tmp/app", ""); got == "" {
		t.Fatal("empty project name must be rejected")
	}
	if got := ValidateBuildInput("/tmp/app", "app"); got != "" {
		t.Fatalf("valid build input rejected: %s", got)
	}
}
