package stacks

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestPackProjectDir_Success(t *testing.T) {
	// Creates estrutura temporária de projeto de teste
	tempDir, err := os.MkdirTemp("", "docksea_archive_test_*")
	if err != nil {
		t.Fatalf("erro ao criar temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Creates arquivos diversos: compose.yaml, Dockerfile, .env, .dockerignore, subpastas com binários e arquivos Unicode
	files := map[string]string{
		"compose.yaml":                 "services:\n  app:\n    image: test",
		"Dockerfile":                   "FROM alpine\nCMD [\"echo\", \"hello\"]",
		".dockerignore":                "node_modules\n.env*\n",
		".env.production":              "DATABASE_URL=postgres://...",
		"config/app.json":              "{\"port\": 8080}",
		"src/index.js":                 "console.log('test');",
		"nested/deep/arquivo_ação.txt": "conteúdo com acentuação e caracteres especiais UTF-8",
	}

	for relPath, content := range files {
		fullPath := filepath.Join(tempDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("erro ao criar pasta: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("erro ao criar arquivo: %v", err)
		}
	}

	// Adiciona um arquivo binário
	binPath := filepath.Join(tempDir, "bin/app.bin")
	_ = os.MkdirAll(filepath.Dir(binPath), 0755)
	_ = os.WriteFile(binPath, []byte{0x00, 0xFF, 0xFE, 0x42, 0x13, 0x37}, 0755)

	var buf bytes.Buffer
	ctx := context.Background()

	if err := PackProjectDir(ctx, tempDir, &buf); err != nil {
		t.Fatalf("PackProjectDir falhou: %v", err)
	}

	// Extrai e valida se todos os arquivos foram incluídos fielmente no tar.gz
	gr, err := gzip.NewReader(&buf)
	if err != nil {
		t.Fatalf("erro ao abrir gzip reader: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	foundFiles := make(map[string]bool)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("erro ao ler tar: %v", err)
		}
		foundFiles[header.Name] = true
	}

	// REGRA CRÍTICA: .dockerignore e .env.production DEVEM estar presentes no pacote geral
	expectedFiles := []string{
		"compose.yaml",
		"Dockerfile",
		".dockerignore",
		".env.production",
		"config/app.json",
		"src/index.js",
		"nested/deep/arquivo_ação.txt",
		"bin/app.bin",
	}

	for _, expected := range expectedFiles {
		if !foundFiles[expected] {
			t.Errorf("arquivo esperado '%s' não encontrado no tar.gz gerado", expected)
		}
	}
}

func TestPackProjectDir_ContextCancellation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "docksea_cancel_test_*")
	if err != nil {
		t.Fatalf("erro ao criar temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	for i := 0; i < 50; i++ {
		p := filepath.Join(tempDir, filepath.Join("sub", "file.txt"))
		_ = os.MkdirAll(filepath.Dir(p), 0755)
		_ = os.WriteFile(p, []byte("data"), 0644)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancela imediatamente antes de iniciar

	var buf bytes.Buffer
	err = PackProjectDir(ctx, tempDir, &buf)
	if err == nil {
		t.Fatalf("esperava erro de cancelamento de contexto, mas retornou nil")
	}
}

func TestPackSingleYaml_Success(t *testing.T) {
	yamlContent := "services:\n  redis:\n    image: redis:alpine\n"
	var buf bytes.Buffer

	if err := PackSingleYaml(yamlContent, &buf); err != nil {
		t.Fatalf("PackSingleYaml falhou: %v", err)
	}

	gr, err := gzip.NewReader(&buf)
	if err != nil {
		t.Fatalf("erro gzip reader: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	header, err := tr.Next()
	if err != nil {
		t.Fatalf("erro ao ler cabeçalho tar: %v", err)
	}

	if header.Name != "compose.yaml" {
		t.Fatalf("esperado nome 'compose.yaml', obtido '%s'", header.Name)
	}

	extractedContent, err := io.ReadAll(tr)
	if err != nil {
		t.Fatalf("erro ao ler conteúdo do tar: %v", err)
	}

	if string(extractedContent) != yamlContent {
		t.Fatalf("conteúdo extraído não corresponde ao original")
	}
}
