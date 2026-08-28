package contracts

type Scenario struct {
    ScenarioID string `json:"scenario_id"`
    Version    int    `json:"version"`
    Name       string `json:"name"`

    ActionOutcomes map[string][]SyntheticOutcome `json:"action_outcomes"`
}

type SyntheticOutcome struct {
    Status         ResultStatus  `json:"status"`
    Classification Classification `json:"classification"`
    Facts          map[string]any `json:"facts,omitempty"`
}

type Simulation struct {
    SimulationID string `json:"simulation_id"`

    WorkflowID      string `json:"workflow_id"`
    WorkflowVersion int    `json:"workflow_version"`

    ScenarioID      string `json:"scenario_id"`
    ScenarioVersion int    `json:"scenario_version"`

    Status string `json:"status"`
    Result string `json:"result,omitempty"`

    Trace []TraceEntry `json:"trace,omitempty"`
}

type TraceEntry struct {
    Sequence int `json:"sequence"`

    NodeID   string   `json:"node_id"`
    NodeType NodeType `json:"node_type"`

    Attempt int `json:"attempt,omitempty"`

    Status             ResultStatus `json:"status,omitempty"`
    ClassificationCode string       `json:"classification_code,omitempty"`

    DecisionOutcome string `json:"decision_outcome,omitempty"`
    MatchedRule     string `json:"matched_rule,omitempty"`

    SideEffectSuppressed bool `json:"side_effect_suppressed,omitempty"`
}
