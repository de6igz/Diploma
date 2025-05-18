package rules_errors

type RuleData struct {
	Id          string  `json:"id"`
	RuleName    string  `json:"rule_name"`
	Description *string `json:"description"`
	RuleType    string  `json:"rule_type"`
}

type RuleById struct {
	RuleName        string   `json:"name"`
	RuleDescription string   `json:"description"`
	RuleType        string   `json:"ruleType"`
	RootNode        Node     `json:"root_node"`
	Actions         []Action `json:"actions"`
}

type Node struct {
	Operator   string      `json:"operator"`
	Conditions []Condition `json:"conditions"`
	Children   []Node      `json:"children"`
}

type Condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Action struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}
