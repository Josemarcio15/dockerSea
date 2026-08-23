package connection

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"go-walis/internal/core/db"

	"golang.org/x/crypto/ssh"
)

// Client representa uma conexão ativa (SSH remota ou Local)
type Client struct {
	server     db.VpsServer
	sshClient  *ssh.Client
	httpClient *http.Client
	mu         sync.Mutex
}

// Manager gerencia instâncias de conexões ativas reutilizáveis
type Manager struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

var globalManager = &Manager{
	clients: make(map[string]*Client),
}

// GetManager retorna a instância global do gerenciador de conexões
func GetManager() *Manager {
	return globalManager
}

// GetClient obtém ou cria uma nova conexão para o servidor
func (m *Manager) GetClient(server db.VpsServer) (*Client, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client, exists := m.clients[server.ID]; exists {
		if client.IsAlive() {
			return client, nil
		}
		client.Close()
		delete(m.clients, server.ID)
	}

	client, err := NewClient(server)
	if err != nil {
		return nil, err
	}

	m.clients[server.ID] = client
	return client, nil
}

// NewClient instancia um novo cliente de conexão
func NewClient(server db.VpsServer) (*Client, error) {
	c := &Client{server: server}

	if isLocal(server) {
		socketPath := server.DockerSocketPath
		if socketPath == "" {
			socketPath = "/var/run/docker.sock"
		}
		c.httpClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, proto, addr string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
			Timeout: 30 * time.Second,
		}
		return c, nil
	}

	// Conexão SSH
	sshClient, err := createSshClient(server)
	if err != nil {
		return nil, err
	}
	c.sshClient = sshClient

	// HTTP Client tunelado pelo SSH para o socket Docker remoto
	c.httpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, proto, addr string) (net.Conn, error) {
				socketPath := server.DockerSocketPath
				if socketPath == "" {
					socketPath = "/var/run/docker.sock"
				}
				return sshClient.Dial("unix", socketPath)
			},
		},
		Timeout: 30 * time.Second,
	}

	return c, nil
}

// GetHttpClient retorna o cliente HTTP com timeout padrão para requisições rápidas
func (c *Client) GetHttpClient() *http.Client {
	return c.httpClient
}

// GetStreamHttpClient retorna um cliente HTTP sem timeout para conexões longas/streaming
func (c *Client) GetStreamHttpClient() *http.Client {
	if c.httpClient == nil {
		return nil
	}
	return &http.Client{
		Transport: c.httpClient.Transport,
		Timeout:   0, // Sem timeout para streaming persistente
	}
}

// ExecCommand executa um comando no servidor (com ou sem elevação sudo inteligente)
func (c *Client) ExecCommand(cmd string, useSudo bool) (string, error) {
	finalCmd := cmd

	if useSudo {
		username := strings.TrimSpace(strings.ToLower(c.server.Username))
		if username == "root" {
			// Já é superusuário root, executa direto sem sudo
			finalCmd = cmd
		} else if strings.TrimSpace(c.server.SudoPassword) != "" {
			// Sudo com senha informada via stdin
			finalCmd = fmt.Sprintf("echo %s | sudo -S -p '' %s", escapeShell(c.server.SudoPassword), cmd)
		} else {
			// Sudo sem senha (NOPASSWD). Usar flag -n (non-interactive) para não travar se o sudoers exigir senha
			if !strings.HasPrefix(cmd, "sudo ") {
				finalCmd = fmt.Sprintf("sudo -n %s", cmd)
			}
		}
	}

	if isLocal(c.server) {
		// Execução local
		return execLocalCommand(finalCmd)
	}

	// Execução remota via SSH
	if c.sshClient == nil {
		return "", fmt.Errorf("cliente SSH não conectado")
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("falha ao criar sessão SSH: %w", err)
	}
	defer session.Close()

	out, err := session.CombinedOutput(finalCmd)
	if err != nil {
		return string(out), fmt.Errorf("comando falhou (%w): %s", err, string(out))
	}

	return string(out), nil
}

// IsAlive verifica se a conexão ainda está saudável
func (c *Client) IsAlive() bool {
	if isLocal(c.server) {
		socketPath := c.server.DockerSocketPath
		if socketPath == "" {
			socketPath = "/var/run/docker.sock"
		}
		_, err := os.Stat(socketPath)
		return err == nil
	}

	if c.sshClient == nil {
		return false
	}
	_, _, err := c.sshClient.SendRequest("keepalive@openssh.com", true, nil)
	return err == nil
}

// Close encerra as conexões ativas
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sshClient != nil {
		_ = c.sshClient.Close()
		c.sshClient = nil
	}
}

// createSshClient constrói a conexão SSH usando a regra inteligente de autenticação
func createSshClient(server db.VpsServer) (*ssh.Client, error) {
	if strings.TrimSpace(server.Host) == "" {
		return nil, fmt.Errorf("o host ou IP da VPS é obrigatório")
	}

	port := server.Port
	if port == 0 {
		port = 22
	}
	targetAddr := fmt.Sprintf("%s:%d", strings.TrimSpace(server.Host), port)

	var authMethods []ssh.AuthMethod

	// 1. Chave privada (se informada)
	keyPath := strings.TrimSpace(server.SshKeyPath)
	if keyPath != "" {
		if strings.HasPrefix(keyPath, "~/") {
			home, _ := os.UserHomeDir()
			keyPath = home + keyPath[1:]
		}

		keyBytes, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("não foi possível ler o arquivo da chave SSH '%s': %w", keyPath, err)
		}

		var signer ssh.Signer
		passphrase := strings.TrimSpace(server.SshKeyPassphrase)
		if passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(keyBytes)
		}

		if err != nil {
			return nil, fmt.Errorf("chave privada SSH inválida ou senha incorreta: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// 2. Senha SSH (se informada)
	password := strings.TrimSpace(server.SshPassword)
	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	// Se nada foi informado
	if len(authMethods) == 0 {
		return nil, fmt.Errorf("nenhum método de autenticação SSH fornecido (preencha a chave privada ou senha)")
	}

	username := strings.TrimSpace(server.Username)
	if username == "" {
		username = "root"
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         8 * time.Second,
	}

	return ssh.Dial("tcp", targetAddr, config)
}

func isLocal(server db.VpsServer) bool {
	if server.ConnectionType == "local" {
		return true
	}
	host := strings.TrimSpace(strings.ToLower(server.Host))
	return host == "" || host == "localhost" || host == "127.0.0.1"
}

func escapeShell(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
