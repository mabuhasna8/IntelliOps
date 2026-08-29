package contracts

type DecisionPolicy struct {
	ID      string       `json:"id"`
	Version int          `json:"version"`
	Rules   []PolicyRule `json:"rules"`
	Default Decision     `json:"default"`
}

type PolicyRule struct {
	ID         string      `json:"id"`
	Conditions []Condition `json:"conditions"`
	Decision   Decision    `json:"decision"`
	Priority   int         `json:"priority"`
}

type Condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

type DecisionEvaluation struct {
	Decision       Decision `json:"decision"`
	MatchedRuleID  string   `json:"matched_rule_id,omitempty"`
	EvaluatedRules []string `json:"evaluated_rules"`
}
