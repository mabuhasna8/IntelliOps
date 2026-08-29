package adapter

import "github.com/mabuhasna8/IntelliOps/internal/contracts"

type Adapter interface {
	Adapt(contracts.ExecutionEnvelope) (contracts.AdapterResult, error)
}
