package interfaces

import (
	v1 "aletheia-public-api/interfaces/types/v1"
	"context"
)

// Projects
// @tg http-server log metrics
// @tg http-prefix=v1
type Projects interface {
	// GetProjects
	// @tg summary=`Получить проекты`
	// @tg desc=`Возвращает список проектов`
	// @tg http-method=GET
	// @tg http-path=/projects
	// @tg http-headers=userId|X-User-Id
	GetProjects(ctx context.Context, userId int64) (items v1.ProjectsResponse, err error)

	// GetProjectByID
	// @tg summary=`Получить детали проекта`
	// @tg desc=`Возвращает подробную информацию о проекте по его ID`
	// @tg http-method=GET
	// @tg http-path=/project/:projectID
	// @tg http-headers=userId|X-User-Id
	GetProjectByID(ctx context.Context, projectID string, userId int64) (project v1.ProjectDetailResponse, err error)

	// DeleteProjectByID
	// @tg summary=`Удалить проект по Id`
	// @tg desc=`Удалить проект по Id`
	// @tg http-method=DELETE
	// @tg http-path=/project/:projectID
	// @tg http-headers=userId|X-User-Id
	DeleteProjectByID(ctx context.Context, projectID string, userId int64) (status bool, err error)

	// CreateProject
	// @tg summary=`Создать проект`
	// @tg desc=`Создать новый проект`
	// @tg http-method=POST
	// @tg http-path=/project/create
	// @tg http-headers=userId|X-User-Id
	CreateProject(ctx context.Context, project *v1.CreateProjectRequest, userId int64) (status bool, err error)

	// UpdateProject
	// @tg summary=` Обновить проект`
	// @tg desc=`Обновить проект`
	// @tg http-method=PUT
	// @tg http-path=/project/:projectID
	// @tg http-headers=userId|X-User-Id
	UpdateProject(ctx context.Context, project *v1.UpdateProjectRequest, projectID string, userId int64) (status bool, err error)
}
