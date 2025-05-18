package v1

import "time"

type EventsResponse struct {
	Events []*Event `json:"events"`
}

type Event struct {
	EventType   *string `json:"eventType,omitempty"`
	ServiceName string  `json:"serviceName"`
	EventsCount int     `json:"eventsCount"`
	Assignee    string  `json:"assignee"`
	Language    string  `json:"language"`
}

type MeResponse struct {
	Username string `json:"username"`
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

type EventDetailResponse struct {
	Event *EventDetail `json:"event,omitempty"`
}

type EventDetail struct {
	ID          string       `json:"id"`
	Log         string       `json:"log"`
	ServiceName string       `json:"serviceName"`
	EventType   *string      `json:"eventType,omitempty"`
	Timestamp   time.Time    `json:"timestamp"`
	UsedRules   []UsedRule   `json:"usedRules"`   // список использованных правил
	UsedActions []UsedAction `json:"usedActions"` // список использованных действий
	Language    string       `json:"language"`
}

type EventsByEventTypeResponse struct {
	Events []EventByType `json:"events"`
}

type EventByType struct {
	// uuid, service_name, timestamp, language, log
	Id          string    `json:"id"`
	ServiceName string    `json:"serviceName"`
	Timestamp   time.Time `json:"timestamp"`
	Language    string    `json:"language"`
	Log         string    `json:"log"`
	EventType   string    `json:"eventType"`
}

type ProjectsResponse struct {
	Projects []Project `json:"projects"`
}

type Project struct {
	ID          string `json:"id"`
	ProjectName string `json:"projectName"`
	Description string `json:"description"`
}

type ProjectDetailResponse struct {
	Project *ProjectDetail `json:"project,omitempty"`
}

type ProjectDetail struct {
	Id          string    `json:"id"`
	ProjectName string    `json:"projectName"`
	Services    []Service `json:"services"`
}

type Service struct {
	ServiceName   string     `json:"serviceName"`
	ErrorRules    []RuleData `json:"errorRules"`
	ResourceRules []RuleData `json:"resourceRules"`
	Events        []*Event   `json:"events,omitempty"`
}

type RuleData struct {
	RuleId   string  `json:"ruleId"`
	RuleName string  `json:"ruleName"`
	RuleType *string `json:"ruleType,omitempty"`
}

type RulesResponse struct {
	Rules []*Rule `json:"rules,omitempty"`
}

type Rule struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	RuleType    *string `json:"ruleType,omitempty"`
	Description *string `json:"description,omitempty"`
}

type RuleDetailResponse struct {
	Name        string   `json:"name"`
	RuleType    *string  `json:"ruleType,omitempty"`
	Description *string  `json:"description,omitempty"`
	RootNode    Node     `json:"root_node"`
	Actions     []Action `json:"actions"`
}

type Node struct {
	Operator   string      `json:"operator"`
	Conditions []Condition `json:"conditions"`
	Children   []Node      `json:"children"`
}

type Condition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type Action struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}

type CreateProjectRequest struct {
	ProjectName string                 `json:"projectName"`
	Description string                 `json:"description"`
	Services    []ServiceCreateRequest `json:"services"`
}

// UpdateProjectRequest описывает входной JSON для обновления проекта.
type UpdateProjectRequest struct {
	ProjectId   string                 `json:"projectId"`
	ProjectName string                 `json:"projectName"`
	Description string                 `json:"description"`
	Services    []ServiceCreateRequest `json:"services"`
}

type ServiceCreateRequest struct {
	ServiceName   string `json:"serviceName"`
	ErrorRules    []int  `json:"errorRules"`
	ResourceRules []int  `json:"resourceRules"`
}

type DeleteRuleRequest struct {
	RuleId   string `json:"ruleId"`
	RuleType string `json:"ruleType"`
}

type CreateRuleRequest struct {
	RuleName        string   `json:"name"`
	RuleDescription string   `json:"description"`
	RuleType        string   `json:"ruleType"`
	RootNode        Node     `json:"root_node"`
	Actions         []Action `json:"actions"`
}

type UpdateRuleRequest struct {
	RuleId          string   `json:"ruleId"`
	RuleName        string   `json:"name"`
	RuleDescription string   `json:"description"`
	RuleType        string   `json:"ruleType"`
	RootNode        Node     `json:"root_node"`
	Actions         []Action `json:"actions"`
}

type RuleByIdRequest struct {
	RuleId   string `json:"ruleId"`
	RuleType string `json:"ruleType"`
}
