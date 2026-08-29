package contracts

import "testing"

func TestDecisionOutcomeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		outcome DecisionOutcome
		want    string
	}{
		{
			name:    "continue",
			outcome: OutcomeContinue,
			want:    string(OutcomeContinue),
		},
		{
			name:    "stop",
			outcome: OutcomeStop,
			want:    string(OutcomeStop),
		},
		{
			name:    "empty outcome",
			outcome: DecisionOutcome(""),
			want:    "",
		},
		{
			name:    "custom outcome",
			outcome: DecisionOutcome("retry"),
			want:    "retry",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.outcome.String(); got != tt.want {
				t.Fatalf(
					"DecisionOutcome(%q).String() = %q, want %q",
					tt.outcome,
					got,
					tt.want,
				)
			}
		})
	}
}
