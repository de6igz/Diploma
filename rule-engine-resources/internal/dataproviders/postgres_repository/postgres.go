package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"rule-engine-resources/internal/config"
	"strconv"

	"github.com/rs/zerolog"
	"rule-engine-resources/internal/domain"
	"rule-engine-resources/internal/usecases"
)

type PostgresRuleRepository struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

// NewPostgresRuleRepository создаёт и инициализирует подключение к PostgreSQL
// и возвращает репозиторий, реализующий usecases.RuleRepository.
func NewPostgresRuleRepository(logger *zerolog.Logger, cfg *config.Config) (usecases.RuleRepository, error) {
	user := cfg.Postgres.User
	pass := cfg.Postgres.Password
	host := cfg.Postgres.Host
	dbName := cfg.Postgres.DBName
	port := cfg.Postgres.Port
	if port == 0 {
		port = 5432
	}

	// Формируем DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, pass, host, port, dbName)

	// Подключаемся через sqlx
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to Postgres")
		return nil, err
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		logger.Error().Err(err).Msg("Failed to ping Postgres")
		return nil, err
	}

	logger.Info().Msgf("Connected to Postgres successfully %v", dsn)
	return &PostgresRuleRepository{
		db:     db,
		logger: logger,
	}, nil
}

// GetRulesByUserAndServiceAndProjectId возвращает error-правила для заданного userID, serviceName и projectId.
func (pr *PostgresRuleRepository) GetRulesByUserAndServiceAndProjectId(
	ctx context.Context,
	userID string,
	serviceName string,
	projectId string,
) ([]domain.Rule, error) {

	pr.logger.Debug().Msgf("PostgresRuleRepository.GetRulesByUserAndServiceAndProjectId: user=%s, service=%s, project=%s", userID, serviceName, projectId)

	// Преобразуем идентификаторы в целочисленный формат
	userIdInt, err := strconv.Atoi(userID)
	if err != nil {
		pr.logger.Error().Err(err).Msg("Failed to parse userID")
		return nil, err
	}

	projIdInt, err := strconv.Atoi(projectId)
	if err != nil {
		pr.logger.Error().Err(err).Msg("Failed to parse projectId")
		return nil, err
	}

	// Запрос для получения error-правил из таблицы error_rules, с join по services для фильтрации по service_name и project_id
	query := `
		SELECT r.id, r.name, r.actions, r.root_node, r.user_id, s.service_name
		FROM rule_engine.resource_rules r
		JOIN rule_engine.services s ON r.service_id = s.id
		WHERE r.user_id = $1 AND s.service_name = $2 AND s.project_id = $3;
	`
	rows, err := pr.db.QueryContext(ctx, query, userIdInt, serviceName, projIdInt)
	if err != nil {
		pr.logger.Error().Err(err).Msg("Failed to fetch error rules")
		return nil, err
	}
	defer rows.Close()

	var rules []domain.Rule
	for rows.Next() {
		var (
			id                int
			name              string
			actionsRaw        []byte
			rootNodeRaw       []byte
			userIdFromDB      int
			serviceNameFromDB string
		)
		if err := rows.Scan(&id, &name, &actionsRaw, &rootNodeRaw, &userIdFromDB, &serviceNameFromDB); err != nil {
			pr.logger.Warn().Err(err).Msg("Failed to scan rule row")
			continue
		}

		// Декодируем JSON-поля
		var actions []domain.Action
		if err := json.Unmarshal(actionsRaw, &actions); err != nil {
			pr.logger.Warn().Err(err).Msg("Failed to unmarshal actions JSON")
			continue
		}

		var rootNode domain.LogicNode
		if err := json.Unmarshal(rootNodeRaw, &rootNode); err != nil {
			pr.logger.Warn().Err(err).Msg("Failed to unmarshal root_node JSON")
			continue
		}

		rule := domain.Rule{
			ID:          strconv.Itoa(id),
			UserID:      userIdFromDB,
			ServiceName: serviceNameFromDB,
			Name:        name,
			Actions:     actions,
			RootNode:    rootNode,
			Conditions:  rootNode.Conditions, // при необходимости
		}
		rules = append(rules, rule)
	}

	pr.logger.Debug().Msgf("Fetched %d error rules", len(rules))
	return rules, nil
}

var _ usecases.RuleRepository = (*PostgresRuleRepository)(nil)
