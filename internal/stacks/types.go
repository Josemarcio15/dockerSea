package stacks

import (
	"time"
)

// StackSourceType define os tipos de origem da stack
type StackSourceType string

const (
	SourceTypeEditor StackSourceType = "editor"
	SourceTypeFolder StackSourceType = "folder"
)

// Stack representa uma stack no banco de dados e no domínio
type Stack struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	ProjectName string          `json:"projectName"`
	SourceType  StackSourceType `json:"sourceType"` // 'editor' | 'folder'
	FolderPath  string          `json:"folderPath"`
	YamlContent string          `json:"yamlContent"`
	ProfileID   string          `json:"profileId"`

	// Último deploy bem-sucedido
	LastDeployedYAML       string     `json:"lastDeployedYaml,omitempty"`
	LastDeployedConfigYAML string     `json:"lastDeployedConfigYaml,omitempty"`
	LastDeployedConfigJSON string     `json:"lastDeployedConfigJson,omitempty"`
	LastDeployedRemoteDir  string     `json:"lastDeployedRemoteDir,omitempty"`
	LastDeployedAt         *time.Time `json:"lastDeployedAt,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// StackItem é a representação DTO utilizada pelo frontend Wails
type StackItem struct {
	ID                    string     `json:"id"`
	Name                  string     `json:"name"`
	ProjectName           string     `json:"projectName"`
	SourceType            string     `json:"sourceType"` // 'editor' | 'folder'
	FolderPath            string     `json:"folderPath"`
	YamlContent           string     `json:"yamlContent"`
	ProfileID             string     `json:"profileId"`
	LastDeployedRemoteDir string     `json:"lastDeployedRemoteDir,omitempty"`
	LastDeployedAt        *time.Time `json:"lastDeployedAt,omitempty"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
}

// StackActionResult representa o resultado de uma operação executada
type StackActionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Logs    string `json:"logs,omitempty"`
}

// StackProgressEvent representa o payload de evento de streaming para o Wails
type StackProgressEvent struct {
	StackID  string `json:"stackId"`
	DeployID string `json:"deployId"`
	Phase    string `json:"phase"` // preparing, uploading, validating, building, starting, checking, cleaning, complete, error
	Message  string `json:"message"`
	Success  bool   `json:"success,omitempty"`
}

// ServerCapabilities representa o suporte a Docker e Compose na VPS remota
type ServerCapabilities struct {
	DockerAvailable  bool   `json:"dockerAvailable"`
	ComposeAvailable bool   `json:"composeAvailable"`
	DockerVersion    string `json:"dockerVersion"`
	ComposeVersion   string `json:"composeVersion"`
	BuildxAvailable  bool   `json:"buildxAvailable"`
	Architecture     string `json:"architecture"`
	OperatingSystem  string `json:"operatingSystem"`
}
