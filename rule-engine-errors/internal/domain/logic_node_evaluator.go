package domain

// ConditionEvaluator – интерфейс/или структурный тип
// для проверки одиночного Condition (как в предыдущих примерах).
//type ConditionEvaluator interface {
//	EvaluateCondition(e Event, c Condition) bool
//}

// EvaluateLogicNode – рекурсивно проверяет LogicNode.
// Возвращает true, если узел "выполнился".
func EvaluateLogicNode(e *Event, node LogicNode, evaluator ConditionEvaluator, r Rule) bool {
	op := node.Operator // "AND" или "OR"

	// 1. Проверим все conditions в этом узле
	var condResult bool
	if op == "AND" {
		// Нужно, чтобы все были true
		condResult = true
		for _, c := range node.Conditions {
			if !evaluator.Evaluate(e, c, r) {
				condResult = false
				break
			}
		}
	} else {
		// Считаем "OR" по умолчанию
		condResult = false
		for _, c := range node.Conditions {
			if evaluator.Evaluate(e, c, r) {
				condResult = true
				break
			}
		}
	}

	// 2. Проверим children
	var childResult bool
	if op == "AND" {
		childResult = true
		for _, child := range node.Children {
			if !EvaluateLogicNode(e, child, evaluator, r) {
				childResult = false
				break
			}
		}
	} else {
		childResult = false
		for _, child := range node.Children {
			if EvaluateLogicNode(e, child, evaluator, r) {
				childResult = true
				break
			}
		}
	}

	// 3. Соединяем results:
	// При "AND" нужно, чтобы condResult и childResult были true.
	// При "OR" – достаточно одного.
	if op == "AND" {
		return condResult && childResult
	} else {
		return condResult || childResult
	}
}
