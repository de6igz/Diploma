// handler.go
package events

import (
	"context"

	types "aletheia-public-api/interfaces/types/v1"
	"aletheia-public-api/internal/dataproviders/timescale"
	"aletheia-public-api/internal/dataproviders/timescale/repositories/logs_errors"
)

type Events struct {
	usecase    EventsUsecase
	serializer Serializer
}

// NewEvents создаёт новый обработчик событий с преднастроенными usecase и serializer.
func NewEvents() *Events {
	repo := logs_errors.NewProvider(timescale.GlobalInstance)
	usecase := NewEventsUsecase(repo)
	serializer := NewSerializer()
	return &Events{
		usecase:    usecase,
		serializer: serializer,
	}
}

// GetEvents получает события через usecase, затем сериализует их и оборачивает в EventsResponse.
func (e *Events) GetEvents(ctx context.Context, pastHours int, userId int64) (types.EventsResponse, error) {
	events, err := e.usecase.GetEvents(ctx, pastHours, userId)
	if err != nil {
		return types.EventsResponse{}, err
	}
	serialized := e.serializer.SerializeEvents(events)
	return types.EventsResponse{
		Events: serialized,
	}, nil
}

// GetEventByID получает событие по ID через usecase, затем сериализует его и оборачивает в EventDetailResponse.
func (e *Events) GetEventByID(ctx context.Context, eventId string, userId int64) (types.EventDetailResponse, error) {
	event, err := e.usecase.GetEventByID(ctx, eventId, userId)
	if err != nil {
		return types.EventDetailResponse{}, err
	}
	if event == nil {
		return types.EventDetailResponse{}, nil
	}
	serialized := e.serializer.SerializeEventDetail(*event)
	return types.EventDetailResponse{
		Event: &serialized,
	}, nil
}

// GetMostRecentEvent получает последнее событие определённого типа через usecase, затем сериализует его и оборачивает в EventDetailResponse.
func (e *Events) GetMostRecentEvent(ctx context.Context, userId int64, eventType string) (types.EventDetailResponse, error) {
	event, err := e.usecase.GetMostRecentEvent(ctx, userId, eventType)
	if err != nil {
		return types.EventDetailResponse{}, err
	}
	if event == nil {
		return types.EventDetailResponse{}, nil
	}
	serialized := e.serializer.SerializeEventDetail(*event)
	return types.EventDetailResponse{
		Event: &serialized,
	}, nil
}

func (e *Events) GetEventsByType(ctx context.Context, eventType string, userId int64) (resp types.EventsByEventTypeResponse, err error) {

	events, err := e.usecase.GetEventsByEventType(ctx, userId, eventType)
	if err != nil {
		return types.EventsByEventTypeResponse{}, err
	}

	serialized := e.serializer.SerializeEventsByEventType(events)

	return types.EventsByEventTypeResponse{
		Events: serialized,
	}, nil
}
