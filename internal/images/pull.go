package images

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go-walis/internal/core/connection"
)

// PullImageStream executa o pull de uma imagem via POST /images/create e envia o streaming de progresso linha a linha
func PullImageStream(ctx context.Context, client *connection.Client, imageName string, onProgress func(data DockerPullProgress)) error {
	httpClient := client.GetStreamHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	// Normalizar nome da imagem e tag
	fromImage := imageName
	tag := "latest"
	if strings.Contains(imageName, ":") {
		parts := strings.Split(imageName, ":")
		fromImage = parts[0]
		tag = parts[1]
	}

	endpoint := fmt.Sprintf("http://localhost/images/create?fromImage=%s&tag=%s", url.QueryEscape(fromImage), url.QueryEscape(tag))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao iniciar pull da imagem: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Docker API /images/create retornou status %d: %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	// Buffer expandido para mensagens JSON longas de progresso do Docker
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var lastError string

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			var progress DockerPullProgress
			if errJson := json.Unmarshal(line, &progress); errJson == nil {
				if progress.Error != "" {
					lastError = progress.Error
					onProgress(progress)
					return fmt.Errorf("%s", progress.Error)
				}
				onProgress(progress)
			}
		}
	}

	if errScan := scanner.Err(); errScan != nil {
		return fmt.Errorf("erro no streaming do Docker: %w", errScan)
	}

	if lastError != "" {
		return fmt.Errorf("%s", lastError)
	}

	return nil
}
