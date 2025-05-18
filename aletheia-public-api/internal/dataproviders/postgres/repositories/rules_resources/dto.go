package rules_resources

type RuleData struct {
	Id          string  `json:"id"`
	RuleName    string  `json:"rule_name"`
	Description *string `json:"description"`
	RuleType    string  `json:"rule_type"`
}
