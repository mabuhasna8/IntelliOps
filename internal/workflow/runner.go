package workflow

import (
	"context"
	"errors"
	"fmt"

	"github.com/mabuhasna8/IntelliOps/apps/agent/executor"
	"github.com/mabuhasna8/IntelliOps/internal/contracts"
	"github.com/mabuhasna8/IntelliOps/internal/diagnosis"
	"github.com/mabuhasna8/IntelliOps/internal/pipeline"
	"github.com/mabuhasna8/IntelliOps/internal/policy"
)

type ExecutorFactory interface {
	ExecutorFor(node contracts.WorkflowNode) (executor.Executor, error)
}

type PolicyResolver interface {
	Resolve(
		ctx context.Context,
		policyID string,
		version int,
	) (contracts.DecisionPolicy, error)
}

type Runner struct {
	Executors ExecutorFactory
	Policies  PolicyResolver
	Diagnoser diagnosis.Diagnoser
}

type RunResult struct {
	Decisions []contracts.Decision
	Nodes     []string
}

func (r Runner) Run(
	ctx context.Context,
	version contracts.WorkflowVersion,
	executionID string,
) (RunResult, error) {
	if r.Executors == nil {
		return RunResult{}, errors.New("workflow runner: executor factory is required")
	}

	if r.Policies == nil {
		return RunResult{}, errors.New("workflow runner: policy resolver is required")
	}

	if r.Diagnoser == nil {
		return RunResult{}, errors.New("workflow runner: diagnoser is required")
	}

	if err := Validate(version); err != nil {
		return RunResult{}, err
	}

	nodes := make(map[string]contracts.WorkflowNode, len(version.Nodes))
	for _, node := range version.Nodes {
		nodes[node.ID] = node
	}

	current := version.EntryNode
	result := RunResult{}

	visited := make(map[string]bool)

	for current != "" {
		if visited[current] {
			return result, fmt.Errorf(
				"workflow runner: cycle detected at node %q",
				current,
			)
		}
		visited[current] = true

		node := nodes[current]
		result.Nodes = append(result.Nodes, node.ID)

		if node.Type != contracts.NodeTypeAction {
			return result, fmt.Errorf(
				"workflow runner: unsupported node type %q at node %q",
				node.Type,
				node.ID,
			)
		}

		if node.Action == nil {
			return result, fmt.Errorf(
				"workflow runner: action missing at node %q",
				node.ID,
			)
		}

		exec, err := r.Executors.ExecutorFor(node)
		if err != nil {
			return result, fmt.Errorf(
				"workflow runner: create executor for %q: %w",
				node.ID,
				err,
			)
		}

		request := executor.Request{
			ExecutionID: executionID,
			Command:     node.Action.Command,
			Args:        node.Action.Args,
			Env:         node.Action.Environment,
			Workspace:   node.Action.WorkingDirectory,
		}

		executionResult, executionErr := exec.Execute(ctx, request)

		decisionPolicy, err := r.Policies.Resolve(
			ctx,
			node.PolicyID,
			node.PolicyVersion,
		)
		if err != nil {
			return result, fmt.Errorf(
				"workflow runner: resolve policy for %q: %w",
				node.ID,
				err,
			)
		}

		p := pipeline.Pipeline{
			Diagnoser: r.Diagnoser,
			Evaluator: policy.NewEvaluator(),
			Policy:    decisionPolicy,
		}

		processed, err := p.Process(ctx, executionResult)
		if err != nil {
			return result, fmt.Errorf(
				"workflow runner: process node %q: %w",
				node.ID,
				err,
			)
		}

		processed.Decision.NodeID = node.ID
		result.Decisions = append(result.Decisions, processed.Decision)

		if executionErr != nil &&
			processed.Decision.Outcome == contracts.OutcomeContinue {
			return result, fmt.Errorf(
				"workflow runner: node %q failed: %w",
				node.ID,
				executionErr,
			)
		}

		next, ok := nextNode(version.Edges, node.ID, processed.Decision)
		if !ok {
			break
		}

		current = next
	}

	return result, nil
}

func nextNode(
	edges []contracts.WorkflowEdge,
	from string,
	decision contracts.Decision,
) (string, bool) {
	for _, edge := range edges {
		if edge.From != from {
			continue
		}

		if edge.When == "" ||
			edge.When == decision.Outcome.String() ||
			edge.When == decision.Action {
			return edge.To, true
		}
	}

	return "", false
}
