package errorEvent

type ErrorEvent struct {
	ProjectID    string                 `json:"project_id"`
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

type ErrorResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message"`
}
