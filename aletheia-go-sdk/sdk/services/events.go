package services

import (
	"aletheia-go-sdk/sdk/config"
	"aletheia-go-sdk/sdk/models"
	"aletheia-go-sdk/sdk/transport"
	"aletheia-go-sdk/sdk/utils"
	"runtime"
	"runtime/debug"
	"time"
)

const (
	Golang = "golang"
)

type EventService struct {
	transport *transport.HTTPTransport
	config    *config.AletheiaConfig
	logger    config.Logger
}

func NewEventService(cfg config.AletheiaConfig, tp *transport.HTTPTransport, logger config.Logger) *EventService {
	return &EventService{
		transport: tp,
		config:    &cfg,
		logger:    logger,
	}
}

func (s *EventService) CaptureMessage(
	message string,
	level models.LogLevel,
	context map[string]interface{},
	tags map[string]string,
	userID string,
	userInfo map[string]interface{},
) {
	eventCtx := utils.BuildContext(context, userID, userInfo)

	req := models.EventRequest{
		ProjectId:    s.config.ProjectId,
		ServiceName:  s.config.ServiceName,
		Environment:  s.config.Environment,
		Version:      s.config.Version,
		GoVersion:    runtime.Version(),
		Os:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		EventType:    string(level),
		EventMessage: message,
		Tags:         utils.FormatTags(tags),
		Timestamp:    time.Now().Format(time.RFC3339),
		ContextJson:  eventCtx,
		Language:     Golang,
	}

	if err := s.transport.SendEvent(req); err != nil {
		s.logger.Errorf("Failed to send message: %v", err)
	}
}

func (s *EventService) CaptureException(
	exception error,
	eventMessage string,
	eventType string,
	context map[string]interface{},
	tags map[string]string,
	userID string,
	userInfo map[string]interface{},
	fields map[string]interface{},
) {
	stackTrace := string(debug.Stack())
	eventCtx := utils.BuildContext(context, userID, userInfo)

	var eventTypeModel string
	if eventType != "" {
		eventTypeModel = ":" + eventType
	}

	req := models.ErrorEvent{
		ProjectId:    s.config.ProjectId,
		ServiceName:  s.config.ServiceName,
		Environment:  s.config.Environment,
		Version:      s.config.Version,
		ErrorMessage: exception.Error(), //
		GoVersion:    runtime.Version(),
		Os:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		EventType:    string(models.EventError) + eventTypeModel,
		EventMessage: eventMessage,
		StackTrace:   stackTrace,
		Tags:         utils.FormatTags(tags),
		Timestamp:    time.Now().Format(time.RFC3339),
		ContextJson:  eventCtx,
		Language:     Golang,
		Fields:       fields,
	}

	if err := s.transport.SendEvent(req); err != nil {
		s.logger.Errorf("Failed to send exception: %v", err)
		return
	}
	s.logger.Infof("Captured exception: %v", exception)
}
