package adapter

import (
	"testing"
	"time"

	"github.com/mabuhasna8/IntelliOps/apps/agent/executor"
)

func TestToExecutionEnvelope_MapsExecutorResult(t *testing.T) {
	startedAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	input := executor.Result{
		ExecutionID: "exec-123",
		Command:     "echo",
		ExitCode:    0,
		Stdout:      "hello stdout",
		Stderr:      "hello stderr",
		StartedAt:   startedAt,
		FinishedAt:  startedAt.Add(2 * time.Second),
		Duration:    2 * time.Second,
		Simulated:   true,
		Scenario:    "success",
		Metadata: map[string]any{
			"host": "test-host",
		},
	}

	envelope := ToExecutionEnvelope(input)

	if envelope.ExecutionID != input.ExecutionID {
		t.Fatalf(
			"ExecutionID = %q, want %q",
			envelope.ExecutionID,
			input.ExecutionID,
		)
	}

	if !envelope.StartedAt.Equal(input.StartedAt) {
		t.Fatalf(
			"StartedAt = %v, want %v",
			envelope.StartedAt,
			input.StartedAt,
		)
	}
}
