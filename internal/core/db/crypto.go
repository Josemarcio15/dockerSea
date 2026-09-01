package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"golang.org/x/crypto/argon2"
)

const (
	appSecretNamespace = "docksea-app-secrets-v1"
)

func getMachineKey() string {
	id, err := machineid.ProtectedID(appSecretNamespace)
	if err != nil || strings.TrimSpace(id) == "" {
		hostname, _ := os.Hostname()
		return "docksea-fallback-" + hostname
	}
	return id
}

func encryptionCipherWithKey(masterKey string) (cipher.AEAD, error) {
	salt := sha256.Sum256([]byte("docksea-v1-salt"))
	key := argon2.IDKey([]byte(masterKey), salt[:], 1, 64*1024, 4, 32)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func encryptionCipher() (cipher.AEAD, error) {
	return encryptionCipherWithKey(getMachineKey())
}

func encryptSecret(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	aead, err := encryptionCipher()
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aead.Seal(nil, nonce, []byte(value), nil)
	encoded := append(nonce, ciphertext...)
	return "enc:v1:" + base64.RawStdEncoding.EncodeToString(encoded), nil
}

func decryptSecret(value string) (string, error) {
	if value == "" || !strings.HasPrefix(value, "enc:v1:") {
		return value, nil
	}

	data, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(value, "enc:v1:"))
	if err != nil {
		return "", nil
	}

	// Tentar descriptografar com a chave exclusiva da máquina atual
	aead, err := encryptionCipher()
	if err == nil && len(data) >= aead.NonceSize() {
		plaintext, err := aead.Open(nil, data[:aead.NonceSize()], data[aead.NonceSize():], nil)
		if err == nil {
			return string(plaintext), nil
		}
	}

	// Se falhou (ex: banco importado de outra máquina/outro SO),
	// retorna string vazia para que o usuário redigite a senha na nova máquina
	return "", nil
}
