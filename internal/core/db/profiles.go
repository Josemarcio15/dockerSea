package db

import (
	"fmt"
	"time"
)

type Profile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Locale    string    `json:"locale"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
	if err := rows.Err(); err != nil {
		return nil, err
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
