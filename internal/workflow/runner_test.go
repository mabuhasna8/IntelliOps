package workflow

import (
	"context"
	"errors"
	"testing"

	"github.com/mabuhasna8/IntelliOps/apps/agent/executor"
	"github.com/mabuhasna8/IntelliOps/internal/contracts"
	"github.com/mabuhasna8/IntelliOps/internal/diagnosis"
)

type fakeExecutor struct {
	result executor.Result
	err    error
	calls  int
}

func (f *fakeExecutor) Execute(
	_ context.Context,
	_ executor.Request,
) (executor.Result, error) {
	f.calls++
	return f.result, f.err
}

type fakeExecutorFactory struct {
	executor executor.Executor
	nodes    []string
}

func (f *fakeExecutorFactory) ExecutorFor(
	node contracts.WorkflowNode,
) (executor.Executor, error) {
	f.nodes = append(f.nodes, node.ID)
	return f.executor, nil
}

type fakePolicyResolver struct {
	policy contracts.DecisionPolicy
	nodes  []string
}

func (f *fakePolicyResolver) Resolve(
	_ context.Context,
	policyID string,
	_ int,
) (contracts.DecisionPolicy, error) {
	f.nodes = append(f.nodes, policyID)
	return f.policy, nil
}

type fakeDiagnoser struct {
	result diagnosis.Diagnosis
	err    error
}

func (f *fakeDiagnoser) Diagnose(
	contracts.AdapterResult,
) (diagnosis.Diagnosis, error) {
	return f.result, f.err
}

func testWorkflow() contracts.WorkflowVersion {
	return contracts.WorkflowVersion{
		ID:         "workflow-version-1",
		WorkflowID: "workflow-1",
		Version:    1,
		EntryNode:  "build",
		Nodes: []contracts.WorkflowNode{
			{
				ID:            "build",
				Type:          contracts.NodeTypeAction,
				PolicyID:      "build-policy",
				PolicyVersion: 1,
				Action: &contracts.ActionSpec{
					Command: "build",
				},
			},
			{
				ID:            "test",
				Type:          contracts.NodeTypeAction,
				PolicyID:      "test-policy",
				PolicyVersion: 1,
				Action: &contracts.ActionSpec{
					Command: "test",
				},
			},
		},
		Edges: []contracts.WorkflowEdge{
			{
				From: "build",
				To:   "test",
			},
		},
	}
}

func testRunner() (*Runner, *fakeExecutorFactory, *fakePolicyResolver) {
	exec := &fakeExecutor{
		result: executor.Result{
			ExecutionID: "execution-1",
			ExitCode:    0,
		},
	}

	executors := &fakeExecutorFactory{
		executor: exec,
	}

	policies := &fakePolicyResolver{
		policy: contracts.DecisionPolicy{
			ID:      "test-policy",
			Version: 1,
			Default: contracts.Decision{
				Action:  "continue",
				Outcome: contracts.OutcomeContinue,
			},
		},
	}

	runner := &Runner{
		Executors: executors,
		Policies:  policies,
		Diagnoser: diagnosis.BasicDiagnoser{},
	}

	return runner, executors, policies
}

func TestRunnerExecutionSequence(t *testing.T) {
	runner, executors, policies := testRunner()

	result, err := runner.Run(
		context.Background(),
		testWorkflow(),
		"execution-1",
	)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	wantNodes := []string{"build", "test"}

	if len(result.Nodes) != len(wantNodes) {
		t.Fatalf("Nodes = %v, want %v", result.Nodes, wantNodes)
	}

	for i, want := range wantNodes {
		if result.Nodes[i] != want {
			t.Fatalf("Nodes = %v, want %v", result.Nodes, wantNodes)
		}
	}

	if len(result.Decisions) != 2 {
		t.Fatalf("expected two decisions, got %d", len(result.Decisions))
	}

	if result.Decisions[0].NodeID != "build" {
		t.Fatalf("first decision belongs to %q, want %q",
			result.Decisions[0].NodeID, "build")
	}

	if result.Decisions[1].NodeID != "test" {
		t.Fatalf("second decision belongs to %q, want %q",
			result.Decisions[1].NodeID, "test")
	}

	if len(executors.nodes) != 2 {
		t.Fatalf("executor factory called for %d nodes, want 2",
			len(executors.nodes))
	}

	if len(policies.nodes) != 2 {
		t.Fatalf("policy resolver called %d times, want 2",
			len(policies.nodes))
	}
}

func TestRunnerRequiresDependencies(t *testing.T) {
	version := testWorkflow()

	tests := []struct {
		name string
		r    Runner
		want string
	}{
		{
			name: "executor factory",
			r: Runner{
				Policies:  &fakePolicyResolver{},
				Diagnoser: diagnosis.BasicDiagnoser{},
			},
			want: "workflow runner: executor factory is required",
		},
		{
			name: "policy resolver",
			r: Runner{
				Executors: &fakeExecutorFactory{},
				Diagnoser: diagnosis.BasicDiagnoser{},
			},
			want: "workflow runner: policy resolver is required",
		},
		{
			name: "diagnoser",
			r: Runner{
				Executors: &fakeExecutorFactory{},
				Policies:  &fakePolicyResolver{},
			},
			want: "workflow runner: diagnoser is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.Run(
				context.Background(),
				version,
				"execution-1",
			)
			if err == nil {
				t.Fatal("Run returned nil error")
			}

			if err.Error() != tt.want {
				t.Fatalf("error = %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

func TestRunnerDetectsCycle(t *testing.T) {
	runner, _, _ := testRunner()
	version := testWorkflow()

	version.Edges = []contracts.WorkflowEdge{
		{From: "build", To: "test"},
		{From: "test", To: "build"},
	}

	result, err := runner.Run(
		context.Background(),
		version,
		"execution-1",
	)
	if err == nil {
		t.Fatal("expected cycle detection error")
	}

	if !errors.Is(err, err) {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Nodes) != 2 {
		t.Fatalf("Nodes = %v, want [build test]", result.Nodes)
	}
}

func TestRunnerRejectsUnsupportedNodeType(t *testing.T) {
	runner, _, _ := testRunner()
	version := testWorkflow()

	version.Nodes[0].Type = contracts.NodeTypeManual

	_, err := runner.Run(
		context.Background(),
		version,
		"execution-1",
	)
	if err == nil {
		t.Fatal("expected unsupported node type error")
	}
}
