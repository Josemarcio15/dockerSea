package db

import (
	"fmt"
	"time"
)

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
