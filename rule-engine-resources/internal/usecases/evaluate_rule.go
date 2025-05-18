package usecases

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"rule-engine-resources/internal/dataproviders/redis_repository"
	"rule-engine-resources/internal/dataproviders/timescale_repository"
	"rule-engine-resources/internal/domain"
)

const (
	ENGINE = "resources"
)

type RuleRepository interface {
	GetRulesByUserAndServiceAndProjectId(ctx context.Context, userID string, serviceName, projectId string) ([]domain.Rule, error)
}

type AlertDispatcher interface {
	DispatchActions(ctx context.Context, e *domain.Event, actions []domain.Action) error
}

type EvaluateRulesUseCase struct {
	ruleRepo        RuleRepository
	timeScaleRepo   timescale_repository.TimescaleRepository
	alertDispatcher AlertDispatcher
	redisCounter    *redis_repository.RedisRepeatCounter
	redisCache      *redis_repository.RedisCache
	logger          *zerolog.Logger
}

func NewEvaluateRulesUseCase(
	rr RuleRepository,
	ts timescale_repository.TimescaleRepository,
	ad AlertDispatcher,
	rc *redis_repository.RedisRepeatCounter,
	rd *redis_repository.RedisCache,
	logger *zerolog.Logger,
) *EvaluateRulesUseCase {
	return &EvaluateRulesUseCase{
		ruleRepo:        rr,
		timeScaleRepo:   ts,
		alertDispatcher: ad,
		redisCounter:    rc,
		redisCache:      rd,
		logger:          logger,
	}
}

// Evaluate – читает event.UserID, event.ServiceName -> загружает только нужные правила
func (uc *EvaluateRulesUseCase) Evaluate(ctx context.Context, event *domain.Event) error {
	uc.logger.Debug().Msgf("Evaluate: user=%s, service=%s", event.UserID, event.ServiceName)

	// 1. Фильтруем только правила для (user_id, service_name)
	rules, err := uc.GetRulesByUserAndServiceAndProjectId(ctx, event.UserID, event.ServiceName, event.ProjectId)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to fetch rules by user & service")
		return err
	}
	uc.logger.Debug().Msgf("Got %d rules for user=%s, service=%s", len(rules), event.UserID, event.ServiceName)

	// 2. Нужно ли обновлять Redis (repeat_over)?
	// if hasRepeatOverRule(rules) {
	//     if err = uc.redisCounter.UpdateKey(ctx, &event); err != nil {
	//         uc.logger.Error().Err(err).Msg("Failed to update Redis for repeat_over")
	//     }
	// }

	// 3. Готовим evaluator
	evaluator := &RuleConditionEvaluator{
		redisCounter: uc.redisCounter,
		logger:       uc.logger,
	}

	// 4. Для каждого правила EvaluateRule -> собираем actions
	var triggered []domain.Action
	var triggeredRuleNames []domain.Rule
	for _, r := range rules {
		ok := domain.EvaluateRule(event, r, evaluator)
		if ok {
			uc.logger.Debug().Msgf("Rule matched: %s", r.Name)
			triggered = append(triggered, r.Actions...)
			triggeredRuleNames = append(triggeredRuleNames, r)
		}
	}

	// 5. Если есть actions, вызываем dispatcher
	if len(triggered) > 0 {
		uc.logger.Info().Msgf("Triggered %d actions for user=%s, service=%s", len(triggered), event.UserID, event.ServiceName)
		err = uc.alertDispatcher.DispatchActions(ctx, event, triggered)
		if err != nil {
			return err
		}
	}
	uc.logger.Debug().Msg("No rules matched, no actions triggered")
	// не добавляем лог в timescale если не сработало правило. Сейчас такая логика
	if len(triggeredRuleNames) == 0 {
		return nil
	}
	// 6. Сохраняем лог в TimescaleDB
	raw, err := json.Marshal(event)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to marshal event to JSON")
		return err
	}

	userIDInt, err := strconv.Atoi(event.UserID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to convert UserID to int")
		return err
	}

	logEntry := timescale_repository.LogEntry{
		UserID:      userIDInt,
		ServiceName: event.ServiceName,
		Timestamp:   time.Now(),
		Log:         raw,
		EventType:   event.EventType,
		UsedRules:   triggeredRuleNames,
		Language:    event.Language,
		ActionUsed:  triggered,
		ProjectId:   event.ProjectId,
		Engine:      ENGINE,
	}

	err = uc.timeScaleRepo.InsertLog(ctx, logEntry)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to insert log into TimescaleDB")
		// В зависимости от требований, можно не прерывать выполнение
		// return err
	}

	return nil
}

func (uc *EvaluateRulesUseCase) GetRulesByUserAndServiceAndProjectId(ctx context.Context, userID, serviceName, projectId string) ([]domain.Rule, error) {
	// Сначала пробуем получить правила из кеша
	if uc.redisCache != nil {
		if cachedRules, err := uc.redisCache.GetRules(ctx, userID, serviceName, projectId); err == nil && cachedRules != nil {
			uc.logger.Debug().Msg("Returning rules from cache")
			return cachedRules, nil
		}
	}

	// Если кеш промахнулся, достаём правила из MongoDB
	rules, err := uc.ruleRepo.GetRulesByUserAndServiceAndProjectId(ctx, userID, serviceName, projectId)
	if err != nil {
		return nil, err
	}

	// Сохраняем полученные правила в кеш для следующих запросов
	if uc.redisCache != nil {
		if err := uc.redisCache.SetRules(ctx, userID, serviceName, projectId, rules); err != nil {
			uc.logger.Error().Err(err).Msg("Failed to set rules in cache")
		}
	}

	return rules, nil
}
