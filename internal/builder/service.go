package builder

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
	"go-walis/internal/stacks"
)

type Folder struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type FolderListing struct {
	CurrentPath     string   `json:"currentPath"`
	ParentPath      string   `json:"parentPath"`
	Folders         []Folder `json:"folders"`
	HasDockerfile   bool     `json:"hasDockerfile"`
	HasDockerignore bool     `json:"hasDockerignore"`
}

type SavedPath struct {
	Path  string `json:"path"`
	Label string `json:"label"`
}

type Service struct {
	database *db.DB
	mu       sync.Mutex
}

func NewService(database *db.DB) *Service {
	return &Service{database: database}
}

func (s *Service) Browse(path string) (FolderListing, error) {
	if path == "" {
		var err error
		path, err = os.UserHomeDir()
		if err != nil {
			return FolderListing{}, fmt.Errorf("não foi possível localizar a pasta inicial: %w", err)
		}
	}
	path, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return FolderListing{}, err
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return FolderListing{}, fmt.Errorf("não foi possível ler a pasta: %w", err)
	}
	listing := FolderListing{CurrentPath: path, ParentPath: filepath.Dir(path)}
	if listing.ParentPath == path {
		listing.ParentPath = ""
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		listing.Folders = append(listing.Folders, Folder{Name: entry.Name(), Path: filepath.Join(path, entry.Name())})
	}
	sort.Slice(listing.Folders, func(i, j int) bool { return listing.Folders[i].Name < listing.Folders[j].Name })
	listing.HasDockerfile = fileExists(filepath.Join(path, "Dockerfile"))
	listing.HasDockerignore = fileExists(filepath.Join(path, ".dockerignore"))
	return listing, nil
}

func (s *Service) ListSavedPaths() []SavedPath {
	if s.database == nil {
		return nil
	}
	items, err := s.database.ListSavedPaths("")
	if err != nil {
		return nil
	}
	paths := make([]SavedPath, 0, len(items))
	for _, item := range items {
		paths = append(paths, SavedPath{Path: item.Path, Label: item.Label})
	}
	return paths
}

func (s *Service) SavePath(path, label string) error {
	path, err := filepath.Abs(filepath.Clean(path))
	if err != nil || !fileExists(path) {
		return fmt.Errorf("pasta inválida")
	}
	if label == "" {
		label = filepath.Base(path)
	}
	if s.database == nil {
		return fmt.Errorf("banco de dados não inicializado")
	}
	return s.database.SaveSavedPath(path, label, "")
}

func (s *Service) RemoveSavedPath(path string) error {
	if s.database == nil {
		return nil
	}
	return s.database.DeleteSavedPath(path, "")
}

func (s *Service) activeServer() (*db.VpsServer, error) {
	if s.database == nil {
		return nil, fmt.Errorf("banco de dados SQLite não inicializado")
	}
	servers, err := s.database.ListVpsServers()
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar servidores do SQLite: %w", err)
	}

	for _, server := range servers {
		if server.IsActive {
			return &server, nil
		}
	}
	return nil, fmt.Errorf("nenhuma VPS ativa selecionada")
}

func (s *Service) Build(folderPath, projectName, locale string) error {
	folderPath = filepath.Clean(folderPath)
	if !fileExists(filepath.Join(folderPath, "Dockerfile")) {
		return fmt.Errorf("Dockerfile não encontrado na pasta selecionada: %s", folderPath)
	}
	tag := strings.TrimSpace(projectName)
	if tag == "" {
		return fmt.Errorf("nome/tag da imagem não informado")
	}
	if !strings.Contains(tag, ":") {
		tag = tag + ":latest"
	}

	srv, err := s.activeServer()
	if err != nil {
		return err
	}

	app := application.Get()
	emit := func(name string, data map[string]interface{}) {
		if app != nil && app.Event != nil {
			app.Event.Emit(name, data)
		}
	}

	// Executa o build de forma assíncrona para streaming fluido
	go func() {
		ctx := context.Background()
		emitProgress := func(line string) {
			emit("builder:progress", map[string]interface{}{"line": line})
		}
		emitComplete := func(success bool, msg string) {
			result := map[string]interface{}{
				"success": success,
				"image":   tag,
			}
			if msg != "" {
				result["message"] = msg
			}
			emit("builder:complete", result)
		}

		emitProgress("→ [PREPARING] Conectando à VPS remota...")
		client, err := connection.NewClient(*srv)
		if err != nil {
			emitProgress(fmt.Sprintf("✗ Erro de conexão com a VPS: %v", err))
			emitComplete(false, err.Error())
			return
		}
		defer client.Close()

		// Cria diretório temporário para o build na VPS
		buildDir := fmt.Sprintf("$HOME/.docksea/builds/bld_%d", os.Getpid())
		if _, err := client.ExecCommand(fmt.Sprintf("mkdir -p %s", buildDir), false); err != nil {
			emitProgress(fmt.Sprintf("✗ Falha ao criar diretório temporário na VPS: %v", err))
			emitComplete(false, err.Error())
			return
		}
		defer func() {
			_, _ = client.ExecCommand(fmt.Sprintf("rm -rf %s", buildDir), false)
		}()

		// Empacota e envia via SSH Tar Stream (respeitando .dockerignore no BuildKit da VPS)
		emitProgress(fmt.Sprintf("→ [UPLOADING] Enviando arquivos do projeto para a VPS (%s)...", buildDir))
		extractCmd := fmt.Sprintf("tar -xzf - -C %s", buildDir)
		stdinPipe, waitExtract, err := client.StartCommandInput(ctx, extractCmd, false)
		if err != nil {
			emitProgress(fmt.Sprintf("✗ Falha ao abrir stream SSH para a VPS: %v", err))
			emitComplete(false, err.Error())
			return
		}

		if err := stacks.PackProjectDir(ctx, folderPath, stdinPipe); err != nil {
			_ = stdinPipe.Close()
			emitProgress(fmt.Sprintf("✗ Falha ao empacotar projeto: %v", err))
			emitComplete(false, err.Error())
			return
		}
		_ = stdinPipe.Close()

		if err := waitExtract(); err != nil {
			emitProgress(fmt.Sprintf("✗ Falha ao descompactar arquivos na VPS: %v", err))
			emitComplete(false, err.Error())
			return
		}
		emitProgress("→ [BUILDING] Arquivos recebidos na VPS. Iniciando docker build...")

		// Executa docker build na VPS
		dockerBin := strings.TrimSpace(srv.DockerPath)
		if dockerBin == "" {
			dockerBin = "docker"
		}
		buildCmd := fmt.Sprintf("cd %s && %s build -t %s .", buildDir, dockerBin, stacks.ShellQuote(tag))
		stdoutPipe, waitBuild, err := client.StartCommandOutput(ctx, buildCmd, false)
		if err != nil {
			emitProgress(fmt.Sprintf("✗ Falha ao iniciar docker build na VPS: %v", err))
			emitComplete(false, err.Error())
			return
		}

		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			emitProgress(scanner.Text())
		}

		if err := waitBuild(); err != nil {
			emitProgress(fmt.Sprintf("✗ Build falhou na VPS: %v", err))
			emitComplete(false, err.Error())
			return
		}

		emitProgress(fmt.Sprintf("✓ Imagem '%s' construída com sucesso na VPS!", tag))
		emitComplete(true, "")
	}()

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
