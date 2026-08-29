package pipeline

import (
	"context"
	"errors"
	"testing"

	"github.com/mabuhasna8/IntelliOps/apps/agent/executor"
	"github.com/mabuhasna8/IntelliOps/internal/contracts"
	"github.com/mabuhasna8/IntelliOps/internal/diagnosis"
)

type fakeDiagnoser struct {
	result diagnosis.Diagnosis
	err    error
	called bool
}

func (f *fakeDiagnoser) Diagnose(
	contracts.AdapterResult,
) (diagnosis.Diagnosis, error) {
	f.called = true
	return f.result, f.err
}

type fakeEvaluator struct {
	result contracts.DecisionEvaluation
	called bool
}

func (f *fakeEvaluator) Evaluate(
	contracts.DecisionPolicy,
	map[string]any,
) contracts.DecisionEvaluation {
	f.called = true
	return f.result
}

func TestPipelineProcess(t *testing.T) {
	t.Parallel()

	diagnoser := &fakeDiagnoser{
		result: diagnosis.Diagnosis{
			Code:     "execution.passed",
			Severity: "info",
			Details: map[string]any{
				"source": "test",
			},
		},
	}

	evaluator := &fakeEvaluator{}

	p := Pipeline{
		Diagnoser: diagnoser,
		Evaluator: evaluator,
		Policy: contracts.DecisionPolicy{
			ID:      "test-policy",
			Version: 1,
		},
	}

	result, err := p.Process(
		context.Background(),
		executor.Result{
			ExecutionID: "execution-1",
		},
	)
	if err != nil {
		t.Fatalf("Process returned an error: %v", err)
	}

	if result.Envelope.ExecutionID != "execution-1" {
		t.Fatalf(
			"expected execution ID %q, got %q",
			"execution-1",
			result.Envelope.ExecutionID,
		)
	}

	if result.Adapter.Classification.Code != "execution.passed" {
		t.Fatalf(
			"expected classification code %q, got %q",
			"execution.passed",
			result.Adapter.Classification.Code,
		)
	}

	if result.Diagnosis.Code != "execution.passed" {
		t.Fatalf(
			"expected diagnosis code %q, got %q",
			"execution.passed",
			result.Diagnosis.Code,
		)
	}

	if !diagnoser.called {
		t.Fatal("expected diagnoser to be called")
	}

	if !evaluator.called {
		t.Fatal("expected evaluator to be called")
	}
}

func TestPipelineRequiresDiagnoser(t *testing.T) {
	t.Parallel()

	p := Pipeline{
		Evaluator: &fakeEvaluator{},
	}

	_, err := p.Process(
		context.Background(),
		executor.Result{},
	)
	if err == nil {
		t.Fatal("expected an error")
	}

	if err.Error() != "pipeline: diagnoser is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPipelineRequiresEvaluator(t *testing.T) {
	t.Parallel()

	p := Pipeline{
		Diagnoser: &fakeDiagnoser{},
	}

	_, err := p.Process(
		context.Background(),
		executor.Result{},
	)
	if err == nil {
		t.Fatal("expected an error")
	}

	if err.Error() != "pipeline: evaluator is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPipelinePropagatesDiagnosisError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("diagnosis failed")

	p := Pipeline{
		Diagnoser: &fakeDiagnoser{
			err: expectedErr,
		},
		Evaluator: &fakeEvaluator{},
	}

	_, err := p.Process(
		context.Background(),
		executor.Result{
			ExecutionID: "execution-1",
		},
	)
	if err == nil {
		t.Fatal("expected an error")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected diagnosis error, got: %v", err)
	}
}

func TestPipelineEvaluatorReceivesDiagnosis(t *testing.T) {
	t.Parallel()

	diagnoser := &fakeDiagnoser{
		result: diagnosis.Diagnosis{
			Code:     "execution.failed",
			Severity: "error",
			Details: map[string]any{
				"reason": "non-zero exit code",
			},
		},
	}

	evaluator := &fakeEvaluator{}

	p := Pipeline{
		Diagnoser: diagnoser,
		Evaluator: evaluator,
	}

	_, err := p.Process(
		context.Background(),
		executor.Result{
			ExecutionID: "execution-1",
		},
	)
	if err != nil {
		t.Fatalf("Process returned an error: %v", err)
	}

	if !evaluator.called {
		t.Fatal("expected evaluator to be called")
	}
}
