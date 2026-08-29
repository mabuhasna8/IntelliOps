package workflow

import (
	"fmt"
	"sort"

	"github.com/mabuhasna8/IntelliOps/internal/contracts"
)

type ValidationError struct {
	Problems []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("workflow validation failed: %v", e.Problems)
}

func Validate(version contracts.WorkflowVersion) error {
	var problems []string

	if version.ID == "" {
		problems = append(problems, "workflow version ID is required")
	}

	if version.WorkflowID == "" {
		problems = append(problems, "workflow ID is required")
	}

	if version.Version <= 0 {
		problems = append(problems, "version must be greater than zero")
	}

	if version.EntryNode == "" {
		problems = append(problems, "entry node is required")
	}

	nodes := make(map[string]contracts.WorkflowNode, len(version.Nodes))

	for _, node := range version.Nodes {
		if node.ID == "" {
			problems = append(problems, "node ID is required")
			continue
		}

		if _, exists := nodes[node.ID]; exists {
			problems = append(problems, "duplicate node ID: "+node.ID)
			continue
		}

		if !validNodeType(node.Type) {
			problems = append(problems,
				fmt.Sprintf("node %q has invalid type %q", node.ID, node.Type))
		}

		if node.Type == contracts.NodeTypeAction {
			if node.Action == nil {
				problems = append(problems,
					fmt.Sprintf("action node %q requires an action", node.ID))
			} else if node.Action.Command == "" {
				problems = append(problems,
					fmt.Sprintf("action node %q requires a command", node.ID))
			}
		}

		if node.Type == contracts.NodeTypeDecision &&
			node.PolicyID == "" {
			problems = append(problems,
				fmt.Sprintf("decision node %q requires a policy ID", node.ID))
		}

		nodes[node.ID] = node
	}

	if version.EntryNode != "" {
		if _, exists := nodes[version.EntryNode]; !exists {
			problems = append(problems,
				"entry node does not exist: "+version.EntryNode)
		}
	}

	edgeKeys := make(map[string]struct{}, len(version.Edges))

	for _, edge := range version.Edges {
		if _, exists := nodes[edge.From]; !exists {
			problems = append(problems,
				"edge source does not exist: "+edge.From)
		}

		if _, exists := nodes[edge.To]; !exists {
			problems = append(problems,
				"edge destination does not exist: "+edge.To)
		}

		key := edge.From + "\x00" + edge.To + "\x00" + edge.When
		if _, exists := edgeKeys[key]; exists {
			problems = append(problems,
				fmt.Sprintf("duplicate edge: %s -> %s", edge.From, edge.To))
		}
		edgeKeys[key] = struct{}{}
	}

	if len(problems) == 0 && !reachableFromEntry(version.EntryNode, nodes, version.Edges) {
		problems = append(problems, "workflow contains unreachable nodes")
	}

	sort.Strings(problems)

	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}

	return nil
}

func validNodeType(nodeType string) bool {
	switch nodeType {
	case contracts.NodeTypeAction,
		contracts.NodeTypeDecision,
		contracts.NodeTypeManual:
		return true
	default:
		return false
	}
}

func reachableFromEntry(
	entry string,
	nodes map[string]contracts.WorkflowNode,
	edges []contracts.WorkflowEdge,
) bool {
	reachable := map[string]bool{entry: true}

	changed := true
	for changed {
		changed = false

		for _, edge := range edges {
			if reachable[edge.From] && !reachable[edge.To] {
				reachable[edge.To] = true
				changed = true
			}
		}
	}

	for nodeID := range nodes {
		if !reachable[nodeID] {
			return false
		}
	}

	return true
}
