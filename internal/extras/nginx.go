package extras

import (
	"fmt"
	"path/filepath"
	"strings"

	"go-walis/internal/core/connection"
)

// ListNginxSites lista os arquivos de configuração do Nginx em sites-available e sites-enabled
func listNginxSites(client *connection.Client) (NginxSitesResult, error) {
	cmd := `ls -1 /etc/nginx/sites-available 2>/dev/null || true; echo "---DOCKSEA-SEP---"; ls -1 /etc/nginx/sites-enabled 2>/dev/null || true`
	out, err := client.ExecCommand(cmd, false)
	if err != nil && strings.TrimSpace(out) == "" {
		return NginxSitesResult{}, fmt.Errorf("falha ao listar diretórios do Nginx: %w", err)
	}

	parts := strings.Split(out, "---DOCKSEA-SEP---")
	var available []string
	var enabled []string

	if len(parts) > 0 {
		for _, line := range strings.Split(strings.TrimSpace(parts[0]), "\n") {
			l := strings.TrimSpace(line)
			if l != "" && !strings.Contains(l, "No such file") && !strings.Contains(l, "cannot access") {
				available = append(available, l)
			}
		}
	}

	if len(parts) > 1 {
		for _, line := range strings.Split(strings.TrimSpace(parts[1]), "\n") {
			l := strings.TrimSpace(line)
			if l != "" && !strings.Contains(l, "No such file") && !strings.Contains(l, "cannot access") {
				enabled = append(enabled, l)
			}
		}
	}

	return NginxSitesResult{
		Available: available,
		Enabled:   enabled,
	}, nil
}

// readNginxSite lê o conteúdo de um arquivo de configuração do Nginx
func readNginxSite(client *connection.Client, filename string, directory string) (string, error) {
	cleanName := filepath.Base(strings.TrimSpace(filename))
	if cleanName == "" || cleanName == "." || cleanName == "/" {
		return "", fmt.Errorf("nome de arquivo inválido")
	}

	dir := "sites-available"
	if directory == "enabled" || directory == "sites-enabled" {
		dir = "sites-enabled"
	}

	path := fmt.Sprintf("/etc/nginx/%s/%s", dir, cleanName)
	cmd := fmt.Sprintf("cat %s 2>/dev/null || true", path)
	out, err := client.ExecCommand(cmd, false)
	if err != nil && strings.TrimSpace(out) == "" {
		return "", fmt.Errorf("falha ao ler arquivo %s: %w", path, err)
	}

	return out, nil
}

// saveNginxSite salva o conteúdo de um arquivo em /etc/nginx/sites-available
func saveNginxSite(client *connection.Client, filename string, content string) ExtraActionResult {
	cleanName := filepath.Base(strings.TrimSpace(filename))
	if cleanName == "" || cleanName == "." || cleanName == "/" {
		return ExtraActionResult{
			Success: false,
			Message: "nome de arquivo inválido",
		}
	}

	path := fmt.Sprintf("/etc/nginx/sites-available/%s", cleanName)
	encodedContent := strings.TrimSpace(strings.ReplaceAll(content, "\r\n", "\n"))
	cmd := fmt.Sprintf("mkdir -p /etc/nginx/sites-available && cat << 'EOF' > %s\n%s\nEOF", path, encodedContent)
	out, err := client.ExecCommand(cmd, true)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("erro ao salvar arquivo %s: %v", cleanName, err),
			Output:  out,
		}
	}

	return ExtraActionResult{
		Success: true,
		Message: fmt.Sprintf("Arquivo %s salvo com sucesso em sites-available.", cleanName),
		Output:  out,
	}
}

// enableNginxSite cria link simbólico de sites-available para sites-enabled
func enableNginxSite(client *connection.Client, filename string) ExtraActionResult {
	cleanName := filepath.Base(strings.TrimSpace(filename))
	if cleanName == "" || cleanName == "." || cleanName == "/" {
		return ExtraActionResult{
			Success: false,
			Message: "nome de arquivo inválido",
		}
	}

	src := fmt.Sprintf("/etc/nginx/sites-available/%s", cleanName)
	dest := fmt.Sprintf("/etc/nginx/sites-enabled/%s", cleanName)

	cmd := fmt.Sprintf("mkdir -p /etc/nginx/sites-enabled && ln -sf %s %s", src, dest)
	out, err := client.ExecCommand(cmd, true)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("erro ao habilitar site %s: %v", cleanName, err),
			Output:  out,
		}
	}

	return ExtraActionResult{
		Success: true,
		Message: fmt.Sprintf("Link simbólico para %s criado em sites-enabled com sucesso.", cleanName),
		Output:  out,
	}
}

// deleteNginxSite apaga o arquivo de sites-available ou o link de sites-enabled
func deleteNginxSite(client *connection.Client, filename string, directory string) ExtraActionResult {
	cleanName := filepath.Base(strings.TrimSpace(filename))
	if cleanName == "" || cleanName == "." || cleanName == "/" {
		return ExtraActionResult{
			Success: false,
			Message: "nome de arquivo inválido",
		}
	}

	dir := "sites-available"
	if directory == "enabled" || directory == "sites-enabled" {
		dir = "sites-enabled"
	}

	path := fmt.Sprintf("/etc/nginx/%s/%s", dir, cleanName)
	cmd := fmt.Sprintf("rm -f %s", path)
	out, err := client.ExecCommand(cmd, true)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: fmt.Sprintf("erro ao apagar %s: %v", cleanName, err),
			Output:  out,
		}
	}

	return ExtraActionResult{
		Success: true,
		Message: fmt.Sprintf("Arquivo %s apagado com sucesso de %s.", cleanName, dir),
		Output:  out,
	}
}

// testNginxConfig executa `nginx -t` para validar sintaxe das configurações
func testNginxConfig(client *connection.Client) ExtraActionResult {
	out, err := client.ExecCommand("nginx -t", true)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: "Erro no teste de sintaxe do Nginx.",
			Output:  out,
		}
	}

	return ExtraActionResult{
		Success: true,
		Message: "Sintaxe do Nginx válida e teste com sucesso.",
		Output:  out,
	}
}

// restartNginx reinicia ou recarrega o serviço Nginx
func restartNginx(client *connection.Client) ExtraActionResult {
	cmd := "systemctl reload nginx 2>/dev/null || systemctl restart nginx 2>/dev/null || service nginx restart"
	out, err := client.ExecCommand(cmd, true)
	if err != nil {
		return ExtraActionResult{
			Success: false,
			Message: "Erro ao reiniciar serviço do Nginx.",
			Output:  out,
		}
	}

	return ExtraActionResult{
		Success: true,
		Message: "Nginx reiniciado/recarregado com sucesso.",
		Output:  out,
	}
}

// getNginxLogs obtém os logs mais recentes de access.log e error.log do Nginx
func getNginxLogs(client *connection.Client, lines int) (string, error) {
	if lines <= 0 {
		lines = 100
	}

	cmd := fmt.Sprintf(`echo "=== /var/log/nginx/error.log (Últimas %d linhas) ==="; tail -n %d /var/log/nginx/error.log 2>/dev/null || echo "(sem error.log)"; echo ""; echo "=== /var/log/nginx/access.log (Últimas %d linhas) ==="; tail -n %d /var/log/nginx/access.log 2>/dev/null || echo "(sem access.log)"`, lines, lines, lines, lines)
	out, err := client.ExecCommand(cmd, true)
	if err != nil && strings.TrimSpace(out) == "" {
		return "", fmt.Errorf("falha ao consultar logs do Nginx: %w", err)
	}

	return out, nil
}
