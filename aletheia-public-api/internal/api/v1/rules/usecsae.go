package rules

import (
	v1 "aletheia-public-api/interfaces/types/v1"
	"aletheia-public-api/internal/dataproviders/postgres/repositories/rules_errors"
	"aletheia-public-api/internal/dataproviders/postgres/repositories/rules_resources"
	"context"
	"fmt"
	"strconv"
)

// RulesUsecase описывает методы для получения проектов.
type RulesUsecase interface {
	GetAvailableRules(ctx context.Context, userId int64) (AvailableRules, error)
	DeleteRuleById(ctx context.Context, userId int64, request v1.DeleteRuleRequest) error
	CreateRule(ctx context.Context, userId int64, request v1.CreateRuleRequest) error
	UpdateRuleById(ctx context.Context, userId int64, request v1.UpdateRuleRequest) error
	GetRules(ctx context.Context, userId int64) (v1.RulesResponse, error)
	GetRuleById(ctx context.Context, userId int64, request v1.RuleByIdRequest) (*v1.RuleDetailResponse, error)
}

type rulesUsecase struct {
	rulesErrorsRepo    rules_errors.Provider
	rulesResourcesRepo rules_resources.Provider
}

// NewRulesUsecase  создаёт usecase с инъекцией репозитория.
func NewRulesUsecase(
	rulesErrorsRepo rules_errors.Provider,
	rulesResourcesRepo rules_resources.Provider,
) RulesUsecase {
	return &rulesUsecase{
		rulesErrorsRepo:    rulesErrorsRepo,
		rulesResourcesRepo: rulesResourcesRepo,
	}
}

func (r *rulesUsecase) GetAvailableRules(ctx context.Context, userId int64) (AvailableRules, error) {
	// Формируем запрос для репозиториев: передаём только userId.

	// Получаем свободные error-правила.
	errorRulesData, err := r.rulesErrorsRepo.GetAvailableRulesData(ctx, rules_errors.Request{UserId: userId})
	if err != nil {
		return AvailableRules{}, fmt.Errorf("error fetching available error rules: %w", err)
	}

	// Получаем свободные resource-правила.
	resourceRulesData, err := r.rulesResourcesRepo.GetAvailableRulesData(ctx, rules_resources.Request{UserId: userId})
	if err != nil {
		return AvailableRules{}, fmt.Errorf("error fetching available resource rules: %w", err)
	}

	// Преобразуем данные из репозиториев в единый список.
	var combined []Rule
	for _, rd := range errorRulesData {
		// Преобразуем id из строки в int64.
		id, err := strconv.ParseInt(rd.Id, 10, 64)
		if err != nil {
			// Если преобразование не удалось, можно пропустить это правило.
			continue
		}
		combined = append(combined, Rule{
			Id:       id,
			RuleName: rd.RuleName,
			RuleType: rd.RuleType,
		})
	}
	for _, rd := range resourceRulesData {
		id, err := strconv.ParseInt(rd.Id, 10, 64)
		if err != nil {
			continue
		}
		combined = append(combined, Rule{
			Id:       id,
			RuleName: rd.RuleName,
			RuleType: rd.RuleType,
		})
	}

	available := AvailableRules{
		Rules: combined,
	}

	// Возвращаем список из одного AvailableRules.
	return available, nil
}

func (r *rulesUsecase) DeleteRuleById(ctx context.Context, userId int64, request v1.DeleteRuleRequest) error {
	if request.RuleType == "" {
		return fmt.Errorf("rule type is required")
	}

	if request.RuleType == "errors" {
		err := r.rulesErrorsRepo.DeleteRuleById(ctx, request.RuleId, request.RuleType, userId)
		if err != nil {
			return fmt.Errorf("error deleting error rule: %w", err)
		}
		return nil
	}

	if request.RuleType == "resources" {
		err := r.rulesResourcesRepo.DeleteRuleById(ctx, request.RuleId, request.RuleType, userId)
		if err != nil {
			return fmt.Errorf("error deleting error rule: %w", err)
		}
		return nil
	}
	return nil
}

func (r *rulesUsecase) CreateRule(ctx context.Context, userId int64, request v1.CreateRuleRequest) error {
	if request.RuleType == "" {
		return fmt.Errorf("rule type is required")
	}
	if request.RuleType == "errors" {
		err := r.rulesErrorsRepo.CreateRule(ctx, userId, request)
		if err != nil {
			return fmt.Errorf("error creating error rule: %w", err)
		}
		return nil
	}
	if request.RuleType == "resources" {
		err := r.rulesResourcesRepo.CreateRule(ctx, userId, request)
		if err != nil {
			return fmt.Errorf("error creating resource rule: %w", err)
		}
		return nil
	} else {
		return fmt.Errorf("invalid rule type")
	}
}

func (r *rulesUsecase) UpdateRuleById(ctx context.Context, userId int64, request v1.UpdateRuleRequest) error {
	if request.RuleType == "" {
		return fmt.Errorf("rule type is required")
	}
	if request.RuleType == "errors" {
		err := r.rulesErrorsRepo.UpdateRuleById(ctx, userId, request)
		if err != nil {
			return fmt.Errorf("error updating error rule: %w", err)
		}
		return nil
	}
	if request.RuleType == "resources" {
		err := r.rulesResourcesRepo.UpdateRuleById(ctx, userId, request)
		if err != nil {
			return fmt.Errorf("error updating resource rule: %w", err)
		}
		return nil
	}
	return nil
}

// GetRules получает все правила для пользователя, объединяя данные из репозиториев error и resource.
func (r *rulesUsecase) GetRules(ctx context.Context, userId int64) (v1.RulesResponse, error) {
	// Получаем все error-правила для пользователя.
	errorRulesData, err := r.rulesErrorsRepo.GetErrorRulesData(ctx, rules_errors.Request{UserId: userId})
	if err != nil {
		return v1.RulesResponse{}, fmt.Errorf("failed to get error rules: %w", err)
	}

	// Получаем все resource-правила для пользователя.
	resourceRulesData, err := r.rulesResourcesRepo.GetResourcesRulesData(ctx, rules_resources.Request{UserId: userId})
	if err != nil {
		return v1.RulesResponse{}, fmt.Errorf("failed to get resource rules: %w", err)
	}

	// Преобразуем данные из репозиториев в единый срез правил для ответа.
	var combined []*v1.Rule
	if len(errorRulesData) > 0 {
		for _, rd := range errorRulesData {
			// Преобразуем id из строки (хранящийся в базе) в нужный формат.
			id, err := strconv.ParseInt(rd.Id, 10, 64)
			if err != nil {
				// Если преобразование не удалось, пропускаем правило.
				continue
			}
			combined = append(combined, &v1.Rule{
				ID:          strconv.FormatInt(id, 10),
				Name:        rd.RuleName,
				RuleType:    &rd.RuleType,
				Description: rd.Description,
			})
		}
	}
	if len(resourceRulesData) > 0 {
		for _, rd := range resourceRulesData {
			id, err := strconv.ParseInt(rd.Id, 10, 64)
			if err != nil {
				continue
			}
			combined = append(combined, &v1.Rule{
				ID:          strconv.FormatInt(id, 10),
				Name:        rd.RuleName,
				RuleType:    &rd.RuleType,
				Description: rd.Description,
			})
		}
	}
	return v1.RulesResponse{Rules: combined}, nil
}

func (r *rulesUsecase) GetRuleById(ctx context.Context, userId int64, request v1.RuleByIdRequest) (*v1.RuleDetailResponse, error) {
	if request.RuleType == "" {
		return nil, fmt.Errorf("rule type is required")
	}
	if request.RuleType == "errors" {
		res, err := r.rulesErrorsRepo.GetRuleById(ctx, request.RuleId, userId)
		if err != nil {
			return nil, fmt.Errorf("error getting rule by id: %w", err)
		}
		return res, nil
	}

	if request.RuleType == "resources" {
		res, err := r.rulesResourcesRepo.GetRuleById(ctx, request.RuleId, userId)
		if err != nil {
			return nil, fmt.Errorf("error getting rule by id: %w", err)
		}
		return res, nil
	}
	return nil, fmt.Errorf("invalid rule type")
}
