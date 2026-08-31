package containers

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go-walis/internal/core/db"
	sharedDocker "go-walis/internal/shared/docker"
	sharedevents "go-walis/internal/shared/events"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type ContainerService struct {
	database *db.DB
	streamMu sync.Mutex
	cancelFn map[string]context.CancelFunc
}

func NewContainerService(database *db.DB) *ContainerService {
	return &ContainerService{
		database: database,
		cancelFn: make(map[string]context.CancelFunc),
	}
}

// StartEventsStream inicia a escuta em tempo real dos eventos da VPS e emite 'docker:container:event' no Wails
func (s *ContainerService) StartEventsStream(server db.VpsServer) error {
	s.streamMu.Lock()
	defer s.streamMu.Unlock()

	if cancel, exists := s.cancelFn[server.ID]; exists {
		cancel()
		delete(s.cancelFn, server.ID)
	}

	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return fmt.Errorf("falha ao conectar no servidor para eventos: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFn[server.ID] = cancel

	go func() {
		defer client.Close()
		fmt.Printf("[DockerEvents] Iniciando streaming de eventos para o servidor %s (%s)...\n", server.Name, server.ID)
		err := StreamEvents(ctx, client, "container", func(ev DockerEvent) {
			fmt.Printf("[DockerEvents] Evento recebido: action=%s name=%s id=%s\n", ev.Action, ev.Actor.Attributes["name"], ev.Actor.ID)
			app := application.Get()
			if app != nil && app.Event != nil {
				app.Event.Emit(sharedevents.ContainerEvent, map[string]interface{}{
					"serverId": server.ID,
					"type":     ev.Type,
					"action":   ev.Action,
					"id":       ev.Actor.ID,
					"name":     ev.Actor.Attributes["name"],
					"image":    ev.Actor.Attributes["image"],
					"time":     ev.Time,
				})
			}
		})
		if err != nil && ctx.Err() == nil {
			fmt.Printf("[DockerEvents] Stream finalizado com erro no servidor %s: %v\n", server.Name, err)
		}
	}()

	return nil
}

// StopEventsStream interrompe a escuta de eventos da VPS
func (s *ContainerService) StopEventsStream(serverId string) {
	s.streamMu.Lock()
	defer s.streamMu.Unlock()

	if cancel, exists := s.cancelFn[serverId]; exists {
		cancel()
		delete(s.cancelFn, serverId)
	}
}

// ListContainers obtém a lista de containers da VPS
func (s *ContainerService) ListContainers(server db.VpsServer, all bool) ([]Container, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return ListContainers(client, all)
}

func (s *ContainerService) CreateContainer(server db.VpsServer, input CreateContainerInput) ContainerActionResult {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Image) == "" {
		return ContainerActionResult{Success: false, Message: "nome e imagem são obrigatórios"}
	}
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return ContainerActionResult{Success: false, Message: fmt.Sprintf("falha ao conectar no servidor: %v", err)}
	}
	defer client.Close()
	if err := CreateContainer(client, input); err != nil {
		return ContainerActionResult{Success: false, Message: err.Error()}
	}
	return ContainerActionResult{Success: true, Message: "Container criado com sucesso!"}
}

// ExecuteAction executa ações em lote: "start" | "stop" | "restart" | "rm"
func (s *ContainerService) ExecuteAction(server db.VpsServer, actionType string, containerNames []string) ContainerActionResult {
	if message := ValidateContainerAction(actionType); message != "" {
		return ContainerActionResult{Success: false, Message: message}
	}
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return ContainerActionResult{
			Success: false,
			Message: fmt.Sprintf("falha ao conectar no servidor: %v", err),
		}
	}
	defer client.Close()

	var errList []string
	for _, name := range containerNames {
		trimmed := NormalizeContainerName(name)
		if trimmed == "" {
			continue
		}

		var opErr error
		switch actionType {
		case "start":
			opErr = StartContainer(client, trimmed)
		case "stop":
			opErr = StopContainer(client, trimmed)
		case "restart":
			opErr = RestartContainer(client, trimmed)
		case "rm":
			opErr = RemoveContainer(client, trimmed, true)
		}

		if opErr != nil {
			errList = append(errList, fmt.Sprintf("%s: %v", trimmed, opErr))
		}
	}

	if len(errList) > 0 {
		return ContainerActionResult{
			Success: false,
			Message: strings.Join(errList, "\n"),
			Errors:  errList,
		}
	}

	return ContainerActionResult{
		Success: true,
		Message: "Operação executada com sucesso!",
	}
}

// GetLogs obtém os logs do container especificado
func (s *ContainerService) GetLogs(server db.VpsServer, containerName string, tail int) (string, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return "", fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return GetContainerLogs(client, containerName, tail)
}

// InspectContainer obtém detalhes completos sob demanda de um container específico
func (s *ContainerService) InspectContainer(server db.VpsServer, id string) (*RawDockerInspect, error) {
	client, err := sharedDocker.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()

	return InspectContainer(client, id)
}
