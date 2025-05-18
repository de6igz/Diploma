package usecases

import (
	"context"
	"strconv"
	"strings"

	"rule-engine-resources/internal/dataproviders/redis_repository"
	"rule-engine-resources/internal/domain"

	"github.com/rs/zerolog"
)

// RuleConditionEvaluator отвечает за проверку одиночного условия (Condition) на событии (Event).
type RuleConditionEvaluator struct {
	redisCounter *redis_repository.RedisRepeatCounter
	logger       *zerolog.Logger
}
type ConditionOperator string

// Evaluate проверяет, выполняется ли условие cond для события e.
// Поддерживаются все операторы, описанные в domain.
func (rce *RuleConditionEvaluator) Evaluate(e *domain.Event, c domain.Condition, r domain.Rule) bool {
	rce.logger.Debug().Msgf("Check condition: field=%s operator=%s value=%v", c.Field, c.Operator, c.Value)
	val := rce.getDynamicField(e, c.Field) // извлекаем значение по dot-path
	switch c.Operator {
	case domain.OpEQ:
		return isEqual(val, c.Value)
	case domain.OpNEQ:
		return !isEqual(val, c.Value)
	case domain.OpGT:
		return compare(val, c.Value) == 1
	case domain.OpGTE:
		cmp := compare(val, c.Value)
		return cmp == 1 || cmp == 0
	case domain.OpLT:
		return compare(val, c.Value) == -1
	case domain.OpLTE:
		cmp := compare(val, c.Value)
		return cmp == -1 || cmp == 0
	case domain.OpIN:
		return inList(val, c.Value, true)
	case domain.OpNIN:
		return inList(val, c.Value, false)
	case domain.OpCont:
		return evaluateContains(val, c.Value)
	case domain.OpRepeatOver:
		// Пример: "value": { "threshold": 3, "minutes": 1 }
		valMap, ok := c.Value.(map[string]interface{})
		if !ok {
			rce.logger.Warn().Msg("repeat_over condition has invalid 'value' (expected map[string]interface{})")
			return false
		}

		// Получаем значения. Учтём, что числа могут прийти как int или float64.
		var threshold int
		switch t := valMap["threshold"].(type) {
		case int:
			threshold = t
		case int32:
			threshold = int(t)
		case int64:
			threshold = int(t)
		case float64:
			threshold = int(t)
		default:
			rce.logger.Warn().Msg("repeat_over condition: invalid type for 'threshold'")
			return false
		}

		var minutes int
		switch m := valMap["minutes"].(type) {
		case int:
			minutes = m
		case int32:
			minutes = int(m)
		case int64:
			minutes = int(m)
		case float64:
			minutes = int(m)
		default:
			rce.logger.Warn().Msg("repeat_over condition: invalid type for 'minutes'")
			return false
		}

		// Получаем счетчик за окно времени minutes
		cnt, err := rce.redisCounter.CountInWindow(context.Background(), e, r, minutes)
		if err != nil {
			rce.logger.Error().Err(err).Msg("CountInWindow failed")
			return false
		}
		e.RepeatCount = cnt

		rce.logger.Debug().Msgf("repeat_over check: count=%d, threshold=%d", cnt, threshold)
		return cnt >= threshold
	default:
		return false
	}
}

// getDynamicField(evt, "fields.memory_alloc_bytes") – вытягивает значение
// Считаем, что всё, что не "user_id"/"service_name" – внутри evt.Fields
func (rce *RuleConditionEvaluator) getDynamicField(evt *domain.Event, fieldPath string) interface{} {

	rce.logger.Debug().Msgf("getDynamicField: fieldPath=%q", fieldPath)

	parts := strings.Split(fieldPath, ".")
	if len(parts) == 0 {
		rce.logger.Debug().Msgf("getDynamicField: parts is empty => return nil")
		return nil
	}

	// пробуем user_id / service_name
	if parts[0] == "user_id" {
		if len(parts) == 1 {
			rce.logger.Debug().Msgf("getDynamicField: return evt.UserID = %s", evt.UserID)
			return evt.UserID
		}
		rce.logger.Debug().Msg("getDynamicField: user_id has more sub-parts => nil")
		return nil
	}
	if parts[0] == "service_name" {
		if len(parts) == 1 {
			rce.logger.Debug().Msgf("getDynamicField: return evt.ServiceName = %s", evt.ServiceName)
			return evt.ServiceName
		}
		rce.logger.Debug().Msg("getDynamicField: service_name has more sub-parts => nil")
		return nil
	}

	// остальное => fields
	if evt.Fields == nil {
		rce.logger.Debug().Msg("getDynamicField: evt.Fields == nil => return nil")
		return nil
	}

	// проверяем, что parts[0] == "fields" (или что угодно)
	// но чаще всего ожидаем "fields" как 1й сегмент
	if parts[0] != "fields" {
		rce.logger.Debug().Msgf("getDynamicField: first part=%q != 'fields', return nil", parts[0])
		return nil
	}

	subParts := parts[1:]
	// рекурсивно
	return rce.traverseMapDebug(evt.Fields, subParts)
}

func (rce *RuleConditionEvaluator) traverseMapDebug(cur interface{}, keys []string) interface{} {
	if len(keys) == 0 {
		rce.logger.Debug().Msgf("traverseMapDebug: no more keys => return %+v", cur)
		return cur
	}

	switch m := cur.(type) {
	case map[string]interface{}:
		k := keys[0]
		val, ok := m[k]
		rce.logger.Debug().Msgf("traverseMapDebug: Looking for key=%q in map => found? %v", k, ok)
		if !ok {
			return nil
		}
		return rce.traverseMapDebug(val, keys[1:])
	default:
		rce.logger.Debug().Msg("traverseMapDebug: current is not map => return nil")
		return nil
	}
}

// isEqual упрощённо сравнивает
func isEqual(a, b interface{}) bool {
	switch av := a.(type) {
	case string:
		bv, ok := b.(string)
		return ok && av == bv
	case float64:
		switch bv := b.(type) {
		case float64:
			return av == bv
		case int:
			return av == float64(bv)
		}
		return false
	case int:
		switch bv := b.(type) {
		case float64:
			return float64(av) == bv
		case int:
			return av == bv
		}
		return false
	}
	// fallback
	return a == b
}

// compare – вернёт -1 (a<b), 0 (a==b), 1 (a>b), или 99 для ошибок
func compare(a, b interface{}) int {
	af, aok := toFloat(a)
	bf, bok := toFloat(b)
	if aok && bok {
		if af < bf {
			return -1
		} else if af > bf {
			return 1
		}
		return 0
	}
	// если строки
	as, saok := a.(string)
	bs, sbok := b.(string)
	if saok && sbok {
		if as < bs {
			return -1
		} else if as > bs {
			return 1
		} else {
			return 0
		}
	}
	return 99
}

func toFloat(x interface{}) (float64, bool) {
	switch v := x.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		f, err := strconv.ParseFloat(v, 64)
		return f, err == nil
	}
	return 0, false
}

func inList(fieldVal, condVal interface{}, in bool) bool {
	switch arr := condVal.(type) {
	case []interface{}:
		for _, el := range arr {
			if isEqual(fieldVal, el) {
				return in
			}
		}
		return !in
	case []string:
		for _, s := range arr {
			if isEqual(fieldVal, s) {
				return in
			}
		}
		return !in
	}
	return false
}

func evaluateContains(fieldVal, condVal interface{}) bool {
	// 1) []string + string
	if arr, ok := fieldVal.([]interface{}); ok {
		cvStr, ok2 := condVal.(string)
		if !ok2 {
			return false
		}
		for _, a := range arr {
			if isEqual(a, cvStr) {
				return true
			}
		}
		return false
	}
	// 1b) []string
	if strArr, ok := fieldVal.([]string); ok {
		cvStr, ok2 := condVal.(string)
		if !ok2 {
			return false
		}
		for _, s := range strArr {
			if s == cvStr {
				return true
			}
		}
		return false
	}
	// 2) string + string => substring
	if fvStr, ok := fieldVal.(string); ok {
		if cvStr, ok2 := condVal.(string); ok2 {
			return strings.Contains(fvStr, cvStr)
		}
	}
	return false
}
