package simulation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mabuhasna8/IntelliOps/apps/agent/executor"
)

type Outcome struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Metadata map[string]any
}

type Executor struct {
	scenarios map[string]Outcome
}

func New(scenarios map[string]Outcome) *Executor {
	copied := make(map[string]Outcome, len(scenarios))

	for name, outcome := range scenarios {
		copied[name] = outcome
	}

	return &Executor{
		scenarios: copied,
	}
}

func (e *Executor) Execute(
	_ context.Context,
	request executor.Request,
) (executor.Result, error) {
	if request.ExecutionID == "" {
		return executor.Result{}, errors.New("execution ID is required")
	}

	if request.Scenario == "" {
		return executor.Result{}, errors.New("scenario is required")
	}

	outcome, ok := e.scenarios[request.Scenario]
	if !ok {
		return executor.Result{}, fmt.Errorf(
			"scenario %q not found",
			request.Scenario,
		)
	}

	startedAt := time.Now().UTC()
	finishedAt := time.Now().UTC()

	result := executor.Result{
		ExecutionID: request.ExecutionID,
		Command:     request.Command,
		ExitCode:    outcome.ExitCode,
		Stdout:      outcome.Stdout,
		Stderr:      outcome.Stderr,
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
		Duration:    finishedAt.Sub(startedAt),
		Simulated:   true,
		Scenario:    request.Scenario,
		Metadata:    cloneMetadata(outcome.Metadata),
	}

	if outcome.ExitCode != 0 {
		return result, fmt.Errorf(
			"simulated command failed with exit code %d",
			outcome.ExitCode,
		)
	}

	return result, nil
}

func cloneMetadata(input map[string]any) map[string]any {
	if input == nil {
		return nil
	}

	output := make(map[string]any, len(input))

	for key, value := range input {
		output[key] = value
	}

	return output
}
