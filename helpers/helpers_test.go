package helpers

import (
	"testing"

	"github.com/shirou/gopsutil/v3/process"
)

func TestProcessChildHandlesMissingProcesses(t *testing.T) {
	children, err := ProcessChild([]*process.Process{nil, {Pid: -1}})
	if err != nil {
		t.Fatalf("ProcessChild returned an error: %v", err)
	}
	if children == nil || len(*children) != 0 {
		t.Fatalf("expected an empty child list, got %v", children)
	}
}
