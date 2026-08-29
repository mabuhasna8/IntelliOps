package policy

import (
	"sort"

	"github.com/mabuhasna8/IntelliOps/internal/contracts"
)

func sortRules(rules []contracts.PolicyRule) []contracts.PolicyRule {
	sorted := append([]contracts.PolicyRule(nil), rules...)

	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Priority != sorted[j].Priority {
			return sorted[i].Priority > sorted[j].Priority
		}

		return sorted[i].ID < sorted[j].ID
	})

	return sorted
}
