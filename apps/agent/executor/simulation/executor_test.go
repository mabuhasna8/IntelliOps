package simulation

import (
	"context"
	"testing"

	"github.com/mabuhasna8/IntelliOps/apps/agent/executor"
)

func TestSimulationReturnsScenarioOutcome(t *testing.T) {
	sim := New(map[string]Outcome{
		"compile-error": {
			ExitCode: 1,
			Stderr:   "undefined: buildTarget",
			Metadata: map[string]any{
				"expected_classification": "build.compilation_error",
			},
		},
	})

	result, err := sim.Execute(context.Background(), executor.Request{
		ExecutionID: "simulation-1",
		Command:     "build",
		Scenario:    "compile-error",
	})

	if err == nil {
		t.Fatal("expected simulated failure")
	}

	if !result.Simulated {
		t.Fatal("expected simulated result")
	}

	if result.ExitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", result.ExitCode)
	}

	if result.Scenario != "compile-error" {
		t.Fatalf("unexpected scenario: %s", result.Scenario)
	}
}

func TestSimulationDoesNotExecuteCommand(t *testing.T) {
	sim := New(map[string]Outcome{
		"success": {
			ExitCode: 0,
			Stdout:   "synthetic output",
		},
	})

	result, err := sim.Execute(context.Background(), executor.Request{
		ExecutionID: "simulation-2",
		Command:     "this-command-must-not-run",
		Scenario:    "success",
	})

	if err != nil {
		t.Fatalf("execute simulation: %v", err)
	}

	if result.Stdout != "synthetic output" {
		t.Fatalf("unexpected output: %q", result.Stdout)
	}
}
