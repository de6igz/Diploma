package logs_errors

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Provider interface {
	GetAllEvents(ctx context.Context, userId int64, pastHours int) ([]*Event, error)
	GetEventById(ctx context.Context, userId int64, eventId string) (*Event, error)
	GetRecentEventByEventType(ctx context.Context, userId int64, eventType string) (*Event, error)
	GetEventsByProjectsId(ctx context.Context, userId, pastHours int64, projectId string) ([]*Event, error)
	GetEventsByEventType(ctx context.Context, userId int64, eventType string) ([]Event, error)
}

type provider struct {
	conn *sql.DB
}

func NewProvider(conn *sql.DB) Provider {
	return &provider{conn}
}

func (p *provider) GetAllEvents(ctx context.Context, userId int64, pastHours int) ([]*Event, error) {
	var query string
	var rows *sql.Rows
	var err error

	query = `
		SELECT event_type, service_name,  COUNT(*) AS event_count,user_id, language
FROM logs_events
where user_id = $1
	`
	addGroupBy := func() string {
		return " GROUP BY event_type,service_name,user_id,language;"
	}
	// Если указан период, добавляем условие по времени
	if pastHours > 0 {
		query += " AND timestamp >= NOW() - ($2 * INTERVAL '1 hour')"
		query += addGroupBy()
		rows, err = p.conn.QueryContext(ctx, query, userId, pastHours)
	} else {
		query += addGroupBy()
		rows, err = p.conn.QueryContext(ctx, query, userId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event

	for rows.Next() {
		var event Event
		//var usedRulesBytes []byte // временная переменная для хранения JSON из БД

		err = rows.Scan(
			&event.EventType,
			&event.ServiceName,
			&event.EventCount,
			&event.UserId,
			&event.Language,
		)
		if err != nil {
			return nil, err
		}

		//// Если в БД поле не пустое, парсим JSON в срез UsedRule
		//if len(usedRulesBytes) > 0 {
		//	if err := json.Unmarshal(usedRulesBytes, &event.UsedRules); err != nil {
		//		return nil, err
		//	}
		//} else {
		//	event.UsedRules = []UsedRule{}
		//}

		events = append(events, &event)
	}

	return events, nil
}

func (p *provider) GetEventById(ctx context.Context, userId int64, eventId string) (*Event, error) {
	query := `
		SELECT uuid, user_id, service_name, timestamp, log, event_type, used_rules, language, used_actions
		FROM logs_events
		WHERE user_id = $1 AND uuid = $2
	`

	var event Event
	var usedRulesBytes, usedActionsBytes []byte

	err := p.conn.QueryRowContext(ctx, query, userId, eventId).Scan(
		&event.Id,
		&event.UserId,
		&event.ServiceName,
		&event.Timestamp,
		&event.Log,
		&event.EventType,
		&usedRulesBytes,
		&event.Language,
		&usedActionsBytes,
	)
	if err != nil {
		return nil, err
	}

	// Десериализуем used_rules, если данные присутствуют
	if len(usedRulesBytes) > 0 {
		if err := json.Unmarshal(usedRulesBytes, &event.UsedRules); err != nil {
			return nil, err
		}
	} else {
		event.UsedRules = []UsedRule{}
	}

	// Десериализуем used_actions, если данные присутствуют
	if len(usedActionsBytes) > 0 {
		if err := json.Unmarshal(usedActionsBytes, &event.UsedActions); err != nil {
			return nil, err
		}
	} else {
		event.UsedActions = []UsedAction{}
	}

	return &event, nil
}

func (p *provider) GetRecentEventByEventType(ctx context.Context, userId int64, eventType string) (*Event, error) {

	query := `
		SELECT uuid, user_id, service_name, timestamp, log, event_type, used_rules, language, used_actions
		FROM logs_events
		WHERE user_id = $1 AND event_type = $2
		order by timestamp desc 
	`

	var event Event
	var usedRulesBytes, usedActionsBytes []byte

	err := p.conn.QueryRowContext(ctx, query, userId, eventType).Scan(
		&event.Id,
		&event.UserId,
		&event.ServiceName,
		&event.Timestamp,
		&event.Log,
		&event.EventType,
		&usedRulesBytes,
		&event.Language,
		&usedActionsBytes,
	)
	if err != nil {
		return nil, err
	}

	// Десериализуем used_rules, если данные присутствуют
	if len(usedRulesBytes) > 0 {
		if err := json.Unmarshal(usedRulesBytes, &event.UsedRules); err != nil {
			return nil, err
		}
	} else {
		event.UsedRules = []UsedRule{}
	}

	// Десериализуем used_actions, если данные присутствуют
	if len(usedActionsBytes) > 0 {
		if err := json.Unmarshal(usedActionsBytes, &event.UsedActions); err != nil {
			return nil, err
		}
	} else {
		event.UsedActions = []UsedAction{}
	}

	return &event, nil
}

func (p *provider) GetEventsByProjectsId(ctx context.Context, userId, pastHours int64, projectId string) ([]*Event, error) {
	var query string
	var rows *sql.Rows
	var err error

	query = `
		SELECT event_type, service_name,  COUNT(*) AS event_count,user_id, language
FROM logs_events
where user_id = $1 and project_id = $2
	`
	addGroupBy := func() string {
		return " GROUP BY event_type,service_name,user_id,language;"
	}
	// Если указан период, добавляем условие по времени
	if pastHours > 0 {
		query += " AND timestamp >= NOW() - ($3 * INTERVAL '1 hour')"
		query += addGroupBy()
		rows, err = p.conn.QueryContext(ctx, query, userId, projectId, pastHours)
	} else {
		query += addGroupBy()
		rows, err = p.conn.QueryContext(ctx, query, userId, projectId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event

	for rows.Next() {
		var event Event
		//var usedRulesBytes []byte // временная переменная для хранения JSON из БД

		err = rows.Scan(
			&event.EventType,
			&event.ServiceName,
			&event.EventCount,
			&event.UserId,
			&event.Language,
		)
		if err != nil {
			return nil, err
		}

		//// Если в БД поле не пустое, парсим JSON в срез UsedRule
		//if len(usedRulesBytes) > 0 {
		//	if err := json.Unmarshal(usedRulesBytes, &event.UsedRules); err != nil {
		//		return nil, err
		//	}
		//} else {
		//	event.UsedRules = []UsedRule{}
		//}

		events = append(events, &event)
	}

	return events, nil
}

func (p *provider) GetEventsByEventType(ctx context.Context, userId int64, eventType string) ([]Event, error) {
	var query string
	var rows *sql.Rows
	var err error

	query = `
		SELECT uuid, service_name, timestamp, language, log, event_type
		FROM logs_events
		where event_type = $1 and user_id = $2
		`
	rows, err = p.conn.QueryContext(ctx, query, eventType, userId)
	if err != nil {
		return nil, err
	}

	var events []Event
	for rows.Next() {
		var event Event

		err = rows.Scan(
			&event.Id,
			&event.ServiceName,
			&event.Timestamp,
			&event.Language,
			&event.Log,
			&event.EventType,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
