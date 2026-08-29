package contracts

type WorkflowStatus string

const (
	WorkflowDraft                WorkflowStatus = "draft"
	WorkflowValidated            WorkflowStatus = "validated"
	WorkflowTested               WorkflowStatus = "tested"
	WorkflowSubmittedForApproval WorkflowStatus = "submitted_for_approval"
	WorkflowApproved             WorkflowStatus = "approved"
	WorkflowPublished            WorkflowStatus = "published"
	WorkflowRejected             WorkflowStatus = "rejected"
)

type NodeType string

const (
	NodeStart        NodeType = "start"
	NodeAction       NodeType = "action"
	NodeDecision     NodeType = "decision"
	NodeManualReview NodeType = "manual_review"
	NodeNotification NodeType = "notification"
	NodeTerminal     NodeType = "terminal"
)

type WorkflowVersion struct {
	WorkflowID string `json:"workflow_id"`
	Version    int    `json:"version"`

	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	Status WorkflowStatus `json:"status"`

	StartNodeID string         `json:"start_node_id"`
	Nodes       []WorkflowNode `json:"nodes"`
	Transitions []Transition   `json:"transitions"`

	CreatedBy string `json:"created_by"`
}

type WorkflowNode struct {
	ID   string   `json:"id"`
	Type NodeType `json:"type"`

	Name string `json:"name,omitempty"`

	ActionID      string `json:"action_id,omitempty"`
	ActionVersion int    `json:"action_version,omitempty"`

	AdapterID      string `json:"adapter_id,omitempty"`
	AdapterVersion int    `json:"adapter_version,omitempty"`

	PolicyID      string `json:"policy_id,omitempty"`
	PolicyVersion int    `json:"policy_version,omitempty"`

	Capabilities map[string]any `json:"capabilities,omitempty"`
}

type Transition struct {
	FromNodeID    string `json:"from_node_id"`
	Outcome       string `json:"outcome"`
	ToNodeID      string `json:"to_node_id"`
	TransitionKey string `json:"transition_key,omitempty"`
}
