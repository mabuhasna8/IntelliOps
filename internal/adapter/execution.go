package adapter

import (
	"github.com/mabuhasna8/IntelliOps/apps/agent/executor"
	"github.com/mabuhasna8/IntelliOps/internal/contracts"
)

func ToExecutionEnvelope(result executor.Result) contracts.ExecutionEnvelope {
	exitCode := result.ExitCode

	mode := contracts.ExecutionModeLive
	if result.Simulated {
		mode = contracts.ExecutionModeSimulation
	}

	return contracts.ExecutionEnvelope{
		SchemaVersion: contracts.ExecutionEnvelopeSchemaV1,
		ExecutionID:   result.ExecutionID,
		Mode:          mode,
		Synthetic:     result.Simulated,
		StartedAt:     result.StartedAt,
		FinishedAt:    result.FinishedAt,

		Termination: contracts.TerminationMetadata{
			State:    contracts.TerminationExited,
			ExitCode: &exitCode,
		},

		Stdout: contracts.LogStream{
			TotalBytes: int64(len(result.Stdout)),
		},

		Stderr: contracts.LogStream{
			TotalBytes: int64(len(result.Stderr)),
		},
	}
}
