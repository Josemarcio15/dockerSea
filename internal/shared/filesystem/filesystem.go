package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

func EnsureDirectory(path string) error {
	if path == "" {
		return fmt.Errorf("caminho vazio")
	}
	return os.MkdirAll(filepath.Clean(path), 0o755)
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ReadFile(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("caminho vazio")
	}
	return os.ReadFile(filepath.Clean(path))
}

func WriteFile(path string, data []byte, mode os.FileMode) error {
	if path == "" {
		return fmt.Errorf("caminho vazio")
	}
	clean := filepath.Clean(path)
	if err := EnsureDirectory(filepath.Dir(clean)); err != nil {
		return err
	}
	return os.WriteFile(clean, data, mode)
}
