package pipeline

import (
	"context"
	"errors"
	"fmt"

	"github.com/mabuhasna8/IntelliOps/apps/agent/executor"
	"github.com/mabuhasna8/IntelliOps/internal/adapter"
	"github.com/mabuhasna8/IntelliOps/internal/contracts"
	"github.com/mabuhasna8/IntelliOps/internal/diagnosis"
)

type DecisionEvaluator interface {
	Evaluate(
		contracts.DecisionPolicy,
		map[string]any,
	) contracts.DecisionEvaluation
}

type Pipeline struct {
	Diagnoser diagnosis.Diagnoser
	Evaluator DecisionEvaluator
	Policy    contracts.DecisionPolicy
}

type Result struct {
	Envelope  contracts.ExecutionEnvelope
	Adapter   contracts.AdapterResult
	Diagnosis diagnosis.Diagnosis
	Decision  contracts.Decision
}

func (p Pipeline) Process(
	_ context.Context,
	execution executor.Result,
) (Result, error) {
	if p.Diagnoser == nil {
		return Result{}, errors.New("pipeline: diagnoser is required")
	}

	if p.Evaluator == nil {
		return Result{}, errors.New("pipeline: evaluator is required")
	}

	envelope := adapter.ToExecutionEnvelope(execution)

	adapterResult, err := adapter.Normalize(envelope)
	if err != nil {
		return Result{}, fmt.Errorf("pipeline: normalize execution: %w", err)
	}

	diagnosisResult, err := p.Diagnoser.Diagnose(adapterResult)
	if err != nil {
		return Result{}, fmt.Errorf("pipeline: diagnose execution: %w", err)
	}

	facts := map[string]any{
		"status":         adapterResult.Status,
		"classification": adapterResult.Classification,
		"facts":          adapterResult.Facts,
		"diagnosis":      diagnosisResult,
	}

	evaluation := p.Evaluator.Evaluate(p.Policy, facts)

	return Result{
		Envelope:  envelope,
		Adapter:   adapterResult,
		Diagnosis: diagnosisResult,
		Decision:  evaluation.Decision,
	}, nil
}
