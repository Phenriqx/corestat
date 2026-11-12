package main

import (
	"os"
	"testing"

	"github.com/shirou/gopsutil/v3/process"
)

func TestProcessInfoUsesProcessPID(t *testing.T) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		t.Fatalf("open current process: %v", err)
	}

	info, err := processInfo(proc)
	if err != nil {
		t.Fatalf("read current process: %v", err)
	}
	if info.PID != int32(os.Getpid()) {
		t.Fatalf("got PID %d, want %d", info.PID, os.Getpid())
	}
	if info.Name == "" {
		t.Fatal("expected the current process name")
	}
}
