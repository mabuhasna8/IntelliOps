package executor

import (
	"context"
	"time"
)

type Request struct {
	ExecutionID string
	Command     string
	Args        []string
	Env         map[string]string
	Workspace   string
	Timeout     time.Duration

	Mode     string
	Scenario string
}

type Result struct {
	ExecutionID string
	Command     string
	ExitCode    int
	Stdout      string
	Stderr      string
	StartedAt   time.Time
	FinishedAt  time.Time
	Duration    time.Duration

	Simulated bool
	Scenario  string
	Metadata  map[string]any
}

type Executor interface {
	Execute(context.Context, Request) (Result, error)
}

const (
	ModeLive       = "live"
	ModeSimulation = "simulation"
)
