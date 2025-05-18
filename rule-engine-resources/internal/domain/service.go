package domain

// EvaluateRule – проверяем, выполняются ли все условия
//
//	func EvaluateRule(e Event, r Rule, evaluator ConditionEvaluator) bool {
//		for _, c := range r.Conditions {
//			if !evaluator.Evaluate(e, c, r) {
//				return false
//			}
//		}
//		return true
//	}
func EvaluateRule(e *Event, r Rule, evaluator ConditionEvaluator) bool {
	return EvaluateLogicNode(e, r.RootNode, evaluator, r)
}

// ConditionEvaluator – интерфейс для проверки одного Condition
// (можем внедрять RedisRepeatCounter или что-то ещё)
type ConditionEvaluator interface {
	Evaluate(e *Event, c Condition, r Rule) bool
}

//func evaluateCondition(event Event, cond Condition) bool {
//	fieldValue := getFieldValue(event, cond.Field)
//	switch cond.Operator {
//	case OpEQ:
//		return fieldValue == cond.Value
//	case OpNEQ:
//		return fieldValue != cond.Value
//	case OpIN:
//		return evaluateIn(fieldValue, cond.Value, true)
//	case OpNIN:
//		return evaluateIn(fieldValue, cond.Value, false)
//	case OpCont:
//		return evaluateContains(fieldValue, cond.Value)
//	default:
//		// упрощенно
//		return false
//	}
//}

//// getFieldValue – извлекаем нужное поле из event
//func getFieldValue(event Event, field string) interface{} {
//	switch field {
//	case "service_name":
//		return event.ServiceName
//	case "environment":
//		return event.Environment
//	case "level":
//		return event.Level
//	case "event_type":
//		return event.EventType
//	case "event_message":
//		return event.EventMessage
//	case "tags":
//		return event.Tags
//	}
//	// ...
//	return nil
//}
//
//func evaluateIn(fieldValue interface{}, condValue interface{}, in bool) bool {
//	arr, ok := condValue.([]interface{})
//	if !ok {
//		return false
//	}
//	for _, v := range arr {
//		if v == fieldValue {
//			return in
//		}
//	}
//	return !in
//}
//
//func evaluateContains(fv interface{}, cv interface{}) bool {
//	if arr, ok := fv.([]string); ok {
//		strVal, ok2 := cv.(string)
//		if !ok2 {
//			return false
//		}
//		for _, tag := range arr {
//			if tag == strVal {
//				return true
//			}
//		}
//		return false
//	}
//	if strField, ok := fv.(string); ok {
//		strCond, ok2 := cv.(string)
//		if !ok2 {
//			return false
//		}
//		return strings.Contains(strField, strCond)
//	}
//	return false
//}
