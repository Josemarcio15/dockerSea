package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

const temporaryEncryptionKey = "docksea-temporary-key-change-before-release"

func encryptionCipher() (cipher.AEAD, error) {
	salt := sha256.Sum256([]byte("docksea-v1-salt"))
	key := argon2.IDKey([]byte(temporaryEncryptionKey), salt[:], 1, 64*1024, 4, 32)
	block, err := aes.NewCipher(key)
	if err != nil { return nil, err }
	return cipher.NewGCM(block)
}

func encryptSecret(value string) (string, error) {
	if value == "" { return "", nil }
	aead, err := encryptionCipher()
	if err != nil { return "", err }
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil { return "", err }
	ciphertext := aead.Seal(nil, nonce, []byte(value), nil)
	encoded := append(nonce, ciphertext...)
	return "enc:v1:" + base64.RawStdEncoding.EncodeToString(encoded), nil
}

func decryptSecret(value string) (string, error) {
	if value == "" || !strings.HasPrefix(value, "enc:v1:") { return value, nil }
	aead, err := encryptionCipher()
	if err != nil { return "", err }
	data, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(value, "enc:v1:"))
	if err != nil || len(data) < aead.NonceSize() { return "", fmt.Errorf("segredo criptografado inválido") }
	plaintext, err := aead.Open(nil, data[:aead.NonceSize()], data[aead.NonceSize():], nil)
	if err != nil { return "", fmt.Errorf("não foi possível descriptografar segredo: %w", err) }
	return string(plaintext), nil
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

type ContainerConfig struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ProfileID     string    `json:"profileId"`
	Image         string    `json:"image"`
	ContainerName string    `json:"containerName"`
	Ports         *string   `json:"ports"`
	Env           *string   `json:"env"`
	Volumes       *string   `json:"volumes"`
	Network       *string   `json:"network"`
	RestartPolicy string    `json:"restartPolicy"`
	ProjectName   string    `json:"projectName"`
	Description   *string   `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Command       *string   `json:"command"`
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
		`CREATE TABLE IF NOT EXISTS image_history (
			id TEXT PRIMARY KEY,
			image_name TEXT NOT NULL,
			profile_id TEXT NOT NULL DEFAULT 'default',
			pulled_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS saved_paths (
			id TEXT PRIMARY KEY,
			path TEXT NOT NULL,
			label TEXT NOT NULL,
			profile_id TEXT NOT NULL DEFAULT 'default',
			created_at DATETIME NOT NULL,
			UNIQUE(profile_id, path)
		);`,
		`CREATE TABLE IF NOT EXISTS stacks (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			project_name TEXT NOT NULL,
			source_type TEXT NOT NULL DEFAULT 'editor',
			folder_path TEXT NOT NULL DEFAULT '',
			yaml_content TEXT NOT NULL DEFAULT '',
			profile_id TEXT NOT NULL DEFAULT 'default',
			last_deployed_yaml TEXT DEFAULT '',
			last_deployed_config_yaml TEXT DEFAULT '',
			last_deployed_config_json TEXT DEFAULT '',
			last_deployed_remote_dir TEXT DEFAULT '',
			last_deployed_at DATETIME,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS container_configs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			profile_id TEXT NOT NULL,
			image TEXT NOT NULL,
			container_name TEXT NOT NULL,
			ports TEXT,
			env TEXT,
			network TEXT,
			volumes TEXT,
			restart_policy TEXT,
			project_name TEXT,
			description TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			command TEXT
		);`,
	}

	for _, query := range queries {
		if _, err := d.conn.Exec(query); err != nil {
			return err
		}
	}

	// Migrações incrementais de colunas para bancos de dados existentes
	alterColumns := []string{
		`ALTER TABLE stacks ADD COLUMN last_deployed_yaml TEXT DEFAULT '';`,
		`ALTER TABLE stacks ADD COLUMN last_deployed_config_yaml TEXT DEFAULT '';`,
		`ALTER TABLE stacks ADD COLUMN last_deployed_config_json TEXT DEFAULT '';`,
		`ALTER TABLE stacks ADD COLUMN last_deployed_remote_dir TEXT DEFAULT '';`,
		`ALTER TABLE stacks ADD COLUMN last_deployed_at DATETIME;`,
	}
	for _, alterQuery := range alterColumns {
		_, _ = d.conn.Exec(alterQuery) // Ignora erro se coluna já existir no SQLite
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

func (d *DB) ListProfiles() ([]Profile, error) {
	rows, err := d.conn.Query(`
		SELECT id, name, locale, is_active, created_at, updated_at
		FROM profiles
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var p Profile
		var isActiveInt int
		if err := rows.Scan(&p.ID, &p.Name, &p.Locale, &isActiveInt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.IsActive = isActiveInt == 1
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (d *DB) GetActiveProfile() (*Profile, error) {
	var p Profile
	var isActiveInt int
	err := d.conn.QueryRow(`
		SELECT id, name, locale, is_active, created_at, updated_at
		FROM profiles
		WHERE is_active = 1
		LIMIT 1
	`).Scan(&p.ID, &p.Name, &p.Locale, &isActiveInt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		// Se nenhum estiver ativo, tenta pegar o primeiro
		err = d.conn.QueryRow(`
			SELECT id, name, locale, is_active, created_at, updated_at
			FROM profiles
			ORDER BY created_at ASC
			LIMIT 1
		`).Scan(&p.ID, &p.Name, &p.Locale, &isActiveInt, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}
	p.IsActive = true
	return &p, nil
}

func (d *DB) SaveProfile(p Profile) error {
	now := time.Now().UTC()
	if p.ID == "" {
		p.ID = fmt.Sprintf("prof_%d", now.UnixNano())
	}
	if p.Locale == "" {
		p.Locale = "pt-BR"
	}

	isActiveInt := 0
	if p.IsActive {
		isActiveInt = 1
	}

	_, err := d.conn.Exec(`
		INSERT INTO profiles (id, name, locale, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			locale = excluded.locale,
			updated_at = excluded.updated_at
	`, p.ID, p.Name, p.Locale, isActiveInt, now, now)
	return err
}

func (d *DB) DeleteProfile(id string) error {
	var count int
	err := d.conn.QueryRow(`SELECT COUNT(*) FROM profiles`).Scan(&count)
	if err != nil {
		return err
	}
	if count <= 1 {
		return fmt.Errorf("não é possível excluir o único perfil existente")
	}

	// Se estiver excluindo o perfil ativo, ativa outro
	var isActive int
	_ = d.conn.QueryRow(`SELECT is_active FROM profiles WHERE id = ?`, id).Scan(&isActive)

	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM profiles WHERE id = ?`, id); err != nil {
		return err
	}
	// Também limpa histórico de imagens e paths atrelados a esse perfil
	_, _ = tx.Exec(`DELETE FROM image_history WHERE profile_id = ?`, id)
	_, _ = tx.Exec(`DELETE FROM saved_paths WHERE profile_id = ?`, id)

	if isActive == 1 {
		if _, err := tx.Exec(`UPDATE profiles SET is_active = 1 WHERE id IN (SELECT id FROM profiles LIMIT 1)`); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *DB) SetActiveProfile(id string) error {
	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`UPDATE profiles SET is_active = 0`); err != nil {
		return err
	}
	res, err := tx.Exec(`UPDATE profiles SET is_active = 1 WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("perfil não encontrado")
	}

	return tx.Commit()
}

func (d *DB) UpdateProfileLocale(id string, locale string) error {
	now := time.Now().UTC()
	_, err := d.conn.Exec(`
		UPDATE profiles SET locale = ?, updated_at = ? WHERE id = ?
	`, locale, now, id)
	return err
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
		if s.SshKeyPassphrase, err = decryptSecret(s.SshKeyPassphrase); err != nil { return nil, err }
		if s.SshPassword, err = decryptSecret(s.SshPassword); err != nil { return nil, err }
		if s.SudoPassword, err = decryptSecret(s.SudoPassword); err != nil { return nil, err }
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
	if s.SshKeyPassphrase, err = encryptSecret(s.SshKeyPassphrase); err != nil { return err }
	if s.SshPassword, err = encryptSecret(s.SshPassword); err != nil { return err }
	if s.SudoPassword, err = encryptSecret(s.SudoPassword); err != nil { return err }

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

type ImageHistoryItem struct {
	ID        string    `json:"id"`
	ImageName string    `json:"imageName"`
	ProfileID string    `json:"profileId"`
	PulledAt  time.Time `json:"pulledAt"`
}

func (d *DB) ListImageHistory(profileID string) ([]ImageHistoryItem, error) {
	if profileID == "" {
		profileID = "default"
	}
	rows, err := d.conn.Query(`
		SELECT id, image_name, profile_id, pulled_at
		FROM image_history
		WHERE profile_id = ?
		ORDER BY pulled_at DESC
		LIMIT 100
	`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ImageHistoryItem
	for rows.Next() {
		var item ImageHistoryItem
		if err := rows.Scan(&item.ID, &item.ImageName, &item.ProfileID, &item.PulledAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *DB) AddImageHistory(imageName, profileID string) error {
	if profileID == "" {
		profileID = "default"
	}
	id := fmt.Sprintf("hist-%d", time.Now().UnixNano())
	now := time.Now().UTC()

	// Remover duplicata se o mesmo nome de imagem já existir para este perfil
	_, _ = d.conn.Exec(`DELETE FROM image_history WHERE image_name = ? AND profile_id = ?`, imageName, profileID)

	_, err := d.conn.Exec(`
		INSERT INTO image_history (id, image_name, profile_id, pulled_at)
		VALUES (?, ?, ?, ?)
	`, id, imageName, profileID, now)
	return err
}

func (d *DB) DeleteImageHistory(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`DELETE FROM image_history WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, id := range ids {
		if _, err := stmt.Exec(id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d *DB) ClearImageHistory(profileID string) error {
	if profileID == "" {
		profileID = "default"
	}
	_, err := d.conn.Exec(`DELETE FROM image_history WHERE profile_id = ?`, profileID)
	return err
}

func (d *DB) ListContainerConfigs(profileID string) ([]ContainerConfig, error) {
	if profileID == "" {
		p, _ := d.GetActiveProfile()
		if p != nil { profileID = p.ID } else { profileID = "default" }
	}
	rows, err := d.conn.Query(`SELECT id,name,profile_id,image,container_name,ports,env,volumes,network,restart_policy,project_name,description,created_at,updated_at,command FROM container_configs WHERE profile_id = ? ORDER BY name ASC`, profileID)
	if err != nil { return nil, err }
	defer rows.Close()
	var items []ContainerConfig
	for rows.Next() {
		var item ContainerConfig
		if err := rows.Scan(&item.ID, &item.Name, &item.ProfileID, &item.Image, &item.ContainerName, &item.Ports, &item.Env, &item.Volumes, &item.Network, &item.RestartPolicy, &item.ProjectName, &item.Description, &item.CreatedAt, &item.UpdatedAt, &item.Command); err != nil { return nil, err }
		items = append(items, item)
	}
	return items, rows.Err()
}

func (d *DB) SaveContainerConfig(item ContainerConfig) error {
	if item.Name == "" || item.Image == "" { return fmt.Errorf("nome e imagem são obrigatórios") }
	if item.ProfileID == "" { p, _ := d.GetActiveProfile(); if p != nil { item.ProfileID = p.ID } else { item.ProfileID = "default" } }
	if item.ID == "" { item.ID = fmt.Sprintf("cfg-%d", time.Now().UnixNano()) }
	now := time.Now().UTC()
	_, err := d.conn.Exec(`INSERT INTO container_configs (id,name,profile_id,image,container_name,ports,env,volumes,network,restart_policy,project_name,description,created_at,updated_at,command) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,profile_id=excluded.profile_id,image=excluded.image,container_name=excluded.container_name,ports=excluded.ports,env=excluded.env,volumes=excluded.volumes,network=excluded.network,restart_policy=excluded.restart_policy,project_name=excluded.project_name,description=excluded.description,updated_at=excluded.updated_at,command=excluded.command`, item.ID,item.Name,item.ProfileID,item.Image,item.ContainerName,item.Ports,item.Env,item.Volumes,item.Network,item.RestartPolicy,item.ProjectName,item.Description,now,now,item.Command)
	return err
}

func (d *DB) DeleteContainerConfig(id string) error {
	if id == "" { return fmt.Errorf("id inválido") }
	_, err := d.conn.Exec(`DELETE FROM container_configs WHERE id = ?`, id)
	return err
}

type SavedPathItem struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	Label     string    `json:"label"`
	ProfileID string    `json:"profileId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (d *DB) ListSavedPaths(profileID string) ([]SavedPathItem, error) {
	if profileID == "" {
		p, _ := d.GetActiveProfile()
		if p != nil {
			profileID = p.ID
		} else {
			profileID = "default"
		}
	}
	rows, err := d.conn.Query(`
		SELECT id, path, label, profile_id, created_at
		FROM saved_paths
		WHERE profile_id = ?
		ORDER BY label ASC
	`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []SavedPathItem
	for rows.Next() {
		var item SavedPathItem
		if err := rows.Scan(&item.ID, &item.Path, &item.Label, &item.ProfileID, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *DB) SaveSavedPath(path, label, profileID string) error {
	if profileID == "" {
		p, _ := d.GetActiveProfile()
		if p != nil {
			profileID = p.ID
		} else {
			profileID = "default"
		}
	}
	id := fmt.Sprintf("path-%d", time.Now().UnixNano())
	now := time.Now().UTC()

	_, err := d.conn.Exec(`
		INSERT INTO saved_paths (id, path, label, profile_id, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(profile_id, path) DO UPDATE SET
			label = excluded.label
	`, id, path, label, profileID, now)
	return err
}

func (d *DB) DeleteSavedPath(path, profileID string) error {
	if profileID == "" {
		p, _ := d.GetActiveProfile()
		if p != nil {
			profileID = p.ID
		} else {
			profileID = "default"
		}
	}
	_, err := d.conn.Exec(`DELETE FROM saved_paths WHERE path = ? AND profile_id = ?`, path, profileID)
	return err
}

type Stack struct {
	ID                     string     `json:"id"`
	Name                   string     `json:"name"`
	ProjectName            string     `json:"projectName"`
	SourceType             string     `json:"sourceType"` // 'editor' | 'folder'
	FolderPath             string     `json:"folderPath"`
	YamlContent            string     `json:"yamlContent"`
	ProfileID              string     `json:"profileId"`
	LastDeployedYAML       string     `json:"lastDeployedYaml,omitempty"`
	LastDeployedConfigYAML string     `json:"lastDeployedConfigYaml,omitempty"`
	LastDeployedConfigJSON string     `json:"lastDeployedConfigJson,omitempty"`
	LastDeployedRemoteDir  string     `json:"lastDeployedRemoteDir,omitempty"`
	LastDeployedAt         *time.Time `json:"lastDeployedAt,omitempty"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
}

func (d *DB) ListStacks(profileID string) ([]Stack, error) {
	if profileID == "" {
		p, _ := d.GetActiveProfile()
		if p != nil {
			profileID = p.ID
		} else {
			profileID = "default"
		}
	}
	rows, err := d.conn.Query(`
		SELECT id, name, project_name, source_type, folder_path, yaml_content, profile_id,
		       COALESCE(last_deployed_yaml, ''), COALESCE(last_deployed_config_yaml, ''),
		       COALESCE(last_deployed_config_json, ''), COALESCE(last_deployed_remote_dir, ''),
		       last_deployed_at, created_at, updated_at
		FROM stacks
		WHERE profile_id = ?
		ORDER BY updated_at DESC
	`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Stack
	for rows.Next() {
		var s Stack
		if err := rows.Scan(
			&s.ID, &s.Name, &s.ProjectName, &s.SourceType, &s.FolderPath, &s.YamlContent, &s.ProfileID,
			&s.LastDeployedYAML, &s.LastDeployedConfigYAML, &s.LastDeployedConfigJSON, &s.LastDeployedRemoteDir,
			&s.LastDeployedAt, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (d *DB) GetStack(id string) (*Stack, error) {
	var s Stack
	err := d.conn.QueryRow(`
		SELECT id, name, project_name, source_type, folder_path, yaml_content, profile_id,
		       COALESCE(last_deployed_yaml, ''), COALESCE(last_deployed_config_yaml, ''),
		       COALESCE(last_deployed_config_json, ''), COALESCE(last_deployed_remote_dir, ''),
		       last_deployed_at, created_at, updated_at
		FROM stacks
		WHERE id = ?
		LIMIT 1
	`, id).Scan(
		&s.ID, &s.Name, &s.ProjectName, &s.SourceType, &s.FolderPath, &s.YamlContent, &s.ProfileID,
		&s.LastDeployedYAML, &s.LastDeployedConfigYAML, &s.LastDeployedConfigJSON, &s.LastDeployedRemoteDir,
		&s.LastDeployedAt, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (d *DB) SaveStack(s Stack) error {
	if s.ProfileID == "" {
		p, _ := d.GetActiveProfile()
		if p != nil {
			s.ProfileID = p.ID
		} else {
			s.ProfileID = "default"
		}
	}
	now := time.Now().UTC()
	if s.ID == "" {
		s.ID = fmt.Sprintf("stk_%d", now.UnixNano())
		s.CreatedAt = now
	}
	s.UpdatedAt = now
	if s.SourceType == "" {
		s.SourceType = "editor"
	}

	_, err := d.conn.Exec(`
		INSERT INTO stacks (
			id, name, project_name, source_type, folder_path, yaml_content, profile_id,
			last_deployed_yaml, last_deployed_config_yaml, last_deployed_config_json,
			last_deployed_remote_dir, last_deployed_at, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			project_name = excluded.project_name,
			source_type = excluded.source_type,
			folder_path = excluded.folder_path,
			yaml_content = excluded.yaml_content,
			last_deployed_yaml = excluded.last_deployed_yaml,
			last_deployed_config_yaml = excluded.last_deployed_config_yaml,
			last_deployed_config_json = excluded.last_deployed_config_json,
			last_deployed_remote_dir = excluded.last_deployed_remote_dir,
			last_deployed_at = excluded.last_deployed_at,
			updated_at = excluded.updated_at
	`, s.ID, s.Name, s.ProjectName, s.SourceType, s.FolderPath, s.YamlContent, s.ProfileID,
		s.LastDeployedYAML, s.LastDeployedConfigYAML, s.LastDeployedConfigJSON,
		s.LastDeployedRemoteDir, s.LastDeployedAt, s.CreatedAt, s.UpdatedAt)
	return err
}

func (d *DB) UpdateStackDeployment(id, deployedYaml, configYaml, configJson, remoteDir string, deployedAt time.Time) error {
	now := time.Now().UTC()
	_, err := d.conn.Exec(`
		UPDATE stacks SET
			last_deployed_yaml = ?,
			last_deployed_config_yaml = ?,
			last_deployed_config_json = ?,
			last_deployed_remote_dir = ?,
			last_deployed_at = ?,
			updated_at = ?
		WHERE id = ?
	`, deployedYaml, configYaml, configJson, remoteDir, deployedAt, now, id)
	return err
}

func (d *DB) DeleteStack(id string) error {
	_, err := d.conn.Exec(`DELETE FROM stacks WHERE id = ?`, id)
	return err
}

func (d *DB) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
