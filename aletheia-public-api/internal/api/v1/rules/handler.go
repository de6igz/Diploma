package rules

import (
	v1 "aletheia-public-api/interfaces/types/v1"
	"aletheia-public-api/internal/dataproviders/postgres"
	rulesErrors "aletheia-public-api/internal/dataproviders/postgres/repositories/rules_errors"
	rulesResources "aletheia-public-api/internal/dataproviders/postgres/repositories/rules_resources"
	"context"
	"fmt"
	"strconv"
)

type Rules struct {
	usecase RulesUsecase
}

// NewRules создаёт новый обработчик с инъекцией репозиториев и юзкейса.
func NewRules() *Rules {
	errorRepo := rulesErrors.NewProvider(postgres.GlobalInstance)
	resourceRepo := rulesResources.NewProvider(postgres.GlobalInstance)
	usecase := NewRulesUsecase(errorRepo, resourceRepo)
	return &Rules{
		usecase: usecase,
	}
}

func (r *Rules) GetRules(ctx context.Context, userId int64) (items v1.RulesResponse, err error) {
	res, err := r.usecase.GetRules(ctx, userId)
	if err != nil {
		return items, err
	}
	return res, nil
}

func (r *Rules) GetRuleByID(ctx context.Context, userId int64, ruleId, ruleType string) (rule *v1.RuleDetailResponse, err error) {
	res, err := r.usecase.GetRuleById(ctx, userId, v1.RuleByIdRequest{RuleId: ruleId, RuleType: ruleType})
	if err != nil {
		return rule, err
	}
	return res, nil
}

// GetAvailableRules получает свободные правила (у которых service_id IS NULL) для пользователя.
func (r *Rules) GetAvailableRules(ctx context.Context, userId int64) (v1.RulesResponse, error) {
	available, err := r.usecase.GetAvailableRules(ctx, userId)
	if err != nil {
		return v1.RulesResponse{}, fmt.Errorf("failed to get available rules: %w", err)
	}

	result := make([]*v1.Rule, 0)
	// Объединяем все полученные правила из каждого элемента available.
	for _, avail := range available.Rules {
		result = append(result, &v1.Rule{
			ID:       strconv.FormatInt(avail.Id, 10),
			Name:     avail.RuleName,
			RuleType: &avail.RuleType,
		})
	}

	return v1.RulesResponse{Rules: result}, nil
}

func (r *Rules) DeleteRuleByID(ctx context.Context, userId int64, req v1.DeleteRuleRequest) (status bool, err error) {

	err = r.usecase.DeleteRuleById(ctx, userId, req)
	if err != nil {
		return false, fmt.Errorf("failed to delete rule: %w", err)
	}

	return true, nil
}

func (r *Rules) CreateRule(ctx context.Context, userId int64, req v1.CreateRuleRequest) (status bool, err error) {
	err = r.usecase.CreateRule(ctx, userId, req)
	if err != nil {
		return false, fmt.Errorf("failed to create rule: %w", err)
	}
	return true, nil

}
func (r *Rules) UpdateRuleById(ctx context.Context, userId int64, request v1.UpdateRuleRequest) (status bool, err error) {
	err = r.usecase.UpdateRuleById(ctx, userId, request)
	if err != nil {
		return false, fmt.Errorf("failed to update rule: %w", err)
	}
	return true, nil
}
