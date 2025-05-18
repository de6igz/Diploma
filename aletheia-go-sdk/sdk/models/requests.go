package models

type EventRequest struct {
	ProjectId    string   `json:"project_id"`
	ServiceName  string   `json:"service_name"`
	Environment  string   `json:"environment"`
	Version      string   `json:"version"`
	GoVersion    string   `json:"go_version"`
	Os           string   `json:"os"`
	Arch         string   `json:"arch"`
	EventType    string   `json:"event_type"`
	EventMessage string   `json:"event_message"`
	StackTrace   string   `json:"stack_trace"`
	Tags         []string `json:"tags"`
	Timestamp    string   `json:"timestamp"`
	Language     string   `json:"language"`
	ContextJson  string   `json:"context_json"`
}

type ResourceRequest struct {
	ProjectId    string                 `json:"project_id"`
	ServiceName  string                 `json:"service_name"`
	Environment  string                 `json:"environment"`
	Version      string                 `json:"version"`
	GoVersion    string                 `json:"go_version"`
	Os           string                 `json:"os"`
	Arch         string                 `json:"arch"`
	EventType    string                 `json:"event_type"`
	EventMessage string                 `json:"event_message"`
	Tags         []string               `json:"tags"`
	Timestamp    string                 `json:"timestamp"`
	Language     string                 `json:"language"`
	Fields       map[string]interface{} `json:"fields"`
}

type ErrorEvent struct {
	ProjectId    string                 `json:"project_id"`
	ServiceName  string                 `json:"service_name"`
	Environment  string                 `json:"environment"`
	Version      string                 `json:"version"`
	ErrorMessage string                 `json:"error_message"`
	GoVersion    string                 `json:"go_version"`
	Os           string                 `json:"os"`
	Arch         string                 `json:"arch"`
	EventType    string                 `json:"event_type"` // Обычно "ERROR"
	EventMessage string                 `json:"event_message"`
	StackTrace   string                 `json:"stack_trace"`
	Tags         []string               `json:"tags"`
	Timestamp    string                 `json:"timestamp"`
	Language     string                 `json:"language"`
	ContextJson  string                 `json:"context_json"`
	Fields       map[string]interface{} `json:"fields"`
}
