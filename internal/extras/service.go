package extras

import (
	"fmt"
	"path/filepath"
	"strings"

	"go-walis/internal/core/db"
	sharedDocker "go-walis/internal/shared/docker"
)

type ExtraService struct {
	database *db.DB
}

func (s *ExtraService) ListDeployTempFilesAt(server db.VpsServer, path string) ([]DeployTempFile, error) {
	if path == "" {
		path = "$HOME/.docksea"
	}
	validDockseaPath := strings.HasPrefix(path, "$HOME/.docksea") ||
		(strings.HasPrefix(path, "/") && (strings.Contains(path, "/.docksea/") || strings.HasSuffix(path, "/.docksea")))
	if !validDockseaPath {
		return nil, fmt.Errorf("caminho de deploy inválido")
	}
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	remotePath := strings.ReplaceAll(path, `"`, `\"`)
	out, err := client.ExecCommand(fmt.Sprintf(`find "%s" -mindepth 1 -maxdepth 1 -printf '%%p|%%s|%%y\n' 2>/dev/null`, remotePath), false)
	if err != nil {
		return nil, err
	}
	var files []DeployTempFile
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		parts := strings.SplitN(line, "|", 3)
		if len(parts) == 3 {
			files = append(files, DeployTempFile{Path: parts[0], Size: parts[1], IsDir: parts[2] == "d"})
		}
	}
	return files, nil
}

func (s *ExtraService) ListDeployTempFiles(server db.VpsServer) ([]DeployTempFile, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()
	out, err := client.ExecCommand("for dir in \"$HOME/.docksea/lifecycle\" \"$HOME/.docksea/deploys\"; do find \"$dir\" -mindepth 1 -maxdepth 1 -printf '%p|%s|%y\\n' 2>/dev/null; done", false)
	if err != nil {
		return nil, err
	}
	var files []DeployTempFile
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		parts := strings.SplitN(line, "|", 3)
		if len(parts) == 3 && parts[0] != "" {
			files = append(files, DeployTempFile{Path: parts[0], Size: parts[1], IsDir: parts[2] == "d"})
		}
	}
	return files, nil
}

func (s *ExtraService) CleanDeployTempFiles(server db.VpsServer) ExtraActionResult {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return ExtraActionResult{Message: err.Error()}
	}
	defer client.Close()
	_, err = client.ExecCommand("for dir in \"$HOME/.docksea/lifecycle\" \"$HOME/.docksea/deploys\"; do find \"$dir\" -type f -mtime +1 -delete 2>/dev/null; find \"$dir\" -type d -empty -delete 2>/dev/null; done", false)
	if err != nil {
		return ExtraActionResult{Message: err.Error()}
	}
	return ExtraActionResult{Success: true, Message: "Arquivos temporários removidos."}
}

func (s *ExtraService) DeleteDeployTempPath(server db.VpsServer, path string) ExtraActionResult {
	valid := strings.HasPrefix(path, "$HOME/.docksea") || (strings.HasPrefix(path, "/") && strings.Contains(path, "/.docksea/"))
	if !valid {
		return ExtraActionResult{Message: "caminho de deploy inválido"}
	}
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return ExtraActionResult{Message: err.Error()}
	}
	defer client.Close()
	quoted := strings.ReplaceAll(path, `'`, `'\''`)
	command := fmt.Sprintf("if [ -d '%s' ]; then rm -rf -- '%s'; else rm -f -- '%s'; fi", quoted, quoted, quoted)
	if _, err := client.ExecCommand(command, false); err != nil {
		return ExtraActionResult{Message: err.Error()}
	}
	return ExtraActionResult{Success: true, Message: "Item removido com sucesso."}
}

func NewExtraService(database *db.DB) *ExtraService {
	return &ExtraService{database: database}
}

// ListNginxSites lista os arquivos de configuração do Nginx em sites-available e sites-enabled
func (s *ExtraService) ListNginxSites(server db.VpsServer) (NginxSitesResult, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return NginxSitesResult{}, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return listNginxSites(client)
}

// ReadNginxSite lê o conteúdo de um arquivo de configuração do Nginx
func (s *ExtraService) ReadNginxSite(server db.VpsServer, filename string, directory string) (string, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return "", fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return readNginxSite(client, filename, directory)
}

// SaveNginxSite salva o conteúdo de um arquivo em /etc/nginx/sites-available
func (s *ExtraService) SaveNginxSite(server db.VpsServer, filename string, content string) ExtraActionResult {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	return saveNginxSite(client, filename, content)
}

// EnableNginxSite cria link simbólico de sites-available para sites-enabled
func (s *ExtraService) EnableNginxSite(server db.VpsServer, filename string) ExtraActionResult {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	return enableNginxSite(client, filename)
}

// DeleteNginxSite apaga o arquivo de sites-available ou o link de sites-enabled
func (s *ExtraService) DeleteNginxSite(server db.VpsServer, filename string, directory string) ExtraActionResult {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	return deleteNginxSite(client, filename, directory)
}

// TestNginxConfig executa `nginx -t` para validar sintaxe das configurações
func (s *ExtraService) TestNginxConfig(server db.VpsServer) ExtraActionResult {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	return testNginxConfig(client)
}

// RestartNginx reinicia ou recarrega o serviço Nginx
func (s *ExtraService) RestartNginx(server db.VpsServer) ExtraActionResult {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	return restartNginx(client)
}

func (s *ExtraService) ListNginxLogFiles(server db.VpsServer) ([]NginxLogFile, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	out, err := client.ExecCommand("find /var/log/nginx -maxdepth 1 -type f -printf '%f|%s\\n' 2>/dev/null | sort", true)
	if err != nil {
		return nil, err
	}
	var result []NginxLogFile
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		p := strings.SplitN(line, "|", 2)
		if len(p) == 2 {
			result = append(result, NginxLogFile{Name: p[0], Size: p[1], Compressed: strings.HasSuffix(p[0], ".gz")})
		}
	}
	return result, nil
}

func (s *ExtraService) ReadNginxLogFile(server db.VpsServer, name string, lines int) (string, error) {
	clean := filepath.Base(strings.TrimSpace(name))
	if clean == "." || clean == "" || clean != name {
		return "", fmt.Errorf("arquivo de log inválido")
	}
	if lines <= 0 {
		lines = 150
	}
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return "", err
	}
	defer client.Close()
	cmd := fmt.Sprintf("zcat /var/log/nginx/%s 2>/dev/null | tail -n %d", clean, lines)
	if !strings.HasSuffix(clean, ".gz") {
		cmd = fmt.Sprintf("tail -n %d /var/log/nginx/%s 2>/dev/null", lines, clean)
	}
	out, err := client.ExecCommand(cmd, true)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(out) == "" {
		return "Arquivo existente, porém sem registros.", nil
	}
	return out, nil
}

// GetNginxLogs obtém os logs mais recentes de access.log e error.log do Nginx
func (s *ExtraService) GetNginxLogs(server db.VpsServer, lines int) (string, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return "", fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return getNginxLogs(client, lines)
}

// ListListeningPorts consulta portas TCP e UDP em escuta no servidor usando ss / netstat / lsof
func (s *ExtraService) ListListeningPorts(server db.VpsServer) ([]PortEntry, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return listListeningPorts(client)
}
