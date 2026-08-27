package filesystem

import (
	"path/filepath"
	"testing"
)

func TestEnsureDirectoryAndExists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested")
	if err := EnsureDirectory(path); err != nil {
		t.Fatal(err)
	}
	ok, err := Exists(path)
	if err != nil || !ok {
		t.Fatalf("exists=%v err=%v", ok, err)
	}
}

func TestWriteAndReadFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "file.txt")
	if err := WriteFile(path, []byte("docksea"), 0o600); err != nil {
		t.Fatal(err)
	}
	data, err := ReadFile(path)
	if err != nil || string(data) != "docksea" {
		t.Fatalf("data=%q err=%v", data, err)
	}
}
