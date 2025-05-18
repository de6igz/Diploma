package projects

import (
	v1 "aletheia-public-api/interfaces/types/v1"
	projectsRepo "aletheia-public-api/internal/dataproviders/postgres/repositories/projects"
)

// ProjectSerializer определяет методы преобразования внутренней модели проекта в внешний тип.
type ProjectSerializer interface {
	SerializeProjects([]*projectsRepo.Project) []v1.Project
	SerializeProjectDetail(*projectsRepo.Project, []*v1.Event) v1.ProjectDetail
}

type defaultProjectSerializer struct{}

// NewProjectSerializer возвращает экземпляр стандартного сериализатора проектов.
func NewProjectSerializer() ProjectSerializer {
	return &defaultProjectSerializer{}
}

// SerializeProjects преобразует срез внутренних проектов в []v1.Project.
func (s *defaultProjectSerializer) SerializeProjects(projects []*projectsRepo.Project) []v1.Project {
	if len(projects) == 0 {
		return []v1.Project{}
	}
	result := make([]v1.Project, 0, len(projects))
	for _, project := range projects {
		result = append(result, v1.Project{
			ID:          project.ID, // Теперь Id уже строка
			ProjectName: project.ProjectName,
			Description: project.Description,
		})
	}
	return result
}

// SerializeProjectDetail преобразует внутреннюю модель проекта в v1.ProjectDetail.
func (s *defaultProjectSerializer) SerializeProjectDetail(project *projectsRepo.Project, events []*v1.Event) v1.ProjectDetail {
	// Преобразуем агрегированные error-правила (срез []*RuleData) в []v1.RuleData.
	var errorRules []v1.RuleData
	for _, rulePtr := range project.ConnectedErrorRules {
		if rulePtr != nil {
			errorRules = append(errorRules, v1.RuleData{
				RuleId:   rulePtr.RuleId,
				RuleName: rulePtr.RuleName,
			})
		}
	}
	// Преобразуем агрегированные resource-правила.
	var resourceRules []v1.RuleData
	for _, rulePtr := range project.ConnectedResourceRules {
		if rulePtr != nil {
			resourceRules = append(resourceRules, v1.RuleData{
				RuleId:   rulePtr.RuleId,
				RuleName: rulePtr.RuleName,
			})
		}
	}

	// Формируем сервис. Если у проекта несколько сервисов, логику группировки можно расширить.
	service := v1.Service{
		ServiceName:   "Example Service", // Здесь можно использовать фактическое имя сервиса, если оно хранится отдельно.
		ErrorRules:    errorRules,
		ResourceRules: resourceRules,
		Events:        events,
	}

	return v1.ProjectDetail{
		Id:          project.ID,
		ProjectName: project.ProjectName,
		Services:    []v1.Service{service},
	}
}
