package projects

import (
	v1 "aletheia-public-api/interfaces/types/v1"
	"aletheia-public-api/internal/api/v1/events"
	"aletheia-public-api/internal/dataproviders/postgres"
	projectsRepo "aletheia-public-api/internal/dataproviders/postgres/repositories/projects"
	"aletheia-public-api/internal/dataproviders/timescale"
	"aletheia-public-api/internal/dataproviders/timescale/repositories/logs_errors"
	"context"
	"fmt"
)

type Projects struct {
	projectUsecase   ProjectsUsecase
	eventsUsecase    events.EventsUsecase
	serializer       ProjectSerializer
	eventsSerializer events.Serializer
}

// NewProjects создаёт новый обработчик, инициализируя usecase и сериализаторы.
// Обратите внимание: теперь используется postgres.GlobalInstance (*sql.DB) вместо mongo-клиента.
func NewProjects() *Projects {
	pgConn := postgres.GlobalInstance
	provider := projectsRepo.NewProvider(pgConn)
	usecase := NewProjectsUsecase(provider)
	serializer := NewProjectSerializer()
	logsErrorRepo := logs_errors.NewProvider(timescale.GlobalInstance)
	eventsUsecase := events.NewEventsUsecase(logsErrorRepo)
	eventsSerializer := events.NewSerializer()

	return &Projects{
		projectUsecase:   usecase,
		serializer:       serializer,
		eventsSerializer: eventsSerializer,
		eventsUsecase:    eventsUsecase,
	}
}

// GetProjects получает проекты для пользователя через usecase, затем сериализует их.
func (p *Projects) GetProjects(ctx context.Context, userId int64) (v1.ProjectsResponse, error) {
	projects, err := p.projectUsecase.GetProjects(ctx, userId)

	if err != nil {
		return v1.ProjectsResponse{}, err
	}

	if projects == nil {
		return v1.ProjectsResponse{}, nil
	}
	return v1.ProjectsResponse{
		Projects: p.serializer.SerializeProjects(projects),
	}, nil
}

// GetProjectByID получает детальную информацию по проекту, а также связанные события.
func (p *Projects) GetProjectByID(ctx context.Context, projectID string, userId int64) (v1.ProjectDetailResponse, error) {
	project, err := p.projectUsecase.GetProjectByID(ctx, userId, projectID)
	if err != nil {
		return v1.ProjectDetailResponse{}, err
	}
	if project == nil {
		return v1.ProjectDetailResponse{}, nil
	}

	eventsData, err := p.eventsUsecase.GetEventsByProjectId(ctx, 24, userId, projectID)
	if err != nil {
		return v1.ProjectDetailResponse{}, err
	}
	serializedEvents := p.eventsSerializer.SerializeEvents(eventsData)

	serialized := p.serializer.SerializeProjectDetail(project, serializedEvents)
	return v1.ProjectDetailResponse{
		Project: &serialized,
	}, nil
}

func (p *Projects) CreateProject(ctx context.Context, project *v1.CreateProjectRequest, userId int64) (status bool, err error) {
	if project == nil {
		return false, fmt.Errorf("project is nil")
	}
	err = p.projectUsecase.CreateProject(ctx, userId, project)

	if err != nil {
		return false, err
	}
	return false, err
}

func (p *Projects) DeleteProjectByID(ctx context.Context, projectID string, userId int64) (status bool, err error) {
	err = p.projectUsecase.DeleteProjectById(ctx, userId, projectID)
	if err != nil {
		return false, err
	}

	return true, nil

}

func (p *Projects) UpdateProject(ctx context.Context, project *v1.UpdateProjectRequest, projectID string, userId int64) (status bool, err error) {
	err = p.projectUsecase.UpdateProject(ctx, project, projectID, userId)
	if err != nil {
		return false, err
	}
	return true, nil
}
