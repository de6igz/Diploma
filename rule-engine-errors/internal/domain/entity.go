package domain

// Event – входящее сообщение, которое нужно проверить правилами.
// Можно добавить много полей (memoryUsage, goroutineCount, timestamp, etc.)
type Event struct {
	ProjectId    string                 `json:"project_id"`
	UserID       string                 `json:"user_id"` // <-- новое поле
	ServiceName  string                 `json:"service_name"`
	Environment  string                 `json:"environment"`
	ErrorMessage string                 `json:"error_message"`
	Version      string                 `json:"version"`
	GoVersion    string                 `json:"go_version"`
	Os           string                 `json:"os"`
	Arch         string                 `json:"arch"`
	EventType    string                 `json:"event_type"`
	Level        string                 `json:"level"`
	EventMessage string                 `json:"event_message"`
	StackTrace   string                 `json:"stack_trace"`
	Tags         []string               `json:"tags"`
	Timestamp    string                 `json:"timestamp"`
	ContextJson  string                 `json:"context_json"`
	Language     string                 `json:"language"`
	Fields       map[string]interface{} `json:"fields"`

	// Поле для повторов (сколько таких ошибок за период)
	// Если мы хотим хранить здесь, а не рассчитывать "на лету".
	RepeatCount int `json:"repeat_count"`
}

// ConditionOperator – тип оператора в правилах.
type ConditionOperator string

const (
	OpEQ         ConditionOperator = "eq"
	OpNEQ        ConditionOperator = "neq"
	OpGT         ConditionOperator = "gt"
	OpGTE        ConditionOperator = "gte"
	OpLT         ConditionOperator = "lt"
	OpLTE        ConditionOperator = "lte"
	OpIN         ConditionOperator = "in"
	OpNIN        ConditionOperator = "nin"
	OpCont       ConditionOperator = "contains"
	OpRepeatOver ConditionOperator = "repeat_over" // нужный нам оператор
)

// Condition – условие
type Condition struct {
	Field    string            `bson:"field"    json:"field"`
	Operator ConditionOperator `bson:"operator" json:"operator"`
	Value    interface{}       `bson:"value"    json:"value"`
}

// LogicNode – узел логического дерева.
// Он может содержать несколько Conditions (одним списком),
// и может содержать дочерние узлы (Children).
// Поле Operator определяет, как связаны *все* эти элементы: AND или OR.
type LogicNode struct {
	Operator   string      `json:"operator"   bson:"operator"`   // "AND" / "OR"
	Conditions []Condition `json:"conditions" bson:"conditions"` // условия на этом уровне
	Children   []LogicNode `json:"children"   bson:"children"`   // подузлы
}

// ActionType – какое действие
type ActionType string

const (
	ActionMail     ActionType = "EMAIL"
	ActionTelegram ActionType = "TELEGRAM"
	ActionDiscord  ActionType = "DISCORD"
	ActionNone     ActionType = "NONE"
)

// Action – одно действие, могут быть параметры.
type Action struct {
	Type   ActionType        `bson:"type"   json:"type"`
	Params map[string]string `bson:"params" json:"params"`
}

// Rule – список условий и действий
type Rule struct {
	ID string `bson:"_id,omitempty" json:"id"`

	UserID      int         `bson:"user_id"        json:"user_id"`
	ServiceName string      `bson:"service_name"   json:"service_name"`
	Name        string      `bson:"name"          json:"name"`
	Conditions  []Condition `bson:"conditions"    json:"conditions"`
	Actions     []Action    `bson:"actions"       json:"actions"`

	// RootNode – корень "дерева" логики (AND/OR + conditions + children)
	RootNode LogicNode `bson:"root_node"     json:"root_node"`
}
