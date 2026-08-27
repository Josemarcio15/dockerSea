package stacks

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var validProjectNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

// ValidateProjectName verifica se o nome do projeto obedece ao padrão do Docker Compose
func ValidateProjectName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fmt.Errorf("nome do projeto não pode ser vazio")
	}
	if len(trimmed) > 63 {
		return fmt.Errorf("nome do projeto muito longo (máximo 63 caracteres)")
	}
	if !validProjectNameRegex.MatchString(trimmed) {
		return fmt.Errorf("nome do projeto inválido '%s': deve conter apenas letras minúsculas, números, hífens ou underlines, e começar com alfanumérico", name)
	}
	return nil
}

// ShellQuote faz o quoting seguro de argumentos para evitar injeção de shell
func ShellQuote(arg string) string {
	if arg == "" {
		return "''"
	}
	// Se começar com $HOME/, expande $HOME e dá quote no restante
	if strings.HasPrefix(arg, "$HOME/") {
		rest := arg[6:]
		return "$HOME/" + ShellQuote(rest)
	}
	if arg == "$HOME" {
		return "$HOME"
	}

	// Se contém apenas caracteres alfanuméricos seguros e pontuação comum sem espaços, não precisa de quotes
	safe := true
	for _, r := range arg {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' || r == '/' || r == ':' || r == '@') {
			safe = false
			break
		}
	}
	if safe {
		return arg
	}
	// Substitui ' por '\'' e encapsula entre '...'
	return "'" + strings.ReplaceAll(arg, "'", "'\\''") + "'"
}

// ComposeVolumeMount representa uma montagem de volume no JSON do Docker Compose
type ComposeVolumeMount struct {
	Type     string `json:"type"`   // 'bind', 'volume', 'tmpfs'
	Source   string `json:"source"` // caminho no host ou nome do volume
	Target   string `json:"target"` // caminho no container
	ReadOnly bool   `json:"read_only"`
}

// ComposeServiceConfig representa a estrutura de um serviço no output do compose config
type ComposeServiceConfig struct {
	Image       string               `json:"image,omitempty"`
	Volumes     []ComposeVolumeMount `json:"volumes,omitempty"`
	Environment map[string]*string   `json:"environment,omitempty"`
}

// ComposeNormalizedConfig representa a raiz do JSON retornado por `docker compose config --format json`
type ComposeNormalizedConfig struct {
	Name     string                          `json:"name"`
	Services map[string]ComposeServiceConfig `json:"services"`
	Volumes  map[string]interface{}          `json:"volumes,omitempty"`
	Networks map[string]interface{}          `json:"networks,omitempty"`
}

// ParseComposeConfigOutput faz o parse do JSON do docker compose config retornado pela VPS
func ParseComposeConfigOutput(jsonOutput string) (*ComposeNormalizedConfig, error) {
	trimmed := strings.TrimSpace(jsonOutput)
	if trimmed == "" {
		return nil, fmt.Errorf("saída de configuração do compose está vazia")
	}

	var config ComposeNormalizedConfig
	if err := json.Unmarshal([]byte(trimmed), &config); err != nil {
		return nil, fmt.Errorf("falha ao interpretar JSON do compose config: %w", err)
	}

	return &config, nil
}

// HasRuntimeBindDependency verifica se o Compose possui bind mounts cujo source aponte para dentro do remoteDeployDir.
// Se apontar para dentro do remoteDeployDir, o diretório de deploy na VPS NÃO pode ser apagado após o deploy.
func HasRuntimeBindDependency(config *ComposeNormalizedConfig, remoteDeployDir string) bool {
	if config == nil || remoteDeployDir == "" {
		return false
	}

	cleanDeployDir := filepath.Clean(remoteDeployDir)

	for _, service := range config.Services {
		for _, v := range service.Volumes {
			if strings.EqualFold(v.Type, "bind") && v.Source != "" {
				cleanSource := filepath.Clean(v.Source)
				// Verifica se cleanSource é o próprio deployDir ou está dentro dele
				if cleanSource == cleanDeployDir || strings.HasPrefix(cleanSource, cleanDeployDir+"/") || strings.HasPrefix(cleanSource, cleanDeployDir+"\\") {
					return true
				}
			}
		}
	}

	return false
}

// BuildComposeConfigCommand gera o comando remoto para validar e extrair o JSON
func BuildComposeConfigCommand(remoteDir, composeBin, projectName string, asJSON bool) string {
	if composeBin == "" {
		composeBin = "docker compose"
	}
	formatFlag := ""
	if asJSON {
		formatFlag = " --format json"
	}
	return fmt.Sprintf("cd %s && %s -p %s config%s", ShellQuote(remoteDir), composeBin, ShellQuote(projectName), formatFlag)
}

// BuildComposeUpCommand gera o comando remoto para subir os containers com build e wait
func BuildComposeUpCommand(remoteDir, composeBin, projectName string) string {
	if composeBin == "" {
		composeBin = "docker compose"
	}
	return fmt.Sprintf("cd %s && %s -p %s up -d --build --remove-orphans --wait", ShellQuote(remoteDir), composeBin, ShellQuote(projectName))
}

// BuildComposeStopCommand gera o comando remoto para parar a stack
func BuildComposeStopCommand(remoteDir, composeBin, projectName string) string {
	if composeBin == "" {
		composeBin = "docker compose"
	}
	return fmt.Sprintf("cd %s && %s -p %s stop", ShellQuote(remoteDir), composeBin, ShellQuote(projectName))
}

// BuildComposeDownCommand gera o comando remoto para down com ou sem remoção de volumes
func BuildComposeDownCommand(remoteDir, composeBin, projectName string, deleteVolumes bool) string {
	if composeBin == "" {
		composeBin = "docker compose"
	}
	volFlag := ""
	if deleteVolumes {
		volFlag = " -v"
	}
	return fmt.Sprintf("cd %s && %s -p %s down%s --remove-orphans", ShellQuote(remoteDir), composeBin, ShellQuote(projectName), volFlag)
}
