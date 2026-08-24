package stacks

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
)

// GenerateOperationID gera um identificador seguro e único para uma operação efêmera de ciclo de vida
func GenerateOperationID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("op_%d_%s", time.Now().Unix(), hex.EncodeToString(b))
}

// BuildRemoteLifecycleDir monta o caminho remoto efêmero para envio mínimo de Compose
func BuildRemoteLifecycleDir(stackID, operationID string) (string, error) {
	sID := strings.TrimSpace(stackID)
	opID := strings.TrimSpace(operationID)

	if sID == "" || !validIDRegex.MatchString(sID) {
		return "", fmt.Errorf("stack ID inválido para diretório de lifecycle: %s", stackID)
	}
	if opID == "" || !validIDRegex.MatchString(opID) {
		return "", fmt.Errorf("operation ID inválido para diretório de lifecycle: %s", operationID)
	}

	return fmt.Sprintf("$HOME/.docksea/lifecycle/%s/%s", sID, opID), nil
}

// ExecuteStopStack executa docker compose stop usando o snapshot canônico sem reenviar todo o projeto
func ExecuteStopStack(
	ctx context.Context,
	client *connection.Client,
	stack *db.Stack,
	composePath string,
	lockMgr *StackLockManager,
) StackActionResult {
	if stack == nil {
		return StackActionResult{Success: false, Message: "Stack não informada"}
	}
	if client == nil {
		return StackActionResult{Success: false, Message: "Conexão com VPS não estabelecida"}
	}

	// 1. Lock de concorrência
	lockKey := stack.ProfileID + ":" + stack.ID
	if lockMgr != nil {
		unlock := lockMgr.Lock(lockKey)
		defer unlock()
	}

	composeBin := strings.TrimSpace(composePath)
	if composeBin == "" {
		composeBin = "docker compose"
	}

	// 2. Se a stack ainda tem um diretório ativo existente (por causa de binds), podemos executar direto lá
	// Caso contrário, criamos um diretório de ciclo de vida efêmero com o Compose canônico
	var targetDir string
	var isEphemeral bool

	if stack.LastDeployedRemoteDir != "" {
		targetDir = stack.LastDeployedRemoteDir
	} else {
		yamlContent := stack.LastDeployedConfigYAML
		if strings.TrimSpace(yamlContent) == "" {
			yamlContent = stack.LastDeployedYAML
		}
		if strings.TrimSpace(yamlContent) == "" {
			yamlContent = stack.YamlContent
		}
		if strings.TrimSpace(yamlContent) == "" {
			return StackActionResult{Success: false, Message: "Nenhuma configuração Compose encontrada para parar a stack"}
		}

		opID := GenerateOperationID()
		lifecycleDir, err := BuildRemoteLifecycleDir(stack.ID, opID)
		if err != nil {
			return StackActionResult{Success: false, Message: fmt.Sprintf("erro ao gerar pasta efêmera: %v", err)}
		}

		// Cria pasta efêmera e envia o YAML canônico
		mkdirCmd := fmt.Sprintf("mkdir -p %s", lifecycleDir)
		if _, err := client.ExecCommand(mkdirCmd, false); err != nil {
			return StackActionResult{Success: false, Message: fmt.Sprintf("falha ao criar pasta efêmera na VPS: %v", err)}
		}

		extractCmd := fmt.Sprintf("tar -xzf - -C %s", lifecycleDir)
		stdinPipe, waitExtract, err := client.StartCommandInput(ctx, extractCmd, false)
		if err != nil {
			_ = cleanRemoteDir(client, lifecycleDir)
			return StackActionResult{Success: false, Message: fmt.Sprintf("falha ao abrir stream para pasta efêmera: %v", err)}
		}

		if err := PackSingleYaml(yamlContent, stdinPipe); err != nil {
			_ = stdinPipe.Close()
			_ = cleanRemoteDir(client, lifecycleDir)
			return StackActionResult{Success: false, Message: fmt.Sprintf("falha ao empacotar Compose efêmero: %v", err)}
		}
		_ = stdinPipe.Close()

		if err := waitExtract(); err != nil {
			_ = cleanRemoteDir(client, lifecycleDir)
			return StackActionResult{Success: false, Message: fmt.Sprintf("falha ao extrair Compose efêmero na VPS: %v", err)}
		}

		targetDir = lifecycleDir
		isEphemeral = true
	}

	if isEphemeral {
		defer cleanRemoteDir(client, targetDir)
	}

	cmdStop := BuildComposeStopCommand(targetDir, composeBin, stack.ProjectName)
	out, err := client.ExecCommand(cmdStop, false)
	if err != nil {
		return StackActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao parar a stack '%s': %v", stack.Name, err),
			Logs:    out,
		}
	}

	return StackActionResult{
		Success: true,
		Message: fmt.Sprintf("Stack '%s' parada com sucesso!", stack.Name),
		Logs:    out,
	}
}

// ExecuteRemoveStackRemote executa docker compose down (com ou sem remoção de volumes) e faz cleanup seguro
func ExecuteRemoveStackRemote(
	ctx context.Context,
	database *db.DB,
	client *connection.Client,
	stack *db.Stack,
	composePath string,
	deleteVolumes bool,
	lockMgr *StackLockManager,
) StackActionResult {
	if stack == nil {
		return StackActionResult{Success: false, Message: "Stack não informada"}
	}
	if client == nil {
		return StackActionResult{Success: false, Message: "Conexão com VPS não estabelecida"}
	}

	// 1. Lock de concorrência
	lockKey := stack.ProfileID + ":" + stack.ID
	if lockMgr != nil {
		unlock := lockMgr.Lock(lockKey)
		defer unlock()
	}

	composeBin := strings.TrimSpace(composePath)
	if composeBin == "" {
		composeBin = "docker compose"
	}

	// 2. Determina diretório de execução
	var targetDir string
	var isEphemeral bool

	if stack.LastDeployedRemoteDir != "" {
		targetDir = stack.LastDeployedRemoteDir
	} else {
		yamlContent := stack.LastDeployedConfigYAML
		if strings.TrimSpace(yamlContent) == "" {
			yamlContent = stack.LastDeployedYAML
		}
		if strings.TrimSpace(yamlContent) == "" {
			yamlContent = stack.YamlContent
		}
		if strings.TrimSpace(yamlContent) == "" {
			return StackActionResult{Success: false, Message: "Nenhuma configuração Compose encontrada para remover a stack"}
		}

		opID := GenerateOperationID()
		lifecycleDir, err := BuildRemoteLifecycleDir(stack.ID, opID)
		if err != nil {
			return StackActionResult{Success: false, Message: fmt.Sprintf("erro ao gerar pasta efêmera: %v", err)}
		}

		mkdirCmd := fmt.Sprintf("mkdir -p %s", lifecycleDir)
		if _, err := client.ExecCommand(mkdirCmd, false); err != nil {
			return StackActionResult{Success: false, Message: fmt.Sprintf("falha ao criar pasta efêmera na VPS: %v", err)}
		}

		extractCmd := fmt.Sprintf("tar -xzf - -C %s", lifecycleDir)
		stdinPipe, waitExtract, err := client.StartCommandInput(ctx, extractCmd, false)
		if err != nil {
			_ = cleanRemoteDir(client, lifecycleDir)
			return StackActionResult{Success: false, Message: fmt.Sprintf("falha ao abrir stream para pasta efêmera: %v", err)}
		}

		if err := PackSingleYaml(yamlContent, stdinPipe); err != nil {
			_ = stdinPipe.Close()
			_ = cleanRemoteDir(client, lifecycleDir)
			return StackActionResult{Success: false, Message: fmt.Sprintf("falha ao empacotar Compose efêmero: %v", err)}
		}
		_ = stdinPipe.Close()

		if err := waitExtract(); err != nil {
			_ = cleanRemoteDir(client, lifecycleDir)
			return StackActionResult{Success: false, Message: fmt.Sprintf("falha ao extrair Compose efêmero na VPS: %v", err)}
		}

		targetDir = lifecycleDir
		isEphemeral = true
	}

	if isEphemeral {
		defer cleanRemoteDir(client, targetDir)
	}

	// 3. Executa docker compose down (com -v somente se deleteVolumes == true)
	cmdDown := BuildComposeDownCommand(targetDir, composeBin, stack.ProjectName, deleteVolumes)
	out, err := client.ExecCommand(cmdDown, false)
	if err != nil {
		return StackActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao remover containers da stack '%s': %v", stack.Name, err),
			Logs:    out,
		}
	}

	// 4. Cleanup pós-down: remove diretórios de deploy antigos da stack se nenhum container restante os utilizar
	deploysDir := fmt.Sprintf("$HOME/.docksea/deploys/%s", stack.ID)
	if canSafelyRemovePreviousDeploy(client, deploysDir) {
		_ = cleanRemoteDir(client, deploysDir)
	}

	// 5. Atualiza o banco de dados SQLite limpando `last_deployed_remote_dir`
	if database != nil {
		_ = database.UpdateStackDeployment(stack.ID, stack.LastDeployedYAML, stack.LastDeployedConfigYAML, stack.LastDeployedConfigJSON, "", time.Now().UTC())
	}

	return StackActionResult{
		Success: true,
		Message: fmt.Sprintf("Stack '%s' removida com sucesso da VPS!", stack.Name),
		Logs:    out,
	}
}

// ExecuteGetStackLogs obtém logs da stack via Docker Compose
func ExecuteGetStackLogs(
	client *connection.Client,
	stack *db.Stack,
	composePath string,
	tail int,
) (string, error) {
	if stack == nil {
		return "", fmt.Errorf("stack não informada")
	}
	if client == nil {
		return "", fmt.Errorf("conexão com VPS não estabelecida")
	}

	if tail <= 0 {
		tail = 200
	}

	composeBin := strings.TrimSpace(composePath)
	if composeBin == "" {
		composeBin = "docker compose"
	}

	targetDir := stack.LastDeployedRemoteDir
	if targetDir != "" {
		cmdLogs := fmt.Sprintf("cd %s && %s -p %s logs --tail %d", ShellQuote(targetDir), composeBin, ShellQuote(stack.ProjectName), tail)
		return client.ExecCommand(cmdLogs, false)
	}

	// Fallback sem cd caso não haja deployDir salvo
	cmdLogs := fmt.Sprintf("%s -p %s logs --tail %d", composeBin, ShellQuote(stack.ProjectName), tail)
	return client.ExecCommand(cmdLogs, false)
}
