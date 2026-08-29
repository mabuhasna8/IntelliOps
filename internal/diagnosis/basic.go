package diagnosis

import (
	"errors"
	"strings"

	"github.com/mabuhasna8/IntelliOps/internal/contracts"
)

type BasicDiagnoser struct{}

func (BasicDiagnoser) Diagnose(
	result contracts.AdapterResult,
) (Diagnosis, error) {
	if strings.TrimSpace(result.Classification.Code) == "" {
		return Diagnosis{}, errors.New(
			"adapter result classification code is required",
		)
	}

	severity := "info"

	if value, ok := result.Facts["severity"].(string); ok &&
		strings.TrimSpace(value) != "" {
		severity = value
	}

	return Diagnosis{
		Code:     result.Classification.Code,
		Severity: severity,
		Details:  result.Facts,
	}, nil
}

