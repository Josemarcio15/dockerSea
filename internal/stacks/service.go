package stacks

import (
	"context"
	"fmt"
	"strings"

	"go-walis/internal/core/db"
	sharedDocker "go-walis/internal/shared/docker"
)

// StackService gerencia todas as operações de Stacks para o frontend Wails
type StackService struct {
	database *db.DB
	lockMgr  *StackLockManager
}

// NewStackService instancia o serviço de Stacks
func NewStackService(database *db.DB) *StackService {
	return &StackService{
		database: database,
		lockMgr:  NewStackLockManager(),
	}
}

// ListStacks lista as stacks gravadas no SQLite para um determinado profile
func (s *StackService) ListStacks(profileID string) ([]StackItem, error) {
	if s.database == nil {
		return nil, fmt.Errorf("banco de dados SQLite não inicializado")
	}
	dbStacks, err := s.database.ListStacks(profileID)
	if err != nil {
		return nil, err
	}
	res := make([]StackItem, 0, len(dbStacks))
	for _, item := range dbStacks {
		res = append(res, StackItem{
			ID:                    item.ID,
			Name:                  item.Name,
			ProjectName:           item.ProjectName,
			SourceType:            item.SourceType,
			FolderPath:            item.FolderPath,
			YamlContent:           item.YamlContent,
			ProfileID:             item.ProfileID,
			LastDeployedRemoteDir: item.LastDeployedRemoteDir,
			LastDeployedAt:        item.LastDeployedAt,
			CreatedAt:             item.CreatedAt,
			UpdatedAt:             item.UpdatedAt,
		})
	}
	return res, nil
}

// GetStack obtém os detalhes de uma stack cadastrada
func (s *StackService) GetStack(stackID string) (*Stack, error) {
	if s.database == nil {
		return nil, fmt.Errorf("banco de dados SQLite não inicializado")
	}
	dbStack, err := s.database.GetStack(stackID)
	if err != nil {
		return nil, err
	}
	if dbStack == nil {
		return nil, fmt.Errorf("stack não encontrada")
	}
	return &Stack{
		ID:                     dbStack.ID,
		Name:                   dbStack.Name,
		ProjectName:            dbStack.ProjectName,
		SourceType:             StackSourceType(dbStack.SourceType),
		FolderPath:             dbStack.FolderPath,
		YamlContent:            dbStack.YamlContent,
		ProfileID:              dbStack.ProfileID,
		LastDeployedYAML:       dbStack.LastDeployedYAML,
		LastDeployedConfigYAML: dbStack.LastDeployedConfigYAML,
		LastDeployedConfigJSON: dbStack.LastDeployedConfigJSON,
		LastDeployedRemoteDir:  dbStack.LastDeployedRemoteDir,
		LastDeployedAt:         dbStack.LastDeployedAt,
		CreatedAt:              dbStack.CreatedAt,
		UpdatedAt:              dbStack.UpdatedAt,
	}, nil
}

// SaveStack salva ou atualiza uma definição de stack no SQLite
func (s *StackService) SaveStack(item StackItem) error {
	if s.database == nil {
		return fmt.Errorf("banco de dados SQLite não inicializado")
	}
	if strings.TrimSpace(item.Name) == "" {
		return fmt.Errorf("nome da stack é obrigatório")
	}
	if strings.TrimSpace(item.ProjectName) == "" {
		return fmt.Errorf("nome do projeto é obrigatório")
	}
	if err := ValidateProjectName(item.ProjectName); err != nil {
		return err
	}

	return s.database.SaveStack(db.Stack{
		ID:          item.ID,
		Name:        item.Name,
		ProjectName: item.ProjectName,
		SourceType:  item.SourceType,
		FolderPath:  item.FolderPath,
		YamlContent: item.YamlContent,
		ProfileID:   item.ProfileID,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}

// DeleteStackDefinition remove exclusivamente o registro da stack no SQLite local (não altera a VPS)
func (s *StackService) DeleteStackDefinition(stackID string) error {
	if s.database == nil {
		return fmt.Errorf("banco de dados SQLite não inicializado")
	}
	return s.database.DeleteStack(stackID)
}

// DeleteStack é mantido por compatibilidade com implementações existentes (apenas SQLite)
func (s *StackService) DeleteStack(id string) error {
	return s.DeleteStackDefinition(id)
}

// DeployStack realiza o deploy seguro e streaming via SSH recebendo apenas profileID e stackID
func (s *StackService) DeployStack(profileID string, stackID string) StackActionResult {
	if s.database == nil {
		return StackActionResult{Success: false, Message: "Banco de dados SQLite não inicializado"}
	}

	stack, err := s.database.GetStack(stackID)
	if err != nil || stack == nil {
		return StackActionResult{Success: false, Message: "Stack não encontrada"}
	}

	server, err := s.resolveServer(profileID)
	if err != nil {
		return StackActionResult{Success: false, Message: fmt.Sprintf("falha ao resolver servidor VPS: %v", err)}
	}

	ctx := context.Background()
	return ExecuteDeploy(ctx, s.database, *server, stack, s.lockMgr)
}

// StopStack pausa/para a execução dos containers da stack na VPS
func (s *StackService) StopStack(profileID string, stackID string) StackActionResult {
	if s.database == nil {
		return StackActionResult{Success: false, Message: "Banco de dados SQLite não inicializado"}
	}

	stack, err := s.database.GetStack(stackID)
	if err != nil || stack == nil {
		return StackActionResult{Success: false, Message: "Stack não encontrada"}
	}

	server, err := s.resolveServer(profileID)
	if err != nil {
		return StackActionResult{Success: false, Message: fmt.Sprintf("falha ao resolver servidor VPS: %v", err)}
	}

	client, err := sharedDocker.NewClient(*server)
	if err != nil {
		return StackActionResult{Success: false, Message: fmt.Sprintf("falha na conexão com VPS: %v", err)}
	}
	defer client.Close()

	ctx := context.Background()
	return ExecuteStopStack(ctx, client, stack, server.DockerComposePath, s.lockMgr)
}

// RemoveStackRemote executa docker compose down na VPS e remove containers e redes (e volumes se solicitado)
func (s *StackService) RemoveStackRemote(profileID string, stackID string, deleteVolumes bool) StackActionResult {
	if s.database == nil {
		return StackActionResult{Success: false, Message: "Banco de dados SQLite não inicializado"}
	}

	stack, err := s.database.GetStack(stackID)
	if err != nil || stack == nil {
		return StackActionResult{Success: false, Message: "Stack não encontrada"}
	}

	server, err := s.resolveServer(profileID)
	if err != nil {
		return StackActionResult{Success: false, Message: fmt.Sprintf("falha ao resolver servidor VPS: %v", err)}
	}

	client, err := sharedDocker.NewClient(*server)
	if err != nil {
		return StackActionResult{Success: false, Message: fmt.Sprintf("falha na conexão com VPS: %v", err)}
	}
	defer client.Close()

	ctx := context.Background()
	return ExecuteRemoveStackRemote(ctx, s.database, client, stack, server.DockerComposePath, deleteVolumes, s.lockMgr)
}

// GetStackLogs obtém os logs do Docker Compose na VPS
func (s *StackService) GetStackLogs(profileID string, stackID string, tail int) (string, error) {
	if s.database == nil {
		return "", fmt.Errorf("banco de dados SQLite não inicializado")
	}

	stack, err := s.database.GetStack(stackID)
	if err != nil || stack == nil {
		return "", fmt.Errorf("stack não encontrada")
	}

	server, err := s.resolveServer(profileID)
	if err != nil {
		return "", fmt.Errorf("falha ao resolver servidor VPS: %w", err)
	}

	client, err := sharedDocker.NewClient(*server)
	if err != nil {
		return "", fmt.Errorf("falha na conexão com VPS: %w", err)
	}
	defer client.Close()

	return ExecuteGetStackLogs(client, stack, server.DockerComposePath, tail)
}

// GetServerCapabilities detecta capacidades do servidor VPS vinculado ao profile
func (s *StackService) GetServerCapabilities(profileID string) (ServerCapabilities, error) {
	var caps ServerCapabilities
	if s.database == nil {
		return caps, fmt.Errorf("banco de dados SQLite não inicializado")
	}

	server, err := s.resolveServer(profileID)
	if err != nil {
		return caps, err
	}

	client, err := sharedDocker.NewClient(*server)
	if err != nil {
		return caps, fmt.Errorf("falha na conexão com VPS: %w", err)
	}
	defer client.Close()

	return DetectServerCapabilities(context.Background(), client, server.DockerComposePath)
}

// resolveServer localiza o servidor cadastrado (preferencialmente ativo)
func (s *StackService) resolveServer(profileID string) (*db.VpsServer, error) {
	servers, err := s.database.ListVpsServers()
	if err != nil {
		return nil, err
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("nenhum servidor VPS cadastrado")
	}
	// Prioriza servidor ativo
	for _, srv := range servers {
		if srv.IsActive {
			return &srv, nil
		}
	}
	return &servers[0], nil
}
