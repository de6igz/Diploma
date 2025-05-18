package usecases

import (
	"context"
	"strings"

	"rule-engine-errors/internal/dataproviders/redis_repository"
	"rule-engine-errors/internal/domain"

	"github.com/rs/zerolog"
)

// RuleConditionEvaluator отвечает за проверку одиночного условия (Condition) на событии (Event).
type RuleConditionEvaluator struct {
	redisCounter *redis_repository.RedisRepeatCounter
	logger       *zerolog.Logger
}

// Evaluate проверяет, выполняется ли условие cond для события e.
// Для фиксированных полей используется старая логика, для динамических (начинающихся с "fields.")
// — логика обхода вложенной мапы (аналогичная resources).
func (rce *RuleConditionEvaluator) Evaluate(e *domain.Event, c domain.Condition, r domain.Rule) bool {
	rce.logger.Debug().Msgf("Check condition: field=%s operator=%s value=%v", c.Field, c.Operator, c.Value)

	switch c.Operator {

	// --- "repeat_over" оператор (подсчёт повторов через Redis) ---
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

	// --- "eq" (равно) ---
	case domain.OpEQ:
		return isEqual(rce.getField(e, c.Field), c.Value)

	// --- "neq" (не равно) ---
	case domain.OpNEQ:
		return !isEqual(rce.getField(e, c.Field), c.Value)

	// --- "gt" (строго больше) ---
	case domain.OpGT:
		return compareNumericOrString(rce.getField(e, c.Field), c.Value) == 1

	// --- "gte" (больше или равно) ---
	case domain.OpGTE:
		cmp := compareNumericOrString(rce.getField(e, c.Field), c.Value)
		return cmp == 1 || cmp == 0

	// --- "lt" (строго меньше) ---
	case domain.OpLT:
		return compareNumericOrString(rce.getField(e, c.Field), c.Value) == -1

	// --- "lte" (меньше или равно) ---
	case domain.OpLTE:
		cmp := compareNumericOrString(rce.getField(e, c.Field), c.Value)
		return cmp == -1 || cmp == 0

	// --- "in" (входит в список) ---
	case domain.OpIN:
		return inList(rce.getField(e, c.Field), c.Value, true)

	// --- "nin" (не входит в список) ---
	case domain.OpNIN:
		return inList(rce.getField(e, c.Field), c.Value, false)

	// --- "contains" (содержит) ---
	case domain.OpCont:
		return evaluateContains(rce.getField(e, c.Field), c.Value)

	// --- неизвестный оператор ---
	default:
		rce.logger.Debug().Msgf("Unsupported operator: %s", c.Operator)
		return false
	}
}

// getField выбирает способ получения значения поля: если имя поля начинается с "fields.",
// то используется динамический поиск в evt.Fields, иначе — фиксированная логика.
func (rce *RuleConditionEvaluator) getField(e *domain.Event, field string) interface{} {
	if strings.HasPrefix(field, "fields.") {
		return rce.getDynamicField(e, field)
	}
	return getFieldValue(e, field)
}

// getFieldValue возвращает значение фиксированного поля события.
func getFieldValue(e *domain.Event, field string) interface{} {
	switch field {
	case "user_id":
		return e.UserID
	case "service_name":
		return e.ServiceName
	case "environment":
		return e.Environment
	case "error_message":
		return e.ErrorMessage
	case "version":
		return e.Version
	case "go_version":
		return e.GoVersion
	case "os":
		return e.Os
	case "arch":
		return e.Arch
	case "event_type":
		return e.EventType
	case "event_message":
		return e.EventMessage
	case "stack_trace":
		return e.StackTrace
	case "tags":
		return e.Tags
	case "timestamp":
		return e.Timestamp
	case "context_json":
		return e.ContextJson
	case "repeat_count":
		return e.RepeatCount
	case "level":
		return e.Level
	default:
		// Если поле не найдено среди фиксированных, возвращаем nil.
		return nil
	}
}

// ===================== Динамическая проверка полей =====================

// getDynamicField извлекает значение динамического поля по dot-path (например, "fields.memory_alloc_bytes").
// Если первый сегмент не равен "fields", возвращается nil.
func (rce *RuleConditionEvaluator) getDynamicField(evt *domain.Event, fieldPath string) interface{} {
	rce.logger.Debug().Msgf("getDynamicField: fieldPath=%q", fieldPath)
	parts := strings.Split(fieldPath, ".")
	if len(parts) == 0 {
		rce.logger.Debug().Msg("getDynamicField: parts is empty => return nil")
		return nil
	}

	// Дополнительная обработка для фиксированных полей (на случай, если кто-то передаст "user_id" с точками)
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

	// Для динамических полей ожидаем, что первый сегмент равен "fields"
	if evt.Fields == nil {
		rce.logger.Debug().Msg("getDynamicField: evt.Fields == nil => return nil")
		return nil
	}
	if parts[0] != "fields" {
		rce.logger.Debug().Msgf("getDynamicField: first part=%q != 'fields', return nil", parts[0])
		return nil
	}

	subParts := parts[1:]
	return rce.traverseMapDebug(evt.Fields, subParts)
}

// traverseMapDebug рекурсивно проходит по вложенным мапам, согласно заданному пути.
func (rce *RuleConditionEvaluator) traverseMapDebug(cur interface{}, keys []string) interface{} {
	if len(keys) == 0 {
		rce.logger.Debug().Msgf("traverseMapDebug: no more keys => return %+v", cur)
		return cur
	}

	switch m := cur.(type) {
	case map[string]interface{}:
		k := keys[0]
		val, ok := m["fields."+k]
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

// ===================== Вспомогательные функции =====================

// isEqual сравнивает простые значения (строки, float64, int и т.д.) на равенство.
// Если типы не совпадают, возвращается false.
func isEqual(a, b interface{}) bool {
	switch av := a.(type) {
	case string:
		bv, ok := b.(string)
		return ok && av == bv
	case float64:
		bv, ok := b.(float64)
		return ok && av == bv
	case int:
		bv, ok := b.(int)
		return ok && av == bv
	default:
		return a == b
	}
}

// compareNumericOrString возвращает -1, 0 или 1 (аналог strcmp).
// Если a < b — возвращается -1, если a > b — 1, если a == b — 0. При несовместимых типах — 99.
func compareNumericOrString(a, b interface{}) int {
	if af, ok := toFloat(a); ok {
		if bf, ok2 := toFloat(b); ok2 {
			switch {
			case af < bf:
				return -1
			case af > bf:
				return 1
			default:
				return 0
			}
		}
		return 99
	}

	as, aok := a.(string)
	bs, bok := b.(string)
	if aok && bok {
		if as < bs {
			return -1
		}
		if as > bs {
			return 1
		}
		return 0
	}

	return 99
}

// toFloat пытается преобразовать значение к float64.
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
	}
	return 0, false
}

// inList проверяет, содержится ли значение поля (fieldVal) в списке значений (condVal).
// Параметр in задаёт: для IN — true, для NIN — false.
func inList(fieldVal, condVal interface{}, in bool) bool {
	switch arr := condVal.(type) {
	case []interface{}:
		for _, el := range arr {
			if isEqual(fieldVal, el) {
				return in // найдено — возвращаем in (true для IN, false для NIN)
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
	default:
		return false
	}
}

// evaluateContains проверяет, содержит ли значение поля подстроку или элемент.
func evaluateContains(fieldVal interface{}, condVal interface{}) bool {
	// Если fieldVal — массив строк
	if arr, ok := fieldVal.([]string); ok {
		cstr, ok2 := condVal.(string)
		if !ok2 {
			return false
		}
		for _, v := range arr {
			if v == cstr {
				return true
			}
		}
		return false
	}

	// Если fieldVal — строка, проверяем наличие подстроки
	if fvStr, ok := fieldVal.(string); ok {
		if cvStr, ok2 := condVal.(string); ok2 {
			return strings.Contains(fvStr, cvStr)
		}
		return false
	}

	return false
}
