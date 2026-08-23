package connection

import (
	"context"
	"net"
	"net/http"
	"os"
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

// GetStreamHttpClient retorna um cliente HTTP com transporte dedicado e sem timeout para streaming longo
func (c *Client) GetStreamHttpClient() *http.Client {
	if isLocal(c.server) {
		socketPath := c.server.DockerSocketPath
		if socketPath == "" {
			socketPath = "/var/run/docker.sock"
		}
		return &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, proto, addr string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
				DisableKeepAlives: true,
			},
			Timeout: 0,
		}
	}

	if c.sshClient == nil {
		return nil
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, proto, addr string) (net.Conn, error) {
				socketPath := c.server.DockerSocketPath
				if socketPath == "" {
					socketPath = "/var/run/docker.sock"
				}
				return c.sshClient.Dial("unix", socketPath)
			},
			DisableKeepAlives: true,
		},
		Timeout: 0,
	}
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
