package policy

import (
	"testing"

	"github.com/mabuhasna8/IntelliOps/internal/contracts"
)

func TestEvaluateMatchesHighestPriorityRule(t *testing.T) {
	policy := contracts.DecisionPolicy{
		ID:      "build-policy",
		Version: 1,
		Rules: []contracts.PolicyRule{
			{
				ID:       "low-priority",
				Priority: 10,
				Conditions: []contracts.Condition{
					{
						Field:    "status",
						Operator: "equals",
						Value:    "failed",
					},
				},
				Decision: contracts.Decision{
					Action: "stop",
				},
			},
			{
				ID:       "high-priority",
				Priority: 20,
				Conditions: []contracts.Condition{
					{
						Field:    "status",
						Operator: "equals",
						Value:    "failed",
					},
				},
				Decision: contracts.Decision{
					Action: "retry",
				},
			},
		},
		Default: contracts.Decision{
			Action: "manual_review",
		},
	}

	result := NewEvaluator().Evaluate(policy, map[string]any{
		"status": "failed",
	})

	if result.Decision.Action != "retry" {
		t.Fatalf("expected retry, got %q", result.Decision.Action)
	}

	if result.MatchedRuleID != "high-priority" {
		t.Fatalf("expected high-priority rule, got %q", result.MatchedRuleID)
	}
}

func TestEvaluateUsesDefaultWhenNothingMatches(t *testing.T) {
	policy := contracts.DecisionPolicy{
		ID:      "test-policy",
		Version: 1,
		Rules: []contracts.PolicyRule{
			{
				ID: "failure-only",
				Conditions: []contracts.Condition{
					{
						Field:    "status",
						Operator: "equals",
						Value:    "failed",
					},
				},
				Decision: contracts.Decision{
					Action: "stop",
				},
			},
		},
		Default: contracts.Decision{
			Action: "manual_review",
		},
	}

	result := NewEvaluator().Evaluate(policy, map[string]any{
		"status": "succeeded",
	})

	if result.Decision.Action != "manual_review" {
		t.Fatalf("expected manual_review, got %q", result.Decision.Action)
	}

	if result.MatchedRuleID != "" {
		t.Fatalf("expected no matched rule, got %q", result.MatchedRuleID)
	}
}

func TestNestedFieldLookup(t *testing.T) {
	policy := contracts.DecisionPolicy{
		ID:      "nested-policy",
		Version: 1,
		Rules: []contracts.PolicyRule{
			{
				ID: "classification-rule",
				Conditions: []contracts.Condition{
					{
						Field:    "classification.code",
						Operator: "equals",
						Value:    "test.failure",
					},
				},
				Decision: contracts.Decision{
					Action: "stop",
				},
			},
		},
		Default: contracts.Decision{
			Action: "continue",
		},
	}

	result := NewEvaluator().Evaluate(policy, map[string]any{
		"classification": map[string]any{
			"code": "test.failure",
		},
	})

	if result.Decision.Action != "stop" {
		t.Fatalf("expected stop, got %q", result.Decision.Action)
	}
}
