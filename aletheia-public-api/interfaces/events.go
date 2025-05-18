package interfaces

import (
	v1 "aletheia-public-api/interfaces/types/v1"
	"context"
)

// Events
// @tg http-server log metrics
// @tg http-prefix=v1
type Events interface {
	// GetEvents
	// @tg summary=`Получить события за последние часы`
	// @tg desc=`Возвращает список событий за указанный промежуток времени (параметр pastHours)`
	// @tg http-method=GET
	// @tg http-path=/events
	// @tg http-headers=userId|X-User-Id
	// @tg http-args=`pastHours|pastHours`
	GetEvents(ctx context.Context, pastHours int, userId int64) (items v1.EventsResponse, err error)

	// GetEventByID
	// @tg summary=`Получить детали события`
	// @tg desc=`Возвращает подробную информацию о событии по его ID`
	// @tg http-method=GET
	// @tg http-path=/event
	// @tg http-args=`eventId|eventId`
	// @tg http-headers=userId|X-User-Id
	GetEventByID(ctx context.Context, eventId string, userId int64) (resp v1.EventDetailResponse, err error)
	// GetMostRecentEvent
	// @tg summary=`Получить детали самого свежего события`
	// @tg desc=`Возвращает подробную информацию о последнем событии по его eventType`
	// @tg http-method=GET
	// @tg http-path=/event/recent
	// @tg http-args=`eventType|eventType`
	// @tg http-headers=userId|X-User-Id
	GetMostRecentEvent(ctx context.Context, userId int64, eventType string) (resp v1.EventDetailResponse, err error)

	// GetEventsByType
	// @tg summary=`Получить события по типу ивента`
	// @tg desc=`Возвращает инфу об ивентах по типу`
	// @tg http-method=GET
	// @tg http-path=/events_by_event_type
	// @tg http-args=`eventType|eventType`
	// @tg http-headers=userId|X-User-Id
	GetEventsByType(ctx context.Context, eventType string, userId int64) (resp v1.EventsByEventTypeResponse, err error)
}
