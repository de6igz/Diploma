package resourceEvent

// ResourceEvent – доменная модель «ресурсного» события.
type ResourceEvent struct {
	ProjectID    string                 `json:"project_id"`
	ServiceName  string                 `json:"service_name"`
	Environment  string                 `json:"environment"`
	Version      string                 `json:"version"`
	GoVersion    string                 `json:"go_version"`
	Os           string                 `json:"os"`
	Arch         string                 `json:"arch"`
	EventType    string                 `json:"event_type"` // Обычно будет "RESOURCE_USAGE"
	EventMessage string                 `json:"event_message"`
	StackTrace   string                 `json:"stack_trace"`
	Tags         []string               `json:"tags"`
	Timestamp    string                 `json:"timestamp"`
	ContextJson  string                 `json:"context_json"`
	Fields       map[string]interface{} `json:"fields"`
	Language     string                 `json:"language"`
}

// ResourceResponse – ответ клиенту.
type ResourceResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message"`
}
