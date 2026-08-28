package contracts

type DecisionOutcome string

const (
    OutcomeContinue     DecisionOutcome = "continue"
    OutcomeRetry        DecisionOutcome = "retry"
    OutcomeStop         DecisionOutcome = "stop"
    OutcomeManualReview DecisionOutcome = "manual_review"
    OutcomeNotify       DecisionOutcome = "notify"
    OutcomeCancel       DecisionOutcome = "cancel"
    OutcomeSkip         DecisionOutcome = "skip"
)

type Decision struct {
    SchemaVersion string `json:"schema_version"`

    DecisionID  string `json:"decision_id"`
    ExecutionID string `json:"execution_id,omitempty"`
    SimulationID string `json:"simulation_id,omitempty"`

    NodeID string `json:"node_id"`

    PolicyID      string `json:"policy_id"`
    PolicyVersion int    `json:"policy_version"`

    Outcome      DecisionOutcome `json:"outcome"`
    TransitionKey string         `json:"transition_key,omitempty"`

    ReasonCode  string `json:"reason_code,omitempty"`
    MatchedRule string `json:"matched_rule,omitempty"`

    Retry *RetryDecision `json:"retry,omitempty"`
}

type RetryDecision struct {
    Allowed       bool `json:"allowed"`
    AttemptNumber int  `json:"attempt_number,omitempty"`
    MaxAttempts   int `json:"max_attempts,omitempty"`
}
