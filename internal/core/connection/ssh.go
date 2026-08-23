package connection

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go-walis/internal/core/db"

	"golang.org/x/crypto/ssh"
)

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
