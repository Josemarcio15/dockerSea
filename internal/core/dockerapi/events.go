package dockerapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go-walis/internal/core/connection"
)

// StreamEvents abre uma conexão de longa duração com o Docker Daemon e chama o handler a cada evento
func StreamEvents(ctx context.Context, client *connection.Client, eventType string, onEvent func(event DockerEvent)) error {
	httpClient := client.GetStreamHttpClient()
	if httpClient == nil {
		return fmt.Errorf("cliente HTTP do Docker não disponível")
	}

	endpoint := "http://localhost/events"
	if eventType != "" {
		endpoint = fmt.Sprintf("http://localhost/events?filters=%%7B%%22type%%22:%%5B%%22%s%%22%%5D%%7D", url.QueryEscape(eventType))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao abrir stream de eventos do Docker: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Docker /events retornou status %d: %s", resp.StatusCode, string(body))
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var ev DockerEvent
			if err := decoder.Decode(&ev); err != nil {
				if err == io.EOF || ctx.Err() != nil {
					return nil
				}
				return fmt.Errorf("erro na leitura do stream de eventos: %w", err)
			}
			onEvent(ev)
		}
	}
}
