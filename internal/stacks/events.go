package stacks

import (
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Constantes de Fases do Pipeline de Deploy
const (
	PhasePreparing  = "preparing"
	PhaseUploading  = "uploading"
	PhaseValidating = "validating"
	PhaseBuilding   = "building"
	PhaseStarting   = "starting"
	PhaseChecking   = "checking"
	PhaseCleaning   = "cleaning"
	PhaseComplete   = "complete"
	PhaseError      = "error"
)

// Constantes de Nomes de Eventos Wails
const (
	EventDeployStarted  = "stacks:deploy:started"
	EventDeployProgress = "stacks:deploy:progress"
	EventDeployComplete = "stacks:deploy:complete"
	EventDeployFailed   = "stacks:deploy:failed"
)

// EventEmitter encapsula a emissão de eventos Wails v3
type EventEmitter interface {
	Emit(event string, data interface{})
}

type wailsEventEmitter struct{}

func (w *wailsEventEmitter) Emit(event string, data interface{}) {
	app := application.Get()
	if app != nil && app.Event != nil {
		app.Event.Emit(event, data)
	}
}

// DeployEventBroadcaster gerencia o envio de eventos de deploy
type DeployEventBroadcaster struct {
	emitter  EventEmitter
	stackID  string
	deployID string
}

// NewDeployEventBroadcaster cria uma nova instância do broadcaster
func NewDeployEventBroadcaster(stackID, deployID string, emitter EventEmitter) *DeployEventBroadcaster {
	if emitter == nil {
		emitter = &wailsEventEmitter{}
	}
	return &DeployEventBroadcaster{
		emitter:  emitter,
		stackID:  stackID,
		deployID: deployID,
	}
}

// EmitStarted notifica o início do deploy
func (b *DeployEventBroadcaster) EmitStarted(projectName string) {
	b.emitter.Emit(EventDeployStarted, map[string]interface{}{
		"stackId":     b.stackID,
		"deployId":    b.deployID,
		"projectName": projectName,
	})
}

// EmitProgress transmite uma atualização de progresso de uma fase
func (b *DeployEventBroadcaster) EmitProgress(phase, message string) {
	b.emitter.Emit(EventDeployProgress, StackProgressEvent{
		StackID:  b.stackID,
		DeployID: b.deployID,
		Phase:    phase,
		Message:  sanitizeLog(message),
	})
}

// EmitComplete notifica que o deploy foi finalizado com sucesso
func (b *DeployEventBroadcaster) EmitComplete(message string) {
	b.emitter.Emit(EventDeployComplete, StackProgressEvent{
		StackID:  b.stackID,
		DeployID: b.deployID,
		Phase:    PhaseComplete,
		Message:  sanitizeLog(message),
		Success:  true,
	})
}

// EmitFailed notifica que o deploy falhou em alguma fase
func (b *DeployEventBroadcaster) EmitFailed(phase, message string) {
	b.emitter.Emit(EventDeployFailed, StackProgressEvent{
		StackID:  b.stackID,
		DeployID: b.deployID,
		Phase:    phase,
		Message:  sanitizeLog(message),
		Success:  false,
	})
}

// sanitizeLog remove eventuais quebras ou dados desnecessários de segurança
func sanitizeLog(line string) string {
	return strings.TrimRight(line, "\r\n")
}
