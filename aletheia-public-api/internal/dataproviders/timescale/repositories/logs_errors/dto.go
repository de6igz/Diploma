package logs_errors

import (
	"encoding/json"
	"time"
)

type Event struct {
	Id          string          // uuid события
	UserId      int64           // id пользователя
	ServiceName string          // имя сервиса
	Timestamp   time.Time       // время события
	Log         json.RawMessage // лог в формате JSON
	EventType   *string         // тип события (может быть nil)
	UsedRules   []UsedRule      // список использованных правил
	Language    string          // язык
	UsedActions []UsedAction    // список использованных действий
	EventCount  int64
	Engine      *string //resource engine / errors engine
}

type UsedRule struct {
	RuleId   string `json:"rule_id"`
	RuleName string `json:"rule_name"`
}

type UsedAction struct {
	Type   string           `json:"type"`
	Params UsedActionParams `json:"params"`
}

type UsedActionParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
