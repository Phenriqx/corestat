package helpers

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func FilterPartitions(fsType string) bool {
	fsType = strings.ToLower(fsType)
	switch fsType {
	case "fuseblk", "tmpfs", "fusectl", "configfs", "tracefs",
		"devpts", "mqueue", "hugetlbfs", "debugfs", "pstorefs",
		"binfmt_misc", "cgroup2fs", "securityfs", "efivarfs", "sysfs", "proc", "devtmpfs":
		return true
	}
	return false
}

func ParseTime(seconds int) string {
	duration := time.Duration(seconds) * time.Second

	hours := int(duration / time.Hour)
	duration -= time.Duration(hours) * time.Hour

	minutes := int(duration / time.Minute)
	duration -= time.Duration(minutes) * time.Minute

	seconds = int(duration / time.Second)

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func GetProcessChildren(p *process.Process) ([]*process.Process, error) {
	children, err := p.Children()
	if err != nil {
		slog.Error("Error loading process children: ", "err", err)
		return nil, err
	}

	return children, nil
}

func ProcessChild(children []*process.Process) (*ChildrenInformation, error) {
	var procs []ChildrenInfo
	
	for _, child := range children {
		name, _ := child.Name()
		cwd, _ := child.Cwd()
		memory, _ := child.MemoryInfo()
		pid, _ := child.Ppid()

		proc := ChildrenInfo {
			Name: name,
			Cwd: cwd,
			MemoryInfo: *memory,
			PID: pid,
		}
		procs = append(procs, proc)
	}
	result := ChildrenInformation(procs)
	if len(result) == 0 {
		return nil, fmt.Errorf("no valid process information found")
	}

	return &result, nil
}