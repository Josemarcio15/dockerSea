package stacks

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// PackProjectDir compacta todo o conteúdo de um diretório de projeto em stream tar.gz.
// REGRA CRÍTICA: Não aplica .dockerignore ao pacote geral, garantindo que arquivos
// de runtime/compose (.env, configs, volumes) cheguem à VPS intactos.
// O .dockerignore será interpretado exclusivamente pelo BuildKit remoto no momento do build.
func PackProjectDir(ctx context.Context, folderPath string, writer io.Writer) error {
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
