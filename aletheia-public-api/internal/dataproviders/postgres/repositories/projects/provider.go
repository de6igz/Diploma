package projects

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

type Request struct {
	UserId    int64  `json:"user_id"`
	ProjectId string `json:"project_id"` // используется для GetProjectById
}

type Provider interface {
	GetProjects(ctx context.Context, req Request) ([]*Project, error)
	GetProjectById(ctx context.Context, req Request) (*Project, error)
	CreateProject(ctx context.Context, req CreateProjectRequest, userId int64) error
	DeleteProjectById(ctx context.Context, req Request) error
	UpdateProject(ctx context.Context, req UpdateProjectRequest, userId int64) error
}

type postgresProvider struct {
	conn *sql.DB
}

func NewProvider(conn *sql.DB) Provider {
	return &postgresProvider{conn: conn}
}

func (p *postgresProvider) GetProjects(ctx context.Context, req Request) ([]*Project, error) {
	// Запрос агрегирует error- и resource-правила через таблицу services.
	query := `
SELECT p.id::text,
       p.project_name,
       p.description
from rule_engine.projects p
WHERE p.user_id = $1;
`
	rows, err := p.conn.QueryContext(ctx, query, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		var proj Project
		if err = rows.Scan(&proj.ID, &proj.ProjectName, &proj.Description); err != nil {
			return nil, err
		}

		projects = append(projects, &proj)
	}
	return projects, nil
}

func (p *postgresProvider) GetProjectById(ctx context.Context, req Request) (*Project, error) {
	// Этот запрос выбирает основные поля проекта и агрегирует связанные правила.
	// Для error_rules и resource_rules используется json_build_object, чтобы сформировать объект с rule_id и rule_name.
	query := `
SELECT p.id::text,
       p.project_name,
       p.description,
       p.user_id,
       COALESCE(
         (SELECT json_agg(json_build_object('rule_id', er.id::text, 'rule_name', er.name))
          FROM rule_engine.error_rules er
          JOIN rule_engine.services s ON er.service_id = s.id
          WHERE s.project_id = p.id), '[]'
       ) AS connected_error_rules,
       COALESCE(
         (SELECT json_agg(json_build_object('rule_id', rr.id::text, 'rule_name', rr.name))
          FROM rule_engine.resource_rules rr
          JOIN rule_engine.services s ON rr.service_id = s.id
          WHERE s.project_id = p.id), '[]'
       ) AS connected_resource_rules
FROM rule_engine.projects p
WHERE p.id::text = $1;
`
	row := p.conn.QueryRowContext(ctx, query, req.ProjectId)

	var proj Project
	var errRulesJSON, resRulesJSON []byte

	if err := row.Scan(&proj.ID, &proj.ProjectName, &proj.Description, &proj.UserId, &errRulesJSON, &resRulesJSON); err != nil {
		return nil, fmt.Errorf("failed to scan project row: %w", err)
	}

	if err := json.Unmarshal(errRulesJSON, &proj.ConnectedErrorRules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal connected_error_rules: %w", err)
	}

	if err := json.Unmarshal(resRulesJSON, &proj.ConnectedResourceRules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal connected_resource_rules: %w", err)
	}

	return &proj, nil
}

func (p *postgresProvider) CreateProject(ctx context.Context, req CreateProjectRequest, userId int64) error {
	// Начинаем транзакцию.
	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Если где-то возникнет ошибка – откатываем транзакцию.
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Создаём проект и получаем его id.
	var projectId int
	projectInsertQuery := `
		INSERT INTO rule_engine.projects (project_name, description, user_id)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	err = tx.QueryRowContext(ctx, projectInsertQuery, req.ProjectName, req.Description, userId).Scan(&projectId)
	if err != nil {
		return fmt.Errorf("failed to insert project: %w", err)
	}

	// Для каждого сервиса в запросе создаём запись в таблице services,
	// затем обновляем правила (error и resource) – привязываем их к новому сервису.
	for _, svc := range req.Services {
		var serviceId int
		serviceInsertQuery := `
			INSERT INTO rule_engine.services (project_id, service_name)
			VALUES ($1, $2)
			RETURNING id;
		`
		err = tx.QueryRowContext(ctx, serviceInsertQuery, projectId, svc.ServiceName).Scan(&serviceId)
		if err != nil {
			return fmt.Errorf("failed to insert service '%s': %w", svc.ServiceName, err)
		}

		// Обновляем error-правила: устанавливаем для каждого правила service_id = serviceId.
		for _, ruleId := range svc.ErrorRules {
			updateErrorRuleQuery := `
				UPDATE rule_engine.error_rules
				SET service_id = $1
				WHERE id = $2 AND user_id = $3;
			`
			_, err = tx.ExecContext(ctx, updateErrorRuleQuery, serviceId, ruleId, userId)
			if err != nil {
				return fmt.Errorf("failed to update error rule id %d for service '%s': %w", ruleId, svc.ServiceName, err)
			}
		}

		// Аналогично обновляем resource-правила.
		for _, ruleId := range svc.ResourceRules {
			updateResourceRuleQuery := `
				UPDATE rule_engine.resource_rules
				SET service_id = $1
				WHERE id = $2 AND user_id = $3;
			`
			_, err = tx.ExecContext(ctx, updateResourceRuleQuery, serviceId, ruleId, userId)
			if err != nil {
				return fmt.Errorf("failed to update resource rule id %d for service '%s': %w", ruleId, svc.ServiceName, err)
			}
		}
	}

	// Если все операции прошли успешно, фиксируем транзакцию.
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteProjectById удаляет проект для заданного пользователя и все связанные с ним сервисы и правила.
func (p *postgresProvider) DeleteProjectById(ctx context.Context, req Request) error {
	// Преобразуем идентификатор проекта из строки в int.
	projectID, err := strconv.Atoi(req.ProjectId)
	if err != nil {
		return fmt.Errorf("invalid project id '%s': %w", req.ProjectId, err)
	}

	// Начинаем транзакцию.
	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Если где-то возникает ошибка, транзакция будет откатана.
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Удаляем error-правила, связанные с сервисами этого проекта.
	deleteErrorRulesQuery := `
		DELETE FROM rule_engine.error_rules
		WHERE service_id IN (
			SELECT id FROM rule_engine.services WHERE project_id = $1
		);
	`
	_, err = tx.ExecContext(ctx, deleteErrorRulesQuery, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete error rules for project %d: %w", projectID, err)
	}

	// 2. Удаляем resource-правила, связанные с сервисами этого проекта.
	deleteResourceRulesQuery := `
		DELETE FROM rule_engine.resource_rules
		WHERE service_id IN (
			SELECT id FROM rule_engine.services WHERE project_id = $1
		);
	`
	_, err = tx.ExecContext(ctx, deleteResourceRulesQuery, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete resource rules for project %d: %w", projectID, err)
	}

	// 3. Удаляем сервисы, принадлежащие проекту.
	deleteServicesQuery := `
		DELETE FROM rule_engine.services
		WHERE project_id = $1;
	`
	_, err = tx.ExecContext(ctx, deleteServicesQuery, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete services for project %d: %w", projectID, err)
	}

	// 4. Удаляем сам проект, убеждаясь, что он принадлежит указанному пользователю.
	deleteProjectQuery := `
		DELETE FROM rule_engine.projects
		WHERE id = $1 AND user_id = $2;
	`
	result, err := tx.ExecContext(ctx, deleteProjectQuery, projectID, req.UserId)
	if err != nil {
		return fmt.Errorf("failed to delete project %d: %w", projectID, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("no project found with id %d for user %d", projectID, req.UserId)
	}

	// Фиксируем транзакцию.
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateProject выполняет обновление проекта, его сервисов и привязанных правил в рамках одной транзакции.
func (p *postgresProvider) UpdateProject(ctx context.Context, req UpdateProjectRequest, userId int64) error {
	// Начинаем транзакцию.
	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Откат транзакции при ошибке.
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Преобразуем идентификатор проекта из строки в int.
	projectID, err := strconv.Atoi(req.ProjectId)
	if err != nil {
		return fmt.Errorf("invalid project id '%s': %w", req.ProjectId, err)
	}

	// Обновляем поля проекта (название, описание) для заданного проекта и пользователя.
	updateProjectQuery := `
		UPDATE rule_engine.projects
		SET project_name = $1, description = $2
		WHERE id = $3 AND user_id = $4;
	`
	_, err = tx.ExecContext(ctx, updateProjectQuery, req.ProjectName, req.Description, projectID, userId)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	// Освобождаем связанные правила: для всех сервисов данного проекта устанавливаем service_id = NULL.
	freeErrorRulesQuery := `
		UPDATE rule_engine.error_rules
		SET service_id = NULL
		WHERE service_id IN (SELECT id FROM rule_engine.services WHERE project_id = $1);
	`
	_, err = tx.ExecContext(ctx, freeErrorRulesQuery, projectID)
	if err != nil {
		return fmt.Errorf("failed to free error rules: %w", err)
	}
	freeResourceRulesQuery := `
		UPDATE rule_engine.resource_rules
		SET service_id = NULL
		WHERE service_id IN (SELECT id FROM rule_engine.services WHERE project_id = $1);
	`
	_, err = tx.ExecContext(ctx, freeResourceRulesQuery, projectID)
	if err != nil {
		return fmt.Errorf("failed to free resource rules: %w", err)
	}

	// Удаляем все существующие сервисы для данного проекта.
	deleteServicesQuery := `
		DELETE FROM rule_engine.services
		WHERE project_id = $1;
	`
	_, err = tx.ExecContext(ctx, deleteServicesQuery, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete services: %w", err)
	}

	// Для каждого сервиса в запросе создаём новый сервис и обновляем привязки правил.
	for _, svc := range req.Services {
		var newServiceId int
		serviceInsertQuery := `
			INSERT INTO rule_engine.services (project_id, service_name)
			VALUES ($1, $2)
			RETURNING id;
		`
		err = tx.QueryRowContext(ctx, serviceInsertQuery, projectID, svc.ServiceName).Scan(&newServiceId)
		if err != nil {
			return fmt.Errorf("failed to insert service '%s': %w", svc.ServiceName, err)
		}

		// Обновляем error-правила: привязываем правило к новому сервису.
		for _, ruleId := range svc.ErrorRules {
			updateErrorRuleQuery := `
				UPDATE rule_engine.error_rules
				SET service_id = $1
				WHERE id = $2 AND user_id = $3;
			`
			_, err = tx.ExecContext(ctx, updateErrorRuleQuery, newServiceId, ruleId, userId)
			if err != nil {
				return fmt.Errorf("failed to update error rule id %d for service '%s': %w", ruleId, svc.ServiceName, err)
			}
		}

		// Обновляем resource-правила: привязываем правило к новому сервису.
		for _, ruleId := range svc.ResourceRules {
			updateResourceRuleQuery := `
				UPDATE rule_engine.resource_rules
				SET service_id = $1
				WHERE id = $2 AND user_id = $3;
			`
			_, err = tx.ExecContext(ctx, updateResourceRuleQuery, newServiceId, ruleId, userId)
			if err != nil {
				return fmt.Errorf("failed to update resource rule id %d for service '%s': %w", ruleId, svc.ServiceName, err)
			}
		}
	}

	// Фиксируем транзакцию.
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
