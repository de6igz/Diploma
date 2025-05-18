package messageEvent

// MessageEvent – доменная модель события типа "MESSAGE".
type MessageEvent struct {
	ProjectID    string   `json:"project_id"`
	ServiceName  string   `json:"service_name"`
	Environment  string   `json:"environment"`
	Version      string   `json:"version"`
	GoVersion    string   `json:"go_version"`
	Os           string   `json:"os"`
	Arch         string   `json:"arch"`
	EventType    string   `json:"event_type"` // Обычно "MESSAGE"
	EventMessage string   `json:"event_message"`
	StackTrace   string   `json:"stack_trace"`
	Tags         []string `json:"tags"`
	Timestamp    string   `json:"timestamp"`
	ContextJson  string   `json:"context_json"`
	Language     string   `json:"language"`
}

// MessageResponse – структура для ответа клиенту при отправке MessageEvent.
type MessageResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message"`
}
