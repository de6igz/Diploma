package projects

// Project описывает модель проекта в Postgres.
type Project struct {
	ID                     string      `json:"id"`
	ConnectedErrorRules    []*RuleData `json:"connected_error_rules"`
	ConnectedResourceRules []*RuleData `json:"connected_resource_rules"`
	UserId                 int64       `json:"user_id"`
	Description            string      `json:"description"`
	ProjectName            string      `json:"project_name"`
}

type RuleData struct {
	RuleId   string `json:"rule_id"`
	RuleName string `json:"rule_name"`
	//RuleType *string `json:"rule_type"`
}

// CreateProjectRequest – запрос на создание проекта.
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

// ServiceCreateRequest – запрос на создание сервиса внутри проекта.
type ServiceCreateRequest struct {
	ServiceName   string `json:"serviceName"`
	ErrorRules    []int  `json:"errorRules"`    // идентификаторы error-правил, которые нужно привязать
	ResourceRules []int  `json:"resourceRules"` // идентификаторы resource-правил, которые нужно привязать
}
