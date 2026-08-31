package db

import (
	"fmt"
	"time"
)

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

func (d *DB) ListContainerConfigs(profileID string) ([]ContainerConfig, error) {
	if profileID == "" {
		p, _ := d.GetActiveProfile()
		if p != nil {
			profileID = p.ID
		} else {
			profileID = "default"
		}
	}
	rows, err := d.conn.Query(`SELECT id,name,profile_id,image,container_name,ports,env,volumes,network,restart_policy,project_name,description,created_at,updated_at,command FROM container_configs WHERE profile_id = ? ORDER BY name ASC`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ContainerConfig
	for rows.Next() {
		var item ContainerConfig
		if err := rows.Scan(&item.ID, &item.Name, &item.ProfileID, &item.Image, &item.ContainerName, &item.Ports, &item.Env, &item.Volumes, &item.Network, &item.RestartPolicy, &item.ProjectName, &item.Description, &item.CreatedAt, &item.UpdatedAt, &item.Command); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (d *DB) SaveContainerConfig(item ContainerConfig) error {
	if item.Name == "" || item.Image == "" {
		return fmt.Errorf("nome e imagem são obrigatórios")
	}
	if item.ProfileID == "" {
		p, _ := d.GetActiveProfile()
		if p != nil {
			item.ProfileID = p.ID
		} else {
			item.ProfileID = "default"
		}
	}
	if item.ID == "" {
		item.ID = fmt.Sprintf("cfg-%d", time.Now().UnixNano())
	}
	now := time.Now().UTC()
	_, err := d.conn.Exec(`INSERT INTO container_configs (id,name,profile_id,image,container_name,ports,env,volumes,network,restart_policy,project_name,description,created_at,updated_at,command) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,profile_id=excluded.profile_id,image=excluded.image,container_name=excluded.container_name,ports=excluded.ports,env=excluded.env,volumes=excluded.volumes,network=excluded.network,restart_policy=excluded.restart_policy,project_name=excluded.project_name,description=excluded.description,updated_at=excluded.updated_at,command=excluded.command`, item.ID, item.Name, item.ProfileID, item.Image, item.ContainerName, item.Ports, item.Env, item.Volumes, item.Network, item.RestartPolicy, item.ProjectName, item.Description, now, now, item.Command)
	return err
}

func (d *DB) DeleteContainerConfig(id string) error {
	if id == "" {
		return fmt.Errorf("id inválido")
	}
	_, err := d.conn.Exec(`DELETE FROM container_configs WHERE id = ?`, id)
	return err
}
