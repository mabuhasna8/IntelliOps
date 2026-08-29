package direct_host

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mabuhasna8/IntelliOps/apps/agent/executor"
)

type Executor struct {
	AllowedCommands map[string]bool
}

func New(allowedCommands []string) *Executor {
	allowed := make(map[string]bool, len(allowedCommands))

	for _, command := range allowedCommands {
		allowed[command] = true
	}

	return &Executor{
		AllowedCommands: allowed,
	}
}

func (e *Executor) Execute(
	ctx context.Context,
	request executor.Request,
) (executor.Result, error) {
	if request.ExecutionID == "" {
		return executor.Result{}, errors.New("execution ID is required")
	}

	if request.Command == "" {
		return executor.Result{}, errors.New("command is required")
	}

	if !e.AllowedCommands[request.Command] {
		return executor.Result{}, fmt.Errorf(
			"command %q is not allowed",
			request.Command,
		)
	}

	workspace, err := prepareWorkspace(request.Workspace)
	if err != nil {
		return executor.Result{}, err
	}

	execContext := ctx

	var cancel context.CancelFunc
	if request.Timeout > 0 {
		execContext, cancel = context.WithTimeout(ctx, request.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(
		execContext,
		request.Command,
		request.Args...,
	)

	cmd.Dir = workspace
	cmd.Env = buildEnvironment(request.Env)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startedAt := time.Now().UTC()
	err = cmd.Run()
	finishedAt := time.Now().UTC()

	result := executor.Result{
		ExecutionID: request.ExecutionID,
		Command:     request.Command,
		ExitCode:    exitCode(cmd, err),
		Stdout:      stdout.String(),
		Stderr:      stderr.String(),
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
		Duration:    finishedAt.Sub(startedAt),
	}

	if err != nil {
		// A non-zero exit status is execution evidence, not necessarily
		// an executor failure. Return both the result and the error.
		return result, err
	}

	return result, nil
}

func prepareWorkspace(path string) (string, error) {
	if path == "" {
		return os.MkdirTemp("", "intelliops-execution-")
	}

	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve workspace: %w", err)
	}

	if err := os.MkdirAll(absolute, 0o750); err != nil {
		return "", fmt.Errorf("create workspace: %w", err)
	}

	return absolute, nil
}

func buildEnvironment(values map[string]string) []string {
	environment := os.Environ()

	for key, value := range values {
		environment = append(environment, key+"="+value)
	}

	return environment
}

func exitCode(cmd *exec.Cmd, runErr error) int {
	if runErr == nil {
		return 0
	}

	if cmd.ProcessState == nil {
		return -1
	}

	return cmd.ProcessState.ExitCode()
}
