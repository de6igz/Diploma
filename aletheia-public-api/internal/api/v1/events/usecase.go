// usecase.go
package events

import (
	"context"

	"aletheia-public-api/internal/dataproviders/timescale/repositories/logs_errors"
)

// EventsUsecase описывает логику получения данных из репозитория.
type EventsUsecase interface {
	GetEvents(ctx context.Context, pastHours int, userId int64) ([]*logs_errors.Event, error)
	GetEventsByProjectId(ctx context.Context, pastHours int64, userId int64, projectId string) ([]*logs_errors.Event, error)
	GetEventByID(ctx context.Context, eventId string, userId int64) (*logs_errors.Event, error)
	GetMostRecentEvent(ctx context.Context, userId int64, eventType string) (*logs_errors.Event, error)
	GetEventsByEventType(ctx context.Context, userId int64, eventType string) ([]logs_errors.Event, error)
}

type eventsUsecase struct {
	logsRepo logs_errors.Provider
}

// NewEventsUsecase создаёт новый usecase с инъекцией репозитория.
func NewEventsUsecase(repo logs_errors.Provider) EventsUsecase {
	return &eventsUsecase{logsRepo: repo}
}

func (uc *eventsUsecase) GetEvents(ctx context.Context, pastHours int, userId int64) ([]*logs_errors.Event, error) {
	return uc.logsRepo.GetAllEvents(ctx, userId, pastHours)
}

func (uc *eventsUsecase) GetEventByID(ctx context.Context, eventId string, userId int64) (*logs_errors.Event, error) {
	return uc.logsRepo.GetEventById(ctx, userId, eventId)
}

func (uc *eventsUsecase) GetMostRecentEvent(ctx context.Context, userId int64, eventType string) (*logs_errors.Event, error) {
	if eventType == "" {
		return nil, nil
	}
	return uc.logsRepo.GetRecentEventByEventType(ctx, userId, eventType)
}

func (uc *eventsUsecase) GetEventsByProjectId(ctx context.Context, pastHours int64, userId int64, projectId string) ([]*logs_errors.Event, error) {
	return uc.logsRepo.GetEventsByProjectsId(ctx, userId, pastHours, projectId)
}

func (uc *eventsUsecase) GetEventsByEventType(ctx context.Context, userId int64, eventType string) ([]logs_errors.Event, error) {
	if eventType == "" {
		return nil, nil
	}
	events, err := uc.logsRepo.GetEventsByEventType(ctx, userId, eventType)
	if err != nil {
		return nil, err
	}

	return events, nil

}
