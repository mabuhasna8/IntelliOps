package contracts

type AdapterInput struct {
	SchemaVersion     string            `json:"schema_version"`
	ExecutionEnvelope ExecutionEnvelope `json:"execution_envelope"`
}

type ResultStatus string

const (
	StatusPassed    ResultStatus = "passed"
	StatusFailed    ResultStatus = "failed"
	StatusTimedOut  ResultStatus = "timed_out"
	StatusCancelled ResultStatus = "cancelled"
	StatusUnknown   ResultStatus = "unknown"
)

type Classification struct {
	Code       string  `json:"code"`
	Confidence float64 `json:"confidence,omitempty"`
}

type AdapterResult struct {
	SchemaVersion string `json:"schema_version"`

	Status         ResultStatus   `json:"status"`
	Classification Classification `json:"classification"`
	Facts          map[string]any `json:"facts,omitempty"`

	EvidenceRefs []BlobRef `json:"evidence_refs,omitempty"`
}

type Diagnosis struct {
	DiagnosisID string `json:"diagnosis_id"`

	ExecutionID string `json:"execution_id"`
	AttemptID   string `json:"attempt_id"`
	EnvelopeID  string `json:"envelope_id"`

	AdapterID      string `json:"adapter_id"`
	AdapterVersion int    `json:"adapter_version"`

	Result AdapterResult `json:"result"`
}
