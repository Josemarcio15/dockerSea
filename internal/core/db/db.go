package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

type VpsServer struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	ConnectionType     string    `json:"connectionType"` // 'local' | 'ssh'
	Host               string    `json:"host"`
	Port               int       `json:"port"`
	Username           string    `json:"username"`
	AuthType           string    `json:"authType"` // 'key' | 'password'
	SshKeyPath         string    `json:"sshKeyPath"`
	SshKeyPassphrase   string    `json:"sshKeyPassphrase"`
	SshPassword        string    `json:"sshPassword"`
	SudoPassword       string    `json:"sudoPassword"`
	DockerSocketPath   string    `json:"dockerSocketPath"`
	DockerPath         string    `json:"dockerPath"`
	DockerComposePath  string    `json:"dockerComposePath"`
	IsActive           bool      `json:"isActive"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type Profile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Locale    string    `json:"locale"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func InitDB() (*DB, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}

	appDir := filepath.Join(configDir, "docksea")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config dir: %w", err)
	}

	dbPath := filepath.Join(appDir, "docksea.db")
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Enable WAL mode & foreign keys
	if _, err := conn.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		return nil, fmt.Errorf("failed to configure sqlite pragma: %w", err)
	}

	d := &DB{conn: conn}
	if err := d.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return d, nil
}

func (d *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS vps_servers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			connection_type TEXT NOT NULL,
			host TEXT,
			port INTEGER DEFAULT 22,
			username TEXT,
			auth_type TEXT,
			ssh_key_path TEXT,
			ssh_key_passphrase TEXT,
			ssh_password TEXT,
			sudo_password TEXT,
			docker_socket_path TEXT NOT NULL DEFAULT '/var/run/docker.sock',
			docker_path TEXT NOT NULL DEFAULT '/usr/bin/docker',
			docker_compose_path TEXT,
			is_active INTEGER DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS profiles (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			locale TEXT NOT NULL DEFAULT 'pt-BR',
			is_active INTEGER DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS app_settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);`,
	}

	for _, query := range queries {
		if _, err := d.conn.Exec(query); err != nil {
			return err
		}
	}



	// Inserir perfil padrão se a tabela estiver vazia
	var profileCount int
	err := d.conn.QueryRow(`SELECT COUNT(*) FROM profiles`).Scan(&profileCount)
	if err == nil && profileCount == 0 {
		now := time.Now().UTC()
		_, _ = d.conn.Exec(`
			INSERT INTO profiles (id, name, locale, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "default", "Perfil Padrão", "pt-BR", 1, now, now)
	}

	return nil
}

func (d *DB) ListVpsServers() ([]VpsServer, error) {
	rows, err := d.conn.Query(`
		SELECT 
			id, name, connection_type, COALESCE(host, ''), COALESCE(port, 22), 
			COALESCE(username, ''), COALESCE(auth_type, ''), COALESCE(ssh_key_path, ''), 
			COALESCE(ssh_key_passphrase, ''), COALESCE(ssh_password, ''), COALESCE(sudo_password, ''), 
			COALESCE(docker_socket_path, '/var/run/docker.sock'), COALESCE(docker_path, '/usr/bin/docker'), 
			COALESCE(docker_compose_path, ''), COALESCE(is_active, 0), created_at, updated_at
		FROM vps_servers
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []VpsServer
	for rows.Next() {
		var s VpsServer
		var isActiveInt int
		err := rows.Scan(
			&s.ID, &s.Name, &s.ConnectionType, &s.Host, &s.Port,
			&s.Username, &s.AuthType, &s.SshKeyPath,
			&s.SshKeyPassphrase, &s.SshPassword, &s.SudoPassword,
			&s.DockerSocketPath, &s.DockerPath,
			&s.DockerComposePath, &isActiveInt, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		s.IsActive = isActiveInt == 1
		servers = append(servers, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return servers, nil
}

func (d *DB) SaveVpsServer(s VpsServer) error {
	now := time.Now().UTC()
	if s.ID == "" {
		s.ID = fmt.Sprintf("vps-%d", time.Now().UnixNano())
	}
	if s.Port == 0 {
		s.Port = 22
	}
	if s.DockerSocketPath == "" {
		s.DockerSocketPath = "/var/run/docker.sock"
	}
	if s.DockerPath == "" {
		s.DockerPath = "/usr/bin/docker"
	}

	_, err := d.conn.Exec(`
		INSERT INTO vps_servers (
			id, name, connection_type, host, port, username, auth_type,
			ssh_key_path, ssh_key_passphrase, ssh_password, sudo_password,
			docker_socket_path, docker_path, docker_compose_path, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			connection_type = excluded.connection_type,
			host = excluded.host,
			port = excluded.port,
			username = excluded.username,
			auth_type = excluded.auth_type,
			ssh_key_path = excluded.ssh_key_path,
			ssh_key_passphrase = excluded.ssh_key_passphrase,
			ssh_password = excluded.ssh_password,
			sudo_password = excluded.sudo_password,
			docker_socket_path = excluded.docker_socket_path,
			docker_path = excluded.docker_path,
			docker_compose_path = excluded.docker_compose_path,
			updated_at = excluded.updated_at
	`, s.ID, s.Name, s.ConnectionType, s.Host, s.Port, s.Username, s.AuthType,
		s.SshKeyPath, s.SshKeyPassphrase, s.SshPassword, s.SudoPassword,
		s.DockerSocketPath, s.DockerPath, s.DockerComposePath, now, now)

	return err
}

func (d *DB) DeleteVpsServer(id string) error {
	_, err := d.conn.Exec(`DELETE FROM vps_servers WHERE id = ?`, id)
	return err
}

func (d *DB) SetActiveVpsServer(id string) error {
	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`UPDATE vps_servers SET is_active = 0`); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE vps_servers SET is_active = 1 WHERE id = ?`, id); err != nil {
		return err
	}

	return tx.Commit()
}

func (d *DB) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
