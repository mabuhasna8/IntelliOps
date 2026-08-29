package adapter

import (
	"fmt"

	"github.com/mabuhasna8/IntelliOps/internal/contracts"
)

func Normalize(
	envelope contracts.ExecutionEnvelope,
) (contracts.AdapterResult, error) {
	if envelope.ExecutionID == "" {
		return contracts.AdapterResult{}, fmt.Errorf(
			"execution ID is required",
		)
	}

	status := contracts.StatusPassed

	if envelope.Termination.TimedOut {
		status = contracts.StatusTimedOut
	} else if envelope.Termination.Cancelled {
		status = contracts.StatusCancelled
	} else if envelope.Termination.ExitCode == nil {
		status = contracts.StatusUnknown
	} else if *envelope.Termination.ExitCode != 0 {
		status = contracts.StatusFailed
	}

	classificationCode := "execution.passed"
	if status == contracts.StatusFailed {
		classificationCode = "execution.failed"
	} else if status == contracts.StatusTimedOut {
		classificationCode = "execution.timed_out"
	} else if status == contracts.StatusCancelled {
		classificationCode = "execution.cancelled"
	} else if status == contracts.StatusUnknown {
		classificationCode = "execution.unknown"
	}

	facts := map[string]any{
		"execution_id": envelope.ExecutionID,
		"status":       string(status),
		"mode":         string(envelope.Mode),
		"synthetic":    envelope.Synthetic,
		"stdout_bytes": envelope.Stdout.TotalBytes,
		"stderr_bytes": envelope.Stderr.TotalBytes,
	}

	if envelope.Termination.ExitCode != nil {
		facts["exit_code"] = *envelope.Termination.ExitCode
	}

	return contracts.AdapterResult{
		SchemaVersion: contracts.AdapterResultSchemaV1,
		Status:        status,
		Classification: contracts.Classification{
			Code: classificationCode,
		},
		Facts: facts,
	}, nil
}

