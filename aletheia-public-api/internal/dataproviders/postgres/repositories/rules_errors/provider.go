package rules_errors

import (
	v1 "aletheia-public-api/interfaces/types/v1"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type Request struct {
	RuleIds []string `json:"rule_ids"`
	UserId  int64    `json:"user_id"`
}

type Provider interface {
	GetErrorRulesData(ctx context.Context, req Request) ([]*RuleData, error)
	GetAvailableRulesData(ctx context.Context, req Request) ([]*RuleData, error)
	DeleteRuleById(ctx context.Context, ruleId, ruleType string, userId int64) error
	CreateRule(ctx context.Context, userId int64, request v1.CreateRuleRequest) error
	UpdateRuleById(ctx context.Context, userId int64, request v1.UpdateRuleRequest) error
	GetRuleById(ctx context.Context, ruleId string, userId int64) (*v1.RuleDetailResponse, error)
}

type postgresProvider struct {
	conn *sql.DB
}

func NewProvider(conn *sql.DB) Provider {
	return &postgresProvider{conn: conn}
}

func (p *postgresProvider) GetErrorRulesData(ctx context.Context, req Request) ([]*RuleData, error) {

	query := `
SELECT id::text, name,  description
FROM rule_engine.error_rules
WHERE user_id = $1;
`
	rows, err := p.conn.QueryContext(ctx, query, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*RuleData
	for rows.Next() {
		var rd RuleData
		if err := rows.Scan(&rd.Id, &rd.RuleName, &rd.Description); err != nil {
			return nil, err
		}
		rd.RuleType = "errors"
		results = append(results, &rd)
	}
	return results, nil
}

func (p *postgresProvider) GetAvailableRulesData(ctx context.Context, req Request) ([]*RuleData, error) {
	query := `
SELECT id::text, name,  description
FROM rule_engine.error_rules
WHERE user_id = $1 AND service_id IS NULL;
`
	rows, err := p.conn.QueryContext(ctx, query, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*RuleData
	for rows.Next() {
		var rd RuleData
		if err := rows.Scan(&rd.Id, &rd.RuleName, &rd.Description); err != nil {
			return nil, err
		}
		rd.RuleType = "errors"
		results = append(results, &rd)
	}
	return results, nil
}

func (p *postgresProvider) DeleteRuleById(ctx context.Context, ruleId, ruleType string, userId int64) error {
	var query string
	if ruleType == "" {
		return fmt.Errorf("invalid rule type")
	}
	if ruleType == "errors" {
		query = `DELETE FROM rule_engine.error_rules WHERE id = $1 and user_id = $2;`
	} else {
		return fmt.Errorf("invalid rule type")
	}

	_, err := p.conn.ExecContext(ctx, query, ruleId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgresProvider) CreateRule(ctx context.Context, userId int64, request v1.CreateRuleRequest) error {
	// Сериализуем actions в JSON.
	actionsJSON, err := json.Marshal(request.Actions)
	if err != nil {
		return fmt.Errorf("failed to marshal actions: %w", err)
	}

	// Сериализуем root_node в JSON.
	rootNodeJSON, err := json.Marshal(request.RootNode)
	if err != nil {
		return fmt.Errorf("failed to marshal root_node: %w", err)
	}

	// Выполняем INSERT. Поскольку free‑правило не привязано ни к какому сервису, передаем NULL для service_id.
	query := `
		INSERT INTO rule_engine.error_rules (name, actions, root_node, user_id, service_id, description)
		VALUES ($1, $2, $3, $4, NULL, $5);
	`

	_, err = p.conn.ExecContext(ctx, query, request.RuleName, actionsJSON, rootNodeJSON, userId, request.RuleDescription)
	if err != nil {
		return fmt.Errorf("failed to create error rule: %w", err)
	}

	return nil
}

func (p *postgresProvider) UpdateRuleById(ctx context.Context, userId int64, request v1.UpdateRuleRequest) error {
	// Сериализуем actions в JSON.
	actionsJSON, err := json.Marshal(request.Actions)
	if err != nil {
		return fmt.Errorf("failed to marshal actions: %w", err)
	}

	// Сериализуем root_node в JSON.
	rootNodeJSON, err := json.Marshal(request.RootNode)
	if err != nil {
		return fmt.Errorf("failed to marshal root_node: %w", err)
	}

	// Выполняем INSERT. Поскольку free‑правило не привязано ни к какому сервису, передаем NULL для service_id.
	//query := `
	//	INSERT INTO rule_engine.error_rules (name, actions, root_node, user_id, service_id, description)
	//	VALUES ($1, $2, $3, $4, NULL, $5);
	//`
	query := `
	update rule_engine.error_rules
	set name = $1, actions = $2, root_node = $3, description = $4
	where id = $5 and user_id = $6;
	`

	_, err = p.conn.ExecContext(ctx, query, request.RuleName, actionsJSON, rootNodeJSON, request.RuleDescription, request.RuleId, userId)
	if err != nil {
		return fmt.Errorf("failed to create error rule: %w", err)
	}

	return nil
}

func (p *postgresProvider) GetRuleById(ctx context.Context, ruleId string, userId int64) (*v1.RuleDetailResponse, error) {
	query := `
		SELECT name, description, actions, root_node 
		FROM rule_engine.error_rules 
		WHERE id = $1 AND user_id = $2;
	`
	rows, err := p.conn.QueryContext(ctx, query, ruleId, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query rule: %w", err)
	}
	defer rows.Close()

	// Если строки не найдены, возвращаем nil без ошибки.
	if !rows.Next() {
		return nil, nil
	}

	var res v1.RuleDetailResponse
	var actionsJSON, rootNodeJSON []byte
	if err := rows.Scan(&res.Name, &res.Description, &actionsJSON, &rootNodeJSON); err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	var actions []v1.Action
	// Если поле JSON пустое, возвращаем пустой слайс.
	if len(actionsJSON) == 0 {
		actions = []v1.Action{}
	} else {
		if err := json.Unmarshal(actionsJSON, &actions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal actions: %w", err)
		}
	}

	var rootNode v1.Node
	// Если поле JSON пустое, возвращаем пустой объект.
	if len(rootNodeJSON) == 0 {
		rootNode = v1.Node{}
	} else {
		if err := json.Unmarshal(rootNodeJSON, &rootNode); err != nil {
			return nil, fmt.Errorf("failed to unmarshal root_node: %w", err)
		}
	}

	res.Actions = actions
	res.RootNode = rootNode

	return &res, nil
}
