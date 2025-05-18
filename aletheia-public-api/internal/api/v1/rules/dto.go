package rules

type Rule struct {
	Id       int64  `json:"id"`
	RuleName string `json:"ruleName"`
	RuleType string `json:"ruleType"`
}

type AvailableRules struct {
	Rules []Rule `json:"rules"`
}

type Request struct {
	UserId   int64  `json:"userId"`
	RuleId   int64  `json:"ruleId"`
	RuleType string `json:"ruleType"`
}
