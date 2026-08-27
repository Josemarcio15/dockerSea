package builder

import (
	"os"
	"path/filepath"
)

func ProjectFileExists(folderPath, name string) bool {
	if folderPath == "" || name == "" {
		return false
	}
	info, err := os.Stat(filepath.Join(folderPath, name))
	return err == nil && !info.IsDir()
}
