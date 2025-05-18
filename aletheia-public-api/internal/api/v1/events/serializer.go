// serializer.go
package events

import (
	"strconv"

	types "aletheia-public-api/interfaces/types/v1"
	"aletheia-public-api/internal/dataproviders/timescale/repositories/logs_errors"
)

// Serializer определяет методы для преобразования внутренних событий в внешний формат.
type Serializer interface {
	SerializeEvents([]*logs_errors.Event) []*types.Event
	SerializeEventDetail(logs_errors.Event) types.EventDetail
	SerializeEventsByEventType([]logs_errors.Event) []types.EventByType
}

// defaultSerializer – стандартная реализация интерфейса.
type defaultSerializer struct{}

// NewSerializer возвращает новый экземпляр стандартного сериализатора.
func NewSerializer() Serializer {
	return &defaultSerializer{}
}

// SerializeEvents преобразует срез внутренних событий в []types.Event.
func (s *defaultSerializer) SerializeEvents(events []*logs_errors.Event) []*types.Event {
	res := make([]*types.Event, 0, len(events))
	for _, event := range events {
		res = append(res, &types.Event{
			EventType:   event.EventType,
			ServiceName: event.ServiceName,
			EventsCount: int(event.EventCount),
			Assignee:    strconv.FormatInt(event.UserId, 10),
			Language:    event.Language,
		})
	}
	return res
}

// SerializeEventDetail преобразует внутреннее событие в внешний тип EventDetail.
func (s *defaultSerializer) SerializeEventDetail(event logs_errors.Event) types.EventDetail {
	return types.EventDetail{
		ID:          event.Id,
		Log:         string(event.Log), // преобразуем json.RawMessage в строку
		ServiceName: event.ServiceName,
		EventType:   event.EventType,
		Timestamp:   event.Timestamp,
		UsedRules:   serializeUsedRules(event.UsedRules),
		UsedActions: serializeUsedActions(event.UsedActions),
		Language:    event.Language,
	}
}

func serializeUsedRules(rules []logs_errors.UsedRule) []types.UsedRule {
	result := make([]types.UsedRule, len(rules))
	for i, rule := range rules {
		result[i] = types.UsedRule{
			RuleId:   rule.RuleId,
			RuleName: rule.RuleName,
		}
	}
	return result
}

func serializeUsedActions(actions []logs_errors.UsedAction) []types.UsedAction {
	result := make([]types.UsedAction, len(actions))
	for i, act := range actions {
		result[i] = types.UsedAction{
			Type: act.Type,
			Params: types.UsedActionParams{
				Key:   act.Params.Key,
				Value: act.Params.Value,
			},
		}
	}
	return result
}

func (s *defaultSerializer) SerializeEventsByEventType(events []logs_errors.Event) []types.EventByType {
	result := make([]types.EventByType, len(events))
	for i, e := range events {
		result[i] = types.EventByType{
			Id:          e.Id,
			ServiceName: e.ServiceName,
			Timestamp:   e.Timestamp,
			Language:    e.Language,
			Log:         string(e.Log), // преобразуем json.RawMessage в строку
			EventType:   *e.EventType,
		}
	}
	return result
}
