package db

import (
	"fmt"
	"time"
)

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
