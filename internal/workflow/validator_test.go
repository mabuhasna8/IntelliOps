package workflow

import (
	"strings"
	"testing"

	"github.com/mabuhasna8/IntelliOps/internal/contracts"
)

func TestValidateWorkflow(t *testing.T) {
	version := contracts.WorkflowVersion{
		ID:         "wf-v1",
		WorkflowID: "wf-1",
		Version:    1,
		EntryNode:  "build",
		Nodes: []contracts.WorkflowNode{
			{
				ID:   "build",
				Type: contracts.NodeTypeAction,
				Action: &contracts.ActionSpec{
					Command: "go",
				},
			},
			{
				ID:            "decide",
				Type:          contracts.NodeTypeDecision,
				PolicyID:      "build-policy",
				PolicyVersion: 1,
			},
		},
		Edges: []contracts.WorkflowEdge{
			{From: "build", To: "decide"},
		},
	}

	if err := Validate(version); err != nil {
		t.Fatalf("expected valid workflow, got %v", err)
	}
}

func TestValidateRejectsUnknownEdgeNode(t *testing.T) {
	version := contracts.WorkflowVersion{
		ID:         "wf-v1",
		WorkflowID: "wf-1",
		Version:    1,
		EntryNode:  "start",
		Nodes: []contracts.WorkflowNode{
			{ID: "start", Type: contracts.NodeTypeManual},
		},
		Edges: []contracts.WorkflowEdge{
			{From: "start", To: "missing"},
		},
	}

	err := Validate(version)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "edge destination does not exist") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRejectsUnreachableNode(t *testing.T) {
	version := contracts.WorkflowVersion{
		ID:         "wf-v1",
		WorkflowID: "wf-1",
		Version:    1,
		EntryNode:  "start",
		Nodes: []contracts.WorkflowNode{
			{ID: "start", Type: contracts.NodeTypeManual},
			{ID: "orphan", Type: contracts.NodeTypeManual},
		},
	}

	err := Validate(version)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "unreachable") {
		t.Fatalf("unexpected error: %v", err)
	}
}
