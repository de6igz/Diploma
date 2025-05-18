package projects

import (
	types "aletheia-public-api/interfaces/types/v1"
	projectsRepo "aletheia-public-api/internal/dataproviders/postgres/repositories/projects"
	"context"
)

// ProjectsUsecase описывает методы для получения проектов.
type ProjectsUsecase interface {
	GetProjects(ctx context.Context, userId int64) ([]*projectsRepo.Project, error)
	GetProjectByID(ctx context.Context, userId int64, projectID string) (*projectsRepo.Project, error)
	CreateProject(ctx context.Context, userId int64, project *types.CreateProjectRequest) error
	DeleteProjectById(ctx context.Context, userId int64, projectID string) error
	UpdateProject(ctx context.Context, project *types.UpdateProjectRequest, projectId string, userId int64) error
}

type projectsUsecase struct {
	projectsRepo projectsRepo.Provider
}

// NewProjectsUsecase создаёт usecase с инъекцией репозитория.
func NewProjectsUsecase(provider projectsRepo.Provider) ProjectsUsecase {
	return &projectsUsecase{
		projectsRepo: provider,
	}
}

// GetProjects получает проекты для заданного пользователя.
func (uc *projectsUsecase) GetProjects(ctx context.Context, userId int64) ([]*projectsRepo.Project, error) {
	req := projectsRepo.Request{UserId: userId}
	return uc.projectsRepo.GetProjects(ctx, req)
}

// GetProjectByID получает проект по ID для заданного пользователя.
func (uc *projectsUsecase) GetProjectByID(ctx context.Context, userId int64, projectID string) (*projectsRepo.Project, error) {
	req := projectsRepo.Request{UserId: userId, ProjectId: projectID}
	return uc.projectsRepo.GetProjectById(ctx, req)
}

func (uc *projectsUsecase) CreateProject(ctx context.Context, userId int64, project *types.CreateProjectRequest) error {

	servs := make([]projectsRepo.ServiceCreateRequest, len(project.Services))
	for i, service := range project.Services {
		servs[i] = projectsRepo.ServiceCreateRequest{
			ServiceName:   service.ServiceName,
			ErrorRules:    service.ErrorRules,
			ResourceRules: service.ResourceRules,
		}
	}
	req := projectsRepo.CreateProjectRequest{
		ProjectName: project.ProjectName,
		Description: project.Description,
		Services:    servs,
	}
	err := uc.projectsRepo.CreateProject(ctx, req, userId)

	return err
}

func (uc *projectsUsecase) DeleteProjectById(ctx context.Context, userId int64, projectID string) error {
	req := projectsRepo.Request{UserId: userId, ProjectId: projectID}
	return uc.projectsRepo.DeleteProjectById(ctx, req)
}

func (uc *projectsUsecase) UpdateProject(ctx context.Context, project *types.UpdateProjectRequest, projectID string, userId int64) error {
	servs := make([]projectsRepo.ServiceCreateRequest, len(project.Services))
	for i, service := range project.Services {
		servs[i] = projectsRepo.ServiceCreateRequest{
			ServiceName:   service.ServiceName,
			ErrorRules:    service.ErrorRules,
			ResourceRules: service.ResourceRules,
		}
	}
	req := projectsRepo.UpdateProjectRequest{
		ProjectName: project.ProjectName,
		Description: project.Description,
		Services:    servs,
		ProjectId:   projectID,
	}

	return uc.projectsRepo.UpdateProject(ctx, req, userId)
}
