package db

import (
	"database/sql"
	"encoding/base64"
	"path/filepath"
	"testing"
)

func TestEncryptDecryptSecret(t *testing.T) {
	original := "MinhaSenhaSuperSecreta123!@#"

	encrypted, err := encryptSecret(original)
	if err != nil {
		t.Fatalf("Erro ao criptografar: %v", err)
	}

	if encrypted == original {
		t.Fatalf("O texto criptografado não deveria ser igual ao original")
	}

	decrypted, err := decryptSecret(encrypted)
	if err != nil {
		t.Fatalf("Erro ao descriptografar: %v", err)
	}

	if decrypted != original {
		t.Fatalf("Esperado %q, obteve %q", original, decrypted)
	}
}

func TestDecryptMismatchMachineGraceful(t *testing.T) {
	// Criptografado com chave de outra máquina ou chave desconhecida
	otherCipher, err := encryptionCipherWithKey("chave-de-outro-computador-qualquer")
	if err != nil {
		t.Fatalf("Erro ao gerar cifra com outra chave: %v", err)
	}

	nonce := make([]byte, otherCipher.NonceSize())
	ciphertext := otherCipher.Seal(nil, nonce, []byte("senha-secreta-em-outro-pc"), nil)
	encoded := append(nonce, ciphertext...)
	secretEnc := "enc:v1:" + string(encodeBase64(encoded))

	decrypted, err := decryptSecret(secretEnc)
	if err != nil {
		t.Fatalf("Não deveria retornar erro fatal, mas sim tratar graciosamente: %v", err)
	}
	if decrypted != "" {
		t.Fatalf("Esperado string vazia ao decodificar em máquina diferente, obteve %q", decrypted)
	}
}

func encodeBase64(data []byte) string {
	return base64.RawStdEncoding.EncodeToString(data)
}

func TestBackupAndReset(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_docksea.db")
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Erro ao abrir banco temporário: %v", err)
	}
	defer conn.Close()

	database := &DB{conn: conn, dbPath: dbPath}
	if err := database.migrate(); err != nil {
		t.Fatalf("Erro ao migrar banco de teste: %v", err)
	}

	// 1. Salvar dados de teste
	testServer := VpsServer{
		Name:           "Servidor de Teste Backup",
		ConnectionType: "ssh",
		Host:           "192.168.1.50",
		Port:           22,
		Username:       "root",
	}
	if err := database.SaveVpsServer(testServer); err != nil {
		t.Fatalf("Erro ao salvar servidor: %v", err)
	}

	// 2. Realizar Backup
	tmpBackup := filepath.Join(tmpDir, "backup_output.db")
	if err := database.Backup(tmpBackup); err != nil {
		t.Fatalf("Erro ao realizar backup: %v", err)
	}

	// 3. Resetar o banco
	if err := database.Reset(); err != nil {
		t.Fatalf("Erro ao resetar banco: %v", err)
	}

	// 4. Verificar que servidores foram zerados
	servers, err := database.ListVpsServers()
	if err != nil {
		t.Fatalf("Erro ao listar servidores após reset: %v", err)
	}
	if len(servers) != 0 {
		t.Fatalf("Esperava 0 servidores após reset, encontrou %d", len(servers))
	}

	// 5. Verificar que perfil padrão foi recriado
	profiles, err := database.ListProfiles()
	if err != nil || len(profiles) == 0 {
		t.Fatalf("Esperava recriação do perfil padrão após reset")
	}
}
