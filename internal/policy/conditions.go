package policy

import (
	"reflect"
	"strings"

	"github.com/mabuhasna8/IntelliOps/internal/contracts"
)

func matchesAll(
	conditions []contracts.Condition,
	input map[string]any,
) bool {
	for _, condition := range conditions {
		value, exists := lookup(input, condition.Field)

		if !matches(condition, value, exists) {
			return false
		}
	}

	return true
}

func lookup(input map[string]any, field string) (any, bool) {
	var current any = input

	for _, part := range strings.Split(field, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		current, ok = object[part]
		if !ok {
			return nil, false
		}
	}

	return current, true
}

func matches(
	condition contracts.Condition,
	actual any,
	exists bool,
) bool {
	switch condition.Operator {
	case "exists":
		expected, ok := condition.Value.(bool)
		return ok && exists == expected

	case "equals":
		return exists && reflect.DeepEqual(actual, condition.Value)

	case "not_equals":
		return !exists || !reflect.DeepEqual(actual, condition.Value)

	case "in":
		values, ok := condition.Value.([]any)
		if !ok || !exists {
			return false
		}

		for _, value := range values {
			if reflect.DeepEqual(actual, value) {
				return true
			}
		}

		return false

	default:
		return false
	}
}
