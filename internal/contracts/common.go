package contracts

import "time"

const (
    ExecutionEnvelopeSchemaV1 = "execution-envelope.v1"
    AdapterInputSchemaV1      = "adapter-input.v1"
    AdapterResultSchemaV1     = "adapter-result.v1"
    DecisionSchemaV1          = "decision.v1"
    SimulationTraceSchemaV1   = "simulation-trace.v1"
    AuditEventSchemaV1        = "audit-event.v1"
)

type Timestamp = time.Time

type BlobRef struct {
    URI        string `json:"uri"`
    SizeBytes  int64  `json:"size_bytes,omitempty"`
    SHA256     string `json:"sha256,omitempty"`
    MediaType  string `json:"media_type,omitempty"`
}

type LogStream struct {
    BlobRefs       []BlobRef `json:"blob_refs,omitempty"`
    TotalBytes     int64     `json:"total_bytes"`
    PersistedBytes int64     `json:"persisted_bytes"`
    Truncated      bool      `json:"truncated"`
}

type ExecutionMode string

const (
    ExecutionModeLive       ExecutionMode = "live"
    ExecutionModeSimulation ExecutionMode = "simulation"
)
