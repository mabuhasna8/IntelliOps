package contracts

import "time"

type ExecutorType string

const (
	ExecutorDirectHost ExecutorType = "direct_host"
	ExecutorSimulation ExecutorType = "simulation"
)

type ExecutionRequest struct {
	SchemaVersion string `json:"schema_version"`

	ExecutionID string `json:"execution_id"`
	AttemptID   string `json:"attempt_id"`

	Mode ExecutionMode `json:"mode"`

	ActionID      string `json:"action_id"`
	ActionVersion int    `json:"action_version"`

	ExecutorType    ExecutorType `json:"executor_type"`
	ExecutorVersion string       `json:"executor_version,omitempty"`

	Command     []string          `json:"command"`
	Environment map[string]string `json:"environment,omitempty"`

	Workspace WorkspacePolicy `json:"workspace"`
	Limits    ResourceLimits  `json:"limits"`

	Cancellation CancellationPolicy `json:"cancellation"`
}

type WorkspacePolicy struct {
	RootPath         string `json:"root_path,omitempty"`
	CleanupOnSuccess bool   `json:"cleanup_on_success"`
	CleanupOnFailure bool   `json:"cleanup_on_failure"`
}

type ResourceLimits struct {
	TimeoutSeconds   int64 `json:"timeout_seconds,omitempty"`
	MaxStdoutBytes   int64 `json:"max_stdout_bytes,omitempty"`
	MaxStderrBytes   int64 `json:"max_stderr_bytes,omitempty"`
	MaxArtifactBytes int64 `json:"max_artifact_bytes,omitempty"`
	MaxArtifactCount int   `json:"max_artifact_count,omitempty"`
}

type CancellationPolicy struct {
	Supported          bool  `json:"supported"`
	GracePeriodSeconds int64 `json:"grace_period_seconds,omitempty"`
	ForceKill          bool  `json:"force_kill"`
}

type TerminationState string

const (
	TerminationExited     TerminationState = "exited"
	TerminationSignaled   TerminationState = "signaled"
	TerminationTimedOut   TerminationState = "timed_out"
	TerminationCancelled  TerminationState = "cancelled"
	TerminationStartError TerminationState = "start_error"
	TerminationUnknown    TerminationState = "unknown"
)

type TerminationMetadata struct {
	State     TerminationState `json:"state"`
	ExitCode  *int             `json:"exit_code,omitempty"`
	Signal    *string          `json:"signal,omitempty"`
	TimedOut  bool             `json:"timed_out"`
	Cancelled bool             `json:"cancelled"`
	ErrorCode string           `json:"error_code,omitempty"`
	Error     string           `json:"error,omitempty"`
}

type Artifact struct {
	Name      string  `json:"name"`
	BlobRef   BlobRef `json:"blob_ref"`
	SizeBytes int64   `json:"size_bytes"`
	SHA256    string  `json:"sha256,omitempty"`
}

type ExecutionEnvelope struct {
	SchemaVersion string `json:"schema_version"`

	ExecutionID string `json:"execution_id"`
	AttemptID   string `json:"attempt_id"`

	Mode      ExecutionMode `json:"mode"`
	Synthetic bool          `json:"synthetic"`

	Executor struct {
		Type    ExecutorType `json:"type"`
		Version string       `json:"version"`
		AgentID string       `json:"agent_id"`
	} `json:"executor"`

	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`

	Termination TerminationMetadata `json:"termination"`

	Stdout LogStream `json:"stdout"`
	Stderr LogStream `json:"stderr"`

	Artifacts []Artifact `json:"artifacts,omitempty"`

	Workspace struct {
		ID       string `json:"id"`
		Retained bool   `json:"retained"`
	} `json:"workspace"`
}
