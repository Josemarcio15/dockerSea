package stacks

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
	sharedDocker "go-walis/internal/shared/docker"
)

var validIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// GenerateDeployID gera um identificador seguro e único para uma tentativa de deploy
func GenerateDeployID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("dep_%d_%s", time.Now().Unix(), hex.EncodeToString(b))
}

// BuildRemoteDeployDir monta o caminho remoto seguro para o diretório temporário/efêmero do deploy
func BuildRemoteDeployDir(stackID, deployID string) (string, error) {
	sID := strings.TrimSpace(stackID)
	dID := strings.TrimSpace(deployID)

	if sID == "" || !validIDRegex.MatchString(sID) {
		return "", fmt.Errorf("stack ID inválido para diretório remoto: %s", stackID)
	}
	if dID == "" || !validIDRegex.MatchString(dID) {
		return "", fmt.Errorf("deploy ID inválido para diretório remoto: %s", deployID)
	}

	return fmt.Sprintf("$HOME/.docksea/deploys/%s/%s", sID, dID), nil
}

// DeployOptions contém parâmetros opcionais ou customizados para o pipeline de deploy
type DeployOptions struct {
	Emitter EventEmitter
}

// ExecuteDeploy executa o pipeline de deploy remoto rigoroso e seguro
func ExecuteDeploy(
	ctx context.Context,
	database *db.DB,
	server db.VpsServer,
	stack *db.Stack,
	lockMgr *StackLockManager,
	opts ...DeployOptions,
) StackActionResult {
	if database == nil {
		return StackActionResult{Success: false, Message: "Banco de dados SQLite não inicializado"}
	}
	if stack == nil {
		return StackActionResult{Success: false, Message: "Stack informada é nula"}
	}

	// 1. Validar nome do projeto Docker Compose
	if err := ValidateProjectName(stack.ProjectName); err != nil {
		return StackActionResult{Success: false, Message: fmt.Sprintf("nome do projeto inválido: %v", err)}
	}

	// 2. Adquirir lock de concorrência por profileID:stackID
	lockKey := stack.ProfileID + ":" + stack.ID
	if lockMgr != nil {
		unlock := lockMgr.Lock(lockKey)
		defer unlock()
	}

	// 3. Inicializar deployID e broadcaster de eventos
	deployID := GenerateDeployID()
	var emitter EventEmitter
	if len(opts) > 0 && opts[0].Emitter != nil {
		emitter = opts[0].Emitter
	}
	broadcaster := NewDeployEventBroadcaster(stack.ID, deployID, emitter)
	broadcaster.EmitStarted(stack.ProjectName)

	var logBuffer bytes.Buffer
	recordLog := func(line string) {
		logBuffer.WriteString(line + "\n")
	}

	// 4. Obter ou criar conexão com o servidor VPS
	broadcaster.EmitProgress(PhasePreparing, "Conectando à VPS remota via SSH...")
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		msg := fmt.Sprintf("falha na conexão com o servidor: %v", err)
		broadcaster.EmitFailed(PhasePreparing, msg)
		return StackActionResult{Success: false, Message: msg}
	}
	defer client.Close()

	// 5. Detectar/validar capabilities do Docker e Docker Compose na VPS
	broadcaster.EmitProgress(PhasePreparing, "Validando capabilities do Docker e Compose no servidor...")
	caps, err := DetectServerCapabilities(ctx, client, server.DockerComposePath)
	if err != nil || !caps.DockerAvailable || !caps.ComposeAvailable {
		msg := "Docker ou Docker Compose não estão disponíveis no servidor remoto"
		if err != nil {
			msg = fmt.Sprintf("%s: %v", msg, err)
		} else if !caps.DockerAvailable {
			msg = "Docker Engine não foi detectado no servidor remoto"
		} else if !caps.ComposeAvailable {
			msg = "Docker Compose não foi detectado no servidor remoto"
		}
		broadcaster.EmitFailed(PhasePreparing, msg)
		return StackActionResult{Success: false, Message: msg}
	}

	// 6. Preparar diretório remoto isolado
	remoteDeployDir, err := BuildRemoteDeployDir(stack.ID, deployID)
	if err != nil {
		msg := fmt.Sprintf("falha ao construir diretório de deploy: %v", err)
		broadcaster.EmitFailed(PhasePreparing, msg)
		return StackActionResult{Success: false, Message: msg}
	}

	composeBin := strings.TrimSpace(server.DockerComposePath)
	if composeBin == "" {
		composeBin = "docker compose"
	}

	mkdirCmd := fmt.Sprintf("mkdir -p %s", remoteDeployDir)
	if _, err := client.ExecCommand(mkdirCmd, false); err != nil {
		msg := fmt.Sprintf("falha ao criar pasta remota de deploy: %v", err)
		broadcaster.EmitFailed(PhasePreparing, msg)
		return StackActionResult{Success: false, Message: msg}
	}

	// 7. Empacotamento e Upload via SSH Stream (tar -xzf - -C <remote_dir>)
	broadcaster.EmitProgress(PhaseUploading, fmt.Sprintf("Empacotando e enviando projeto para a VPS (%s)...", remoteDeployDir))
	extractCmd := fmt.Sprintf("tar -xzf - -C %s", remoteDeployDir)
	stdinPipe, waitExtract, err := client.StartCommandInput(ctx, extractCmd, false)
	if err != nil {
		msg := fmt.Sprintf("falha ao abrir stream de extração remota: %v", err)
		broadcaster.EmitFailed(PhaseUploading, msg)
		_ = cleanRemoteDir(client, remoteDeployDir)
		return StackActionResult{Success: false, Message: msg}
	}

	if stack.SourceType == "folder" && stack.FolderPath != "" {
		if err := PackProjectDir(ctx, stack.FolderPath, stdinPipe); err != nil {
			_ = stdinPipe.Close()
			msg := fmt.Sprintf("falha no empacotamento da pasta local: %v", err)
			broadcaster.EmitFailed(PhaseUploading, msg)
			_ = cleanRemoteDir(client, remoteDeployDir)
			return StackActionResult{Success: false, Message: msg}
		}
	} else {
		if err := PackSingleYaml(stack.YamlContent, stdinPipe); err != nil {
			_ = stdinPipe.Close()
			msg := fmt.Sprintf("falha no empacotamento do arquivo YAML: %v", err)
			broadcaster.EmitFailed(PhaseUploading, msg)
			_ = cleanRemoteDir(client, remoteDeployDir)
			return StackActionResult{Success: false, Message: msg}
		}
	}
	_ = stdinPipe.Close()

	if err := waitExtract(); err != nil {
		msg := fmt.Sprintf("falha ao descompactar arquivos na VPS: %v", err)
		broadcaster.EmitFailed(PhaseUploading, msg)
		_ = cleanRemoteDir(client, remoteDeployDir)
		return StackActionResult{Success: false, Message: msg}
	}
	broadcaster.EmitProgress(PhaseUploading, "Arquivos posicionados com sucesso na VPS.")

	// 8. Pré-Validação com `docker compose config --format json`
	broadcaster.EmitProgress(PhaseValidating, "Validando sintaxe e estrutura do Compose na VPS...")
	configJsonCmd := BuildComposeConfigCommand(remoteDeployDir, composeBin, stack.ProjectName, true)
	jsonOutput, err := client.ExecCommand(configJsonCmd, false)
	if err != nil {
		msg := fmt.Sprintf("validação do Docker Compose falhou na VPS: %v\nSaída: %s", err, jsonOutput)
		broadcaster.EmitFailed(PhaseValidating, msg)
		recordLog(msg)
		_ = cleanRemoteDir(client, remoteDeployDir)
		return StackActionResult{Success: false, Message: msg, Logs: logBuffer.String()}
	}

	parsedConfig, err := ParseComposeConfigOutput(jsonOutput)
	if err != nil {
		msg := fmt.Sprintf("falha ao interpretar retorno JSON do compose config: %v", err)
		broadcaster.EmitFailed(PhaseValidating, msg)
		recordLog(msg)
		_ = cleanRemoteDir(client, remoteDeployDir)
		return StackActionResult{Success: false, Message: msg, Logs: logBuffer.String()}
	}

	// 9. Capturar Snapshot Canônico YAML
	configYamlCmd := BuildComposeConfigCommand(remoteDeployDir, composeBin, stack.ProjectName, false)
	yamlOutput, err := client.ExecCommand(configYamlCmd, false)
	if err != nil {
		msg := fmt.Sprintf("falha ao extrair snapshot canônico YAML: %v", err)
		broadcaster.EmitFailed(PhaseValidating, msg)
		recordLog(msg)
		_ = cleanRemoteDir(client, remoteDeployDir)
		return StackActionResult{Success: false, Message: msg, Logs: logBuffer.String()}
	}

	candidateConfigJSON := jsonOutput
	candidateConfigYAML := yamlOutput

	// 10. Executar Deploy com `docker compose up -d --build --remove-orphans --wait`
	broadcaster.EmitProgress(PhaseBuilding, fmt.Sprintf("Iniciando build e deploy da stack '%s'...", stack.Name))
	upCmd := BuildComposeUpCommand(remoteDeployDir, composeBin, stack.ProjectName)
	stdoutPipe, waitUp, err := client.StartCommandOutput(ctx, upCmd, false)
	if err != nil {
		msg := fmt.Sprintf("falha ao iniciar execução do docker compose up: %v", err)
		broadcaster.EmitFailed(PhaseBuilding, msg)
		recordLog(msg)
		_ = cleanRemoteDir(client, remoteDeployDir)
		return StackActionResult{Success: false, Message: msg, Logs: logBuffer.String()}
	}

	scanner := bufio.NewScanner(stdoutPipe)
	for scanner.Scan() {
		line := scanner.Text()
		recordLog(line)
		broadcaster.EmitProgress(PhaseBuilding, line)
	}

	if err := waitUp(); err != nil {
		msg := fmt.Sprintf("execução do docker compose up falhou: %v", err)
		broadcaster.EmitFailed(PhaseStarting, msg)
		recordLog(msg)
		_ = cleanRemoteDir(client, remoteDeployDir)
		return StackActionResult{Success: false, Message: msg, Logs: logBuffer.String()}
	}

	broadcaster.EmitProgress(PhaseChecking, "Containers iniciados e verificados com sucesso!")

	// 11. Análise de Runtime Dependencies (Bind Mounts)
	broadcaster.EmitProgress(PhaseCleaning, "Analisando dependências de runtime e bind mounts...")
	hasBind := HasRuntimeBindDependency(parsedConfig, remoteDeployDir)

	var finalRemoteDir string
	if hasBind {
		finalRemoteDir = remoteDeployDir
		broadcaster.EmitProgress(PhaseCleaning, "Bind mounts detectados no diretório do projeto. Diretório preservado.")
	} else {
		finalRemoteDir = ""
		broadcaster.EmitProgress(PhaseCleaning, "Nenhum bind mount ativo no diretório. Limpando arquivos temporários...")
		_ = cleanRemoteDir(client, remoteDeployDir)
	}

	// 12. Rotação e Limpeza do Deploy Anterior (se houver e não for mais referenciado)
	if stack.LastDeployedRemoteDir != "" && stack.LastDeployedRemoteDir != finalRemoteDir {
		broadcaster.EmitProgress(PhaseCleaning, "Verificando se deploy anterior ainda é referenciado por outros containers...")
		canCleanPrevious := canSafelyRemovePreviousDeploy(client, stack.LastDeployedRemoteDir)
		if canCleanPrevious {
			_ = cleanRemoteDir(client, stack.LastDeployedRemoteDir)
		}
	}

	// 13. Persistência dos Metadados no SQLite (Somente após sucesso)
	deployedAt := time.Now().UTC()
	deployedYaml := stack.YamlContent
	if stack.SourceType == "folder" {
		deployedYaml = candidateConfigYAML
	}

	if err := database.UpdateStackDeployment(stack.ID, deployedYaml, candidateConfigYAML, candidateConfigJSON, finalRemoteDir, deployedAt); err != nil {
		msg := fmt.Sprintf("deploy concluído na VPS, mas falhou ao salvar metadados no SQLite: %v", err)
		broadcaster.EmitFailed(PhaseError, msg)
		return StackActionResult{Success: false, Message: msg, Logs: logBuffer.String()}
	}

	successMsg := fmt.Sprintf("Stack '%s' implantada e validada com sucesso!", stack.Name)
	broadcaster.EmitComplete(successMsg)

	return StackActionResult{
		Success: true,
		Message: successMsg,
		Logs:    logBuffer.String(),
	}
}

// DetectServerCapabilities detecta versões e capacidades do Docker e Docker Compose na VPS
func DetectServerCapabilities(ctx context.Context, client *connection.Client, composePath string) (ServerCapabilities, error) {
	var caps ServerCapabilities
	if client == nil {
		return caps, fmt.Errorf("cliente de conexão nulo")
	}

	// 1. Docker version
	dockerOut, err := client.ExecCommand("docker version --format '{{.Server.Version}}' 2>/dev/null || docker -v", false)
	if err == nil && strings.TrimSpace(dockerOut) != "" {
		caps.DockerAvailable = true
		caps.DockerVersion = strings.TrimSpace(dockerOut)
	}

	// 2. Compose binary/version
	composeBin := strings.TrimSpace(composePath)
	if composeBin == "" {
		composeBin = "docker compose"
	}
	composeOut, err := client.ExecCommand(fmt.Sprintf("%s version --short 2>/dev/null || %s version", composeBin, composeBin), false)
	if err == nil && strings.TrimSpace(composeOut) != "" {
		caps.ComposeAvailable = true
		caps.ComposeVersion = strings.TrimSpace(composeOut)
	}

	// 3. Buildx version
	buildxOut, err := client.ExecCommand("docker buildx version 2>/dev/null", false)
	if err == nil && strings.TrimSpace(buildxOut) != "" {
		caps.BuildxAvailable = true
	}

	return caps, nil
}

// cleanRemoteDir remove com segurança um diretório na VPS
func cleanRemoteDir(client *connection.Client, dir string) error {
	if client == nil || strings.TrimSpace(dir) == "" {
		return nil
	}
	// Validação de segurança para evitar rm -rf de diretórios raiz ou perigosos
	clean := filepath.Clean(dir)
	if clean == "/" || clean == "." || clean == ".." || clean == "$HOME" || clean == "~" {
		return fmt.Errorf("tentativa de remoção de caminho inseguro: %s", dir)
	}

	cmd := fmt.Sprintf("rm -rf %s", ShellQuote(dir))
	_, err := client.ExecCommand(cmd, false)
	return err
}

// canSafelyRemovePreviousDeploy verifica se algum container existente ainda referencia o diretório antigo
func canSafelyRemovePreviousDeploy(client *connection.Client, oldDir string) bool {
	if client == nil || strings.TrimSpace(oldDir) == "" {
		return false
	}
	// Inspeciona mounts de todos os containers no host remoto
	inspectCmd := "docker ps -a --format '{{.ID}}' | xargs -r docker inspect --format '{{range .Mounts}}{{println .Source}}{{end}}' 2>/dev/null"
	out, err := client.ExecCommand(inspectCmd, false)
	if err != nil {
		// Por precaução, se a checagem falhar, não apaga
		return false
	}

	cleanOld := filepath.Clean(oldDir)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		mountSource := filepath.Clean(strings.TrimSpace(line))
		if mountSource == "" {
			continue
		}
		if mountSource == cleanOld || strings.HasPrefix(mountSource, cleanOld+"/") || strings.HasPrefix(mountSource, cleanOld+"\\") {
			return false // Encontrou container usando esse diretório!
		}
	}

	return true
}
