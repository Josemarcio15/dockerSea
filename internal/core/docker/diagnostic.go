package docker

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"go-walis/internal/core/db"
	sharedDocker "go-walis/internal/shared/docker"
)

type DiagnosticStep struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // "success" | "error" | "warning"
	Message string `json:"message"`
}

type DiagnosticResult struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Steps   []DiagnosticStep `json:"steps"`
}

// RunDiagnostic testa e valida a conexão usando a engine central de connection
func RunDiagnostic(server db.VpsServer) DiagnosticResult {
	if isLocal(server) {
		return testLocalDiagnostic(server)
	}
	return testSshDiagnostic(server)
}

func testLocalDiagnostic(server db.VpsServer) DiagnosticResult {
	var steps []DiagnosticStep
	socketPath := server.DockerSocketPath
	if socketPath == "" {
		socketPath = "/var/run/docker.sock"
	}

	// 1. Tenta inicializar o cliente de conexão local
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		steps = append(steps, DiagnosticStep{
			Name:    "Conexão Docker Local",
			Status:  "error",
			Message: fmt.Sprintf("Não foi possível acessar o socket local '%s': %v", socketPath, err),
		})
		return DiagnosticResult{
			Success: false,
			Message: "Falha ao conectar no Docker local.",
			Steps:   steps,
		}
	}
	defer client.Close()

	steps = append(steps, DiagnosticStep{
		Name:    "Conexão Docker Local",
		Status:  "success",
		Message: fmt.Sprintf("Socket UNIX '%s' acessível.", socketPath),
	})

	// 2. Ping real na API do Docker via UNIX socket
	resp, err := client.GetHttpClient().Get("http://localhost/_ping")
	if err != nil {
		steps = append(steps, DiagnosticStep{
			Name:    "Docker Daemon Ping",
			Status:  "error",
			Message: fmt.Sprintf("Socket encontrado, mas a API do Docker não respondeu: %v", err),
		})
		return DiagnosticResult{
			Success: false,
			Message: "Daemon do Docker local inativo ou sem permissão de acesso.",
			Steps:   steps,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		steps = append(steps, DiagnosticStep{
			Name:    "Docker Daemon Ping",
			Status:  "error",
			Message: fmt.Sprintf("Docker retornou HTTP status inesperado: %d", resp.StatusCode),
		})
		return DiagnosticResult{
			Success: false,
			Message: "Docker Daemon respondeu com erro.",
			Steps:   steps,
		}
	}

	steps = append(steps, DiagnosticStep{
		Name:    "Docker Daemon Ping",
		Status:  "success",
		Message: "API do Docker daemon respondeu com sucesso (HTTP 200 OK).",
	})

	return DiagnosticResult{
		Success: true,
		Message: "Docker local conectado e operacional com sucesso!",
		Steps:   steps,
	}
}

func testSshDiagnostic(server db.VpsServer) DiagnosticResult {
	var steps []DiagnosticStep
	host := strings.TrimSpace(server.Host)
	if host == "" {
		steps = append(steps, DiagnosticStep{
			Name:    "Validação de Parâmetros",
			Status:  "error",
			Message: "O campo Host/IP está em branco. Informe o endereço IP ou domínio da VPS.",
		})
		return DiagnosticResult{
			Success: false,
			Message: "Configuração incompleta: Host/IP é obrigatório.",
			Steps:   steps,
		}
	}

	port := server.Port
	if port == 0 {
		port = 22
	}
	targetAddr := fmt.Sprintf("%s:%d", host, port)

	// 1. Resolução DNS e Teste de Conectividade de Rede (Porta TCP)
	connStart := time.Now()
	tcpConn, err := net.DialTimeout("tcp", targetAddr, 5*time.Second)
	if err != nil {
		steps = append(steps, DiagnosticStep{
			Name:    "Conectividade de Rede (TCP)",
			Status:  "error",
			Message: fmt.Sprintf("Não foi possível alcançar o servidor em %s (porta %d): %v", host, port, err),
		})
		return DiagnosticResult{
			Success: false,
			Message: fmt.Sprintf("Servidor inacessível em %s:%d. Verifique seu firewall ou porta SSH.", host, port),
			Steps:   steps,
		}
	}
	_ = tcpConn.Close()
	tcpDuration := time.Since(connStart).Round(time.Millisecond)

	steps = append(steps, DiagnosticStep{
		Name:    "Conectividade de Rede (TCP)",
		Status:  "success",
		Message: fmt.Sprintf("Porta %d aberta e respondendo em %s (latência: %s).", port, host, tcpDuration),
	})

	// 2. Autenticação e Estabelecimento de Conexão com connection.NewClient
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		steps = append(steps, DiagnosticStep{
			Name:    "Autenticação SSH",
			Status:  "error",
			Message: fmt.Sprintf("Falha ao autenticar no servidor: %v", err),
		})
		return DiagnosticResult{
			Success: false,
			Message: fmt.Sprintf("Falha de autenticação SSH: %v", err),
			Steps:   steps,
		}
	}
	defer client.Close()

	user := strings.TrimSpace(server.Username)
	if user == "" {
		user = "root"
	}
	steps = append(steps, DiagnosticStep{
		Name:    "Autenticação SSH",
		Status:  "success",
		Message: fmt.Sprintf("Autenticado com sucesso como usuário '%s'.", user),
	})

	// 3. Identificação do Sistema Operacional Remoto
	unameOut, err := client.ExecCommand("uname -srm", false)
	if err != nil {
		steps = append(steps, DiagnosticStep{
			Name:    "Execução Remota",
			Status:  "warning",
			Message: fmt.Sprintf("Não foi possível consultar informações do OS: %v", err),
		})
	} else {
		steps = append(steps, DiagnosticStep{
			Name:    "Sistema Operacional",
			Status:  "success",
			Message: fmt.Sprintf("Kernel: %s", strings.TrimSpace(unameOut)),
		})
	}

	// 4. Teste de Elevação de Privilégios (Sudo)
	sudoOut, err := client.ExecCommand("id -u", true)
	if err == nil && strings.TrimSpace(sudoOut) == "0" {
		steps = append(steps, DiagnosticStep{
			Name:    "Permissões Superusuário (Sudo)",
			Status:  "success",
			Message: "Privilégios de superusuário (root / NOPASSWD) confirmados e funcionais.",
		})
	} else {
		msg := "O usuário atual requer senha para executar 'sudo' mas nenhuma 'Senha do Sudo' foi preenchida na aba Docker & Sudo."
		if strings.TrimSpace(server.SudoPassword) != "" {
			msg = "A senha do Sudo informada na aba Docker & Sudo foi recusada pelo servidor."
		}
		steps = append(steps, DiagnosticStep{
			Name:    "Permissões Superusuário (Sudo)",
			Status:  "warning",
			Message: msg,
		})
	}

	// 5. Teste do Socket do Docker Remoto
	dockerSock := server.DockerSocketPath
	if dockerSock == "" {
		dockerSock = "/var/run/docker.sock"
	}

	sockCheck, err := client.ExecCommand(fmt.Sprintf("test -S %s && echo OK || echo FAIL", dockerSock), false)
	if err == nil && strings.Contains(sockCheck, "OK") {
		steps = append(steps, DiagnosticStep{
			Name:    "Docker Daemon Remoto",
			Status:  "success",
			Message: fmt.Sprintf("Socket do Docker encontrado em '%s'.", dockerSock),
		})
	} else {
		steps = append(steps, DiagnosticStep{
			Name:    "Docker Daemon Remoto",
			Status:  "warning",
			Message: fmt.Sprintf("Socket '%s' não encontrado. Verifique se o Docker está instalado e em execução na VPS.", dockerSock),
		})
	}

	// Sucesso geral confirmado
	return DiagnosticResult{
		Success: true,
		Message: fmt.Sprintf("Conexão SSH com '%s' (%s) estabelecida e validada com sucesso!", server.Name, host),
		Steps:   steps,
	}
}

func isLocal(server db.VpsServer) bool {
	if server.ConnectionType == "local" {
		return true
	}
	host := strings.TrimSpace(strings.ToLower(server.Host))
	return host == "" || host == "localhost" || host == "127.0.0.1"
}
