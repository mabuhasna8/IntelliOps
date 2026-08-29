package contracts

import (
	"encoding/json"
	"testing"
)

func TestExecutionEnvelopeRoundTrip(t *testing.T) {
	exitCode := 0

	original := ExecutionEnvelope{
		SchemaVersion: ExecutionEnvelopeSchemaV1,
		ExecutionID:   "exec-1",
		AttemptID:     "attempt-1",
		Mode:          ExecutionModeLive,
		Termination: TerminationMetadata{
			State:    TerminationExited,
			ExitCode: &exitCode,
		},
	}

	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}

	var decoded ExecutionEnvelope
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}

	if decoded.ExecutionID != original.ExecutionID {
		t.Fatalf("execution ID mismatch: got %q", decoded.ExecutionID)
	}

	if decoded.Termination.ExitCode == nil {
		t.Fatal("expected exit code")
	}

	if *decoded.Termination.ExitCode != 0 {
		t.Fatalf("unexpected exit code: %d", *decoded.Termination.ExitCode)
	}
}
