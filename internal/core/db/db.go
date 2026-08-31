package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn   *sql.DB
	dbPath string
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

	d := &DB{conn: conn, dbPath: dbPath}
	if err := d.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return d, nil
}

func (d *DB) GetDBPath() string {
	return d.dbPath
}

// Backup realiza um checkpoint completo do WAL e salva uma cópia limpa do banco no caminho de destino
func (d *DB) Backup(destPath string) error {
	if d.conn == nil {
		return fmt.Errorf("conexão com o banco não inicializada")
	}

	// Forçar flush do WAL antes do backup
	if _, err := d.conn.Exec(`PRAGMA wal_checkpoint(TRUNCATE);`); err != nil {
		return fmt.Errorf("falha ao sincronizar WAL: %w", err)
	}

	// Usa o comando nativo de VACUUM INTO do SQLite para gerar um arquivo consistente
	_ = os.Remove(destPath)
	_, err := d.conn.Exec(`VACUUM INTO ?;`, destPath)
	if err != nil {
		return fmt.Errorf("falha ao exportar banco de dados: %w", err)
	}

	return nil
}

// Restore restaura o banco de dados a partir de um arquivo .db externo
func (d *DB) Restore(sourcePath string) error {
	if sourcePath == "" {
		return fmt.Errorf("caminho de origem inválido")
	}

	// 1. Validar se o arquivo fonte é um SQLite legível
	srcConn, err := sql.Open("sqlite", sourcePath)
	if err != nil {
		return fmt.Errorf("arquivo de backup inválido: %w", err)
	}
	var testCount int
	err = srcConn.QueryRow(`SELECT count(*) FROM sqlite_master;`).Scan(&testCount)
	srcConn.Close()
	if err != nil {
		return fmt.Errorf("o arquivo selecionado não é um banco SQLite válido: %w", err)
	}

	// 2. Fechar conexão atual
	if d.conn != nil {
		_ = d.conn.Close()
	}

	// 3. Remover arquivos auxiliares de WAL e SHM
	_ = os.Remove(d.dbPath + "-wal")
	_ = os.Remove(d.dbPath + "-shm")

	// 4. Copiar o arquivo de backup para o destino real
	srcData, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("falha ao ler arquivo de backup: %w", err)
	}

	if err := os.WriteFile(d.dbPath, srcData, 0644); err != nil {
		return fmt.Errorf("falha ao sobrescrever banco de dados: %w", err)
	}

	// 5. Reabrir a conexão
	conn, err := sql.Open("sqlite", d.dbPath)
	if err != nil {
		return fmt.Errorf("falha ao reabrir banco após restauração: %w", err)
	}

	if _, err := conn.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		return fmt.Errorf("falha ao reconfigurar sqlite: %w", err)
	}

	d.conn = conn
	return d.migrate()
}

// Reset apaga todos os dados de todas as tabelas e re-executa a migração padrão de fábrica
func (d *DB) Reset() error {
	if d.conn == nil {
		return fmt.Errorf("conexão com o banco não inicializada")
	}

	tables := []string{
		"vps_servers",
		"profiles",
		"app_settings",
		"image_history",
		"saved_paths",
		"stacks",
		"container_configs",
	}

	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Desabilita foreign keys temporariamente durante a limpeza
	if _, err := tx.Exec(`PRAGMA foreign_keys = OFF;`); err != nil {
		return err
	}

	for _, table := range tables {
		if _, err := tx.Exec(fmt.Sprintf("DELETE FROM %s;", table)); err != nil {
			return fmt.Errorf("falha ao limpar tabela %s: %w", table, err)
		}
	}

	if _, err := tx.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Recria o perfil padrão inicial
	return d.migrate()
}

func (d *DB) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
