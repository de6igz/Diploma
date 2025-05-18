package utils

import (
	"encoding/json"
	"fmt"
)

func MergeContext(context map[string]interface{}, userID string, userInfo map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	// Копируем исходный контекст
	for k, v := range context {
		merged[k] = v
	}

	// Добавляем пользовательские данные
	if userID != "" {
		merged["user_id"] = userID
	}
	if userInfo != nil {
		merged["user_info"] = userInfo
	}

	return merged
}

func BuildContext(ctx map[string]interface{}, userID string, userInfo map[string]interface{}) string {
	if ctx == nil {
		ctx = make(map[string]interface{})
	}
	if userID != "" {
		ctx["user_id"] = userID
	}
	if userInfo != nil {
		ctx["user_info"] = userInfo
	}

	bytes, _ := json.Marshal(ctx)
	return string(bytes)
}

func FormatTags(tags map[string]string) []string {
	result := make([]string, 0, len(tags))
	for k, v := range tags {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}
