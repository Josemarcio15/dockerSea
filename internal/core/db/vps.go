package db

import (
	"fmt"
	"time"
)

type VpsServer struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	ConnectionType    string    `json:"connectionType"` // 'local' | 'ssh'
	Host              string    `json:"host"`
	Port              int       `json:"port"`
	Username          string    `json:"username"`
	AuthType          string    `json:"authType"` // 'key' | 'password'
	SshKeyPath        string    `json:"sshKeyPath"`
	SshKeyPassphrase  string    `json:"sshKeyPassphrase"`
	SshPassword       string    `json:"sshPassword"`
	SudoPassword      string    `json:"sudoPassword"`
	DockerSocketPath  string    `json:"dockerSocketPath"`
	DockerPath        string    `json:"dockerPath"`
	DockerComposePath string    `json:"dockerComposePath"`
	IsActive          bool      `json:"isActive"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
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
		if s.SshKeyPassphrase, err = decryptSecret(s.SshKeyPassphrase); err != nil {
			return nil, err
		}
		if s.SshPassword, err = decryptSecret(s.SshPassword); err != nil {
			return nil, err
		}
		if s.SudoPassword, err = decryptSecret(s.SudoPassword); err != nil {
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
	var err error
	if s.SshKeyPassphrase, err = encryptSecret(s.SshKeyPassphrase); err != nil {
		return err
	}
	if s.SshPassword, err = encryptSecret(s.SshPassword); err != nil {
		return err
	}
	if s.SudoPassword, err = encryptSecret(s.SudoPassword); err != nil {
		return err
	}

	_, err = d.conn.Exec(`
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
