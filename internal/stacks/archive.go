package stacks

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PackProjectDir compacta todo o conteúdo de um diretório de projeto em stream tar.gz.
func PackProjectDir(ctx context.Context, folderPath string, writer io.Writer) error {
	return packProjectDir(ctx, folderPath, writer, nil)
}

// ReadDockerignorePatterns lê os padrões do .dockerignore da raiz do projeto.
// Padrões de negação são ignorados, como no fluxo de empacotamento do Builder.
func ReadDockerignorePatterns(folderPath string) []string {
	data, err := os.ReadFile(filepath.Join(folderPath, ".dockerignore"))
	if err != nil {
		return nil
	}
	var patterns []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "!") {
			patterns = append(patterns, line)
		}
	}
	return patterns
}

func PackProjectDirWithIgnore(ctx context.Context, folderPath string, writer io.Writer, ignored []string) error {
	return packProjectDir(ctx, folderPath, writer, ignored)
}

func packProjectDir(ctx context.Context, folderPath string, writer io.Writer, ignored []string) error {
	folderPath = filepath.Clean(folderPath)
	info, err := os.Stat(folderPath)
	if err != nil {
		return fmt.Errorf("pasta não encontrada: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("o caminho '%s' não é um diretório", folderPath)
	}

	gw := gzip.NewWriter(writer)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	err = filepath.Walk(folderPath, func(path string, fileInfo os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Checagem de cancelamento por Context
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		relPath, err := filepath.Rel(folderPath, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}
		for _, pattern := range ignored {
			if matchedIgnore(relPath, pattern) {
				if fileInfo.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Prepara cabeçalho Tar com permissões e metadados
		header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
		if err != nil {
			return fmt.Errorf("falha no cabeçalho tar para %s: %w", relPath, err)
		}

		// Garante caminhos no padrão Unix consistente (com barras normais)
		header.Name = filepath.ToSlash(relPath)

		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("falha ao escrever cabeçalho tar para %s: %w", relPath, err)
		}

		if fileInfo.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("falha ao ler arquivo %s: %w", path, err)
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return fmt.Errorf("falha ao copiar arquivo %s para o tar: %w", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("falha ao empacotar pasta: %w", err)
	}

	return nil
}

func matchedIgnore(path, pattern string) bool {
	pattern = strings.TrimSpace(strings.TrimPrefix(pattern, "!"))
	pattern = strings.TrimPrefix(pattern, "/")
	pattern = strings.TrimSuffix(pattern, "/")
	path = filepath.ToSlash(path)
	if ok, _ := filepath.Match(pattern, path); ok {
		return true
	}
	if !strings.Contains(pattern, "/") {
		if ok, _ := filepath.Match(pattern, filepath.Base(path)); ok {
			return true
		}
	}
	return strings.HasPrefix(path, pattern+"/") || strings.Contains(path, "/"+pattern+"/") || strings.HasSuffix(path, "/"+pattern)
}

func IsIgnoredPath(path, pattern string) bool { return matchedIgnore(path, pattern) }

// PackSingleYaml compacta um arquivo Compose YAML em memória para envio via SSH
func PackSingleYaml(yamlContent string, writer io.Writer) error {
	gw := gzip.NewWriter(writer)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	contentBytes := []byte(yamlContent)
	header := &tar.Header{
		Name: "compose.yaml",
		Mode: 0644,
		Size: int64(len(contentBytes)),
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := tw.Write(contentBytes); err != nil {
		return err
	}

	return nil
}
