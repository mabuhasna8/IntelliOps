package diagnosis

import "github.com/mabuhasna8/IntelliOps/internal/contracts"

type Diagnosis struct {
	Code     string         `json:"code"`
	Severity string         `json:"severity"`
	Details  map[string]any `json:"details,omitempty"`
}

type Diagnoser interface {
	Diagnose(contracts.AdapterResult) (Diagnosis, error)
}
