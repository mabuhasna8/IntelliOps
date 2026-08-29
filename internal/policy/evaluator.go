package policy

import "github.com/mabuhasna8/IntelliOps/internal/contracts"

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Evaluate(
	policy contracts.DecisionPolicy,
	input map[string]any,
) contracts.DecisionEvaluation {
	evaluated := make([]string, 0, len(policy.Rules))

	rules := sortRules(policy.Rules)

	for _, rule := range rules {
		evaluated = append(evaluated, rule.ID)

		if matchesAll(rule.Conditions, input) {
			return contracts.DecisionEvaluation{
				Decision:       rule.Decision,
				MatchedRuleID:  rule.ID,
				EvaluatedRules: evaluated,
			}
		}
	}

	return contracts.DecisionEvaluation{
		Decision:       policy.Default,
		EvaluatedRules: evaluated,
	}
}
