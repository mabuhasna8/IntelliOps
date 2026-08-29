package direct_host

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/mabuhasna8/IntelliOps/apps/agent/executor"
)

func TestExecutorRunsAllowedCommand(t *testing.T) {
	command := "printf"

	if runtime.GOOS == "windows" {
		t.Skip("printf test is not portable to Windows")
	}

	exec := New([]string{command})

	result, err := exec.Execute(context.Background(), executor.Request{
		ExecutionID: "execution-1",
		Command:     command,
		Args:        []string{"hello"},
		Timeout:     5 * time.Second,
	})

	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", result.ExitCode)
	}

	if result.Stdout != "hello" {
		t.Fatalf("expected stdout %q, got %q", "hello", result.Stdout)
	}
}

func TestExecutorRejectsDisallowedCommand(t *testing.T) {
	exec := New([]string{"printf"})

	_, err := exec.Execute(context.Background(), executor.Request{
		ExecutionID: "execution-2",
		Command:     "sh",
	})

	if err == nil {
		t.Fatal("expected disallowed command error")
	}
}

func TestExecutorReturnsNonZeroExitCode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell test is not portable to Windows")
	}

	exec := New([]string{"sh"})

	result, err := exec.Execute(context.Background(), executor.Request{
		ExecutionID: "execution-3",
		Command:     "sh",
		Args:        []string{"-c", "echo failed >&2; exit 7"},
	})

	if err == nil {
		t.Fatal("expected command failure")
	}

	if result.ExitCode != 7 {
		t.Fatalf("expected exit code 7, got %d", result.ExitCode)
	}

	if result.Stderr == "" {
		t.Fatal("expected stderr evidence")
	}
}

func TestExecutorHonorsTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("timeout test is not portable to Windows")
	}

	exec := New([]string{"sleep"})

	ctx := context.Background()

	result, err := exec.Execute(ctx, executor.Request{
		ExecutionID: "execution-4",
		Command:     "sleep",
		Args:        []string{"5"},
		Timeout:     10 * time.Millisecond,
	})

	if err == nil {
		t.Fatal("expected timeout error")
	}

	if result.ExitCode == 0 {
		t.Fatal("expected non-zero exit code after timeout")
	}
}
