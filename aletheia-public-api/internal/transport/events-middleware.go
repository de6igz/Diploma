// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"aletheia-public-api/interfaces"
	v1 "aletheia-public-api/interfaces/types/v1"
	"context"
)

type EventsGetEvents func(ctx context.Context, pastHours int, userId int64) (items v1.EventsResponse, err error)
type EventsGetEventByID func(ctx context.Context, eventId string, userId int64) (resp v1.EventDetailResponse, err error)
type EventsGetMostRecentEvent func(ctx context.Context, userId int64, eventType string) (resp v1.EventDetailResponse, err error)
type EventsGetEventsByType func(ctx context.Context, eventType string, userId int64) (resp v1.EventsByEventTypeResponse, err error)

type MiddlewareEvents func(next interfaces.Events) interfaces.Events

type MiddlewareEventsGetEvents func(next EventsGetEvents) EventsGetEvents
type MiddlewareEventsGetEventByID func(next EventsGetEventByID) EventsGetEventByID
type MiddlewareEventsGetMostRecentEvent func(next EventsGetMostRecentEvent) EventsGetMostRecentEvent
type MiddlewareEventsGetEventsByType func(next EventsGetEventsByType) EventsGetEventsByType
