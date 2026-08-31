package db

import (
	"fmt"
	"time"
)

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
