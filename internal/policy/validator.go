package policy

import (
	"fmt"

	"github.com/mabuhasna8/IntelliOps/internal/contracts"
)

func Validate(policy contracts.DecisionPolicy) error {
	if policy.ID == "" {
		return fmt.Errorf("policy ID is required")
	}

	if policy.Version <= 0 {
		return fmt.Errorf("policy version must be positive")
	}

	if policy.Default.Action == "" {
		return fmt.Errorf("default decision action is required")
	}

	seen := make(map[string]struct{})

	for _, rule := range policy.Rules {
		if rule.ID == "" {
			return fmt.Errorf("rule ID is required")
		}

		if _, exists := seen[rule.ID]; exists {
			return fmt.Errorf("duplicate rule ID: %s", rule.ID)
		}

		seen[rule.ID] = struct{}{}

		if rule.Decision.Action == "" {
			return fmt.Errorf("decision action is required for rule %s", rule.ID)
		}

		for _, condition := range rule.Conditions {
			if condition.Field == "" {
				return fmt.Errorf("condition field is required in rule %s", rule.ID)
			}

			switch condition.Operator {
			case "equals", "not_equals", "exists", "in":
			default:
				return fmt.Errorf(
					"unsupported operator %q in rule %s",
					condition.Operator,
					rule.ID,
				)
			}
		}
	}

	return nil
}
