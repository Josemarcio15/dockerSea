package stacks

import (
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	validNames := []string{
		"my-project",
		"project123",
		"backend_api",
		"a",
		"stack-prod-01",
	}

	for _, name := range validNames {
		if err := ValidateProjectName(name); err != nil {
			t.Errorf("nome válido '%s' foi rejeitado: %v", name, err)
		}
	}

	invalidNames := []string{
		"",
		"MyProject",          // maiúsculas
		"-invalid-start",     // começa com hífen
		"_invalid-start",     // começa com underline
		"project with space", // espaço
		"stack/nest",         // barra
		"stack$foo",          // caracteres especiais
	}

	for _, name := range invalidNames {
		if err := ValidateProjectName(name); err == nil {
			t.Errorf("nome inválido '%s' deveria ter sido rejeitado", name)
		}
	}
}

func TestShellQuote(t *testing.T) {
	cases := map[string]string{
		"simple":                 "simple",
		"path/to/file.txt":       "path/to/file.txt",
		"with spaces":            "'with spaces'",
		"with'quote":             "'with'\\''quote'",
		"complex; rm -rf /; echo": "'complex; rm -rf /; echo'",
		"$HOME/.docksea/deploys/1": "$HOME/.docksea/deploys/1",
		"":                       "''",
	}

	for input, expected := range cases {
		quoted := ShellQuote(input)
		if quoted != expected {
			t.Errorf("ShellQuote(%q) = %q, esperado %q", input, quoted, expected)
		}
	}
}

func TestParseComposeConfigOutput_And_BindDependency(t *testing.T) {
	deployDir := "/home/user/.docksea/deploys/stk_123/dep_abc"

	// Caso 1: Com Bind Mount apontando para dentro do deployDir
	jsonWithBind := `{
		"name": "my-stack",
		"services": {
			"web": {
				"image": "nginx:alpine",
				"volumes": [
					{
						"type": "bind",
						"source": "/home/user/.docksea/deploys/stk_123/dep_abc/nginx.conf",
						"target": "/etc/nginx/nginx.conf"
					}
				]
			}
		}
	}`

	cfg1, err := ParseComposeConfigOutput(jsonWithBind)
	if err != nil {
		t.Fatalf("erro ao fazer parse do config: %v", err)
	}

	if !HasRuntimeBindDependency(cfg1, deployDir) {
		t.Errorf("esperava HasRuntimeBindDependency == true para bind mount dentro do deployDir")
	}

	// Caso 2: Com Bind Mount externo (fora do deployDir)
	jsonWithExternalBind := `{
		"name": "my-stack",
		"services": {
			"web": {
				"image": "nginx:alpine",
				"volumes": [
					{
						"type": "bind",
						"source": "/var/log/nginx",
						"target": "/var/log/nginx"
					}
				]
			}
		}
	}`

	cfg2, err := ParseComposeConfigOutput(jsonWithExternalBind)
	if err != nil {
		t.Fatalf("erro ao fazer parse do config: %v", err)
	}

	if HasRuntimeBindDependency(cfg2, deployDir) {
		t.Errorf("esperava HasRuntimeBindDependency == false para bind externo fora do deployDir")
	}

	// Caso 3: Apenas Volumes Nomeados (sem bind mounts)
	jsonWithNamedVolumes := `{
		"name": "my-stack",
		"services": {
			"db": {
				"image": "postgres:17",
				"volumes": [
					{
						"type": "volume",
						"source": "db_data",
						"target": "/var/lib/postgresql/data"
					}
				]
			}
		},
		"volumes": {
			"db_data": {}
		}
	}`

	cfg3, err := ParseComposeConfigOutput(jsonWithNamedVolumes)
	if err != nil {
		t.Fatalf("erro ao fazer parse do config: %v", err)
	}

	if HasRuntimeBindDependency(cfg3, deployDir) {
		t.Errorf("esperava HasRuntimeBindDependency == false para volumes nomeados")
	}
}
