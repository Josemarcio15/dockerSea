package images

import (
	"context"
	"fmt"
	"strings"

	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type ImageService struct {
	database *db.DB
}

func NewImageService(database *db.DB) *ImageService {
	return &ImageService{
		database: database,
	}
}

// ListImages obtém a lista de imagens da VPS e correlaciona os containers em uso
func (s *ImageService) ListImages(server db.VpsServer) ([]DockerImage, error) {
	client, err := connection.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return ListImages(client, true)
}

// DeleteImages remove uma ou mais imagens pelo ID ou nome
func (s *ImageService) DeleteImages(server db.VpsServer, imageIds []string) ImageActionResult {
	client, err := connection.NewClient(server)
	if err != nil {
		return ImageActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	var errList []string
	deletedCount := 0

	for _, id := range imageIds {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}

		if err := RemoveImage(client, trimmed, true); err != nil {
			errList = append(errList, fmt.Sprintf("%s: %v", trimmed, err))
		} else {
			deletedCount++
		}
	}

	if len(errList) > 0 {
		return ImageActionResult{
			Success: deletedCount > 0,
			Message: strings.Join(errList, "\n"),
			Count:   deletedCount,
			Errors:  errList,
		}
	}

	return ImageActionResult{
		Success: true,
		Message: "Imagens removidas com sucesso!",
		Count:   deletedCount,
	}
}

// PullImage inicia o download de uma imagem emitindo eventos via Wails Runtime
func (s *ImageService) PullImage(server db.VpsServer, imageName string, profileId string) ImageActionResult {
	trimmed := strings.TrimSpace(imageName)
	if trimmed == "" {
		return ImageActionResult{
			Success: false,
			Message: "Nome da imagem não informado",
		}
	}

	client, err := connection.NewClient(server)
	if err != nil {
		return ImageActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	app := application.Get()

	// Grava no histórico SQLite
	_ = s.database.AddImageHistory(trimmed, profileId)

	ctx := context.Background()
	err = PullImageStream(ctx, client, trimmed, func(progress DockerPullProgress) {
		if app != nil && app.Event != nil {
			app.Event.Emit("docker:image:pull:progress", map[string]interface{}{
				"serverId":       server.ID,
				"imageName":      trimmed,
				"id":             progress.ID,
				"status":         progress.Status,
				"progress":       progress.Progress,
				"progressDetail": progress.ProgressDetail,
				"error":          progress.Error,
				"line":           progress.Line,
			})
		}
	})

	if err != nil {
		if app != nil && app.Event != nil {
			app.Event.Emit("docker:image:pull:complete", map[string]interface{}{
				"serverId":  server.ID,
				"imageName": trimmed,
				"success":   false,
				"message":   err.Error(),
			})
		}
		return ImageActionResult{
			Success: false,
			Message: err.Error(),
		}
	}

	if app != nil && app.Event != nil {
		app.Event.Emit("docker:image:pull:complete", map[string]interface{}{
			"serverId":  server.ID,
			"imageName": trimmed,
			"success":   true,
			"message":   fmt.Sprintf("Imagem '%s' baixada com sucesso!", trimmed),
		})
	}

	return ImageActionResult{
		Success: true,
		Message: fmt.Sprintf("Imagem '%s' baixada com sucesso!", trimmed),
	}
}

// TransferImages transfere imagens entre dois servidores VPS emitindo eventos de progresso em tempo real
func (s *ImageService) TransferImages(srcServer, dstServer db.VpsServer, imageIds []string) ImageActionResult {
	if len(imageIds) == 0 {
		return ImageActionResult{
			Success: false,
			Message: "Nenhuma imagem selecionada para transferência",
		}
	}
	if srcServer.ID == dstServer.ID {
		return ImageActionResult{
			Success: false,
			Message: "Origem e destino devem ser servidores diferentes",
		}
	}

	srcClient, err := connection.NewClient(srcServer)
	if err != nil {
		return ImageActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor de origem (%s): %v", srcServer.Name, err),
		}
	}
	defer srcClient.Close()

	dstClient, err := connection.NewClient(dstServer)
	if err != nil {
		return ImageActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor de destino (%s): %v", dstServer.Name, err),
		}
	}
	defer dstClient.Close()

	app := application.Get()
	ctx := context.Background()

	err = TransferImages(ctx, srcClient, dstClient, imageIds, func(progress DockerTransferProgress) {
		if app != nil && app.Event != nil {
			app.Event.Emit("docker:image:transfer:progress", map[string]interface{}{
				"srcServerId":         srcServer.ID,
				"dstServerId":         dstServer.ID,
				"stage":               progress.Stage,
				"bytesSent":           progress.BytesSent,
				"formattedBytes":      progress.FormattedBytes,
				"totalBytes":          progress.TotalBytes,
				"formattedTotalBytes": progress.FormattedTotalBytes,
				"speed":               progress.Speed,
				"percent":             progress.Percent,
				"message":             progress.Message,
			})
		}
	})

	if err != nil {
		if app != nil && app.Event != nil {
			app.Event.Emit("docker:image:transfer:complete", map[string]interface{}{
				"srcServerId": srcServer.ID,
				"dstServerId": dstServer.ID,
				"success":     false,
				"message":     err.Error(),
			})
		}
		return ImageActionResult{
			Success: false,
			Message: fmt.Sprintf("Erro durante transferência: %v", err),
		}
	}

	if app != nil && app.Event != nil {
		app.Event.Emit("docker:image:transfer:complete", map[string]interface{}{
			"srcServerId": srcServer.ID,
			"dstServerId": dstServer.ID,
			"success":     true,
			"message":     fmt.Sprintf("%d imagem(ns) transferida(s) com sucesso!", len(imageIds)),
			"count":       len(imageIds),
		})
	}

	return ImageActionResult{
		Success: true,
		Message: fmt.Sprintf("%d imagem(ns) transferida(s) com sucesso!", len(imageIds)),
		Count:   len(imageIds),
	}
}

// ListHistory lista histórico de downloads da imagem
func (s *ImageService) ListHistory(profileId string) ([]db.ImageHistoryItem, error) {
	return s.database.ListImageHistory(profileId)
}

// DeleteHistory remove itens específicos do histórico
func (s *ImageService) DeleteHistory(ids []string) error {
	return s.database.DeleteImageHistory(ids)
}

// ClearHistory limpa todo o histórico de downloads do perfil
func (s *ImageService) ClearHistory(profileId string) error {
	return s.database.ClearImageHistory(profileId)
}
