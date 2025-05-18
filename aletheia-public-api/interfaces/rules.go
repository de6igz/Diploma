package interfaces

import (
	v1 "aletheia-public-api/interfaces/types/v1"
	"context"
)

// Rules
// @tg http-server log metrics
// @tg http-prefix=v1
type Rules interface {
	// GetRules
	// @tg summary=`Получить правила`
	// @tg desc=`Возвращает список правил`
	// @tg http-method=GET
	// @tg http-path=/rules
	// @tg http-headers=userId|X-User-Id
	GetRules(ctx context.Context, userId int64) (items v1.RulesResponse, err error)

	// GetRuleByID
	// @tg summary=`Получить детали правила`
	// @tg desc=`Возвращает подробную информацию о правиле по его ID`
	// @tg http-method=GET
	// @tg http-path=/rule/byId
	// @tg http-headers=userId|X-User-Id
	// @tg http-args=`ruleId|ruleId`
	// @tg http-args=`ruleType|ruleType`
	GetRuleByID(ctx context.Context, userId int64, ruleId, ruleType string) (rule *v1.RuleDetailResponse, err error)
	// GetAvailableRules
	// @tg summary=`Получить свободные правила`
	// @tg desc=`Возвращает подробную информацию о правиле по его ID`
	// @tg http-method=GET
	// @tg http-path=/rules/available
	// @tg http-headers=userId|X-User-Id
	GetAvailableRules(ctx context.Context, userId int64) (items v1.RulesResponse, err error)
	// DeleteRuleByID
	// @tg summary=`Удалить правило по id`
	// @tg desc=`Удалить правило по id`
	// @tg http-method=DELETE
	// @tg http-path=/rules
	// @tg http-headers=userId|X-User-Id
	DeleteRuleByID(ctx context.Context, userId int64, req v1.DeleteRuleRequest) (status bool, err error)
	// CreateRule
	// @tg summary=`Создать правило`
	// @tg desc=`Создать правило`
	// @tg http-method=POST
	// @tg http-path=/rules/create
	// @tg http-headers=userId|X-User-Id
	CreateRule(ctx context.Context, userId int64, request v1.CreateRuleRequest) (status bool, err error)
	// UpdateRuleById
	// @tg summary=`Обновить правило`
	// @tg desc=`Обновить правило`
	// @tg http-method=PUT
	// @tg http-path=/rules/update
	// @tg http-headers=userId|X-User-Id
	UpdateRuleById(ctx context.Context, userId int64, request v1.UpdateRuleRequest) (status bool, err error)
}
