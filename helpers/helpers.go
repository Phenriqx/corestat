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
	procs := make(ChildrenInformation, 0, len(children))

	for _, child := range children {
		if child == nil {
			continue
		}
		memory, err := child.MemoryInfo()
		if err != nil || memory == nil {
			continue // The process may have exited or become inaccessible.
		}
		name, err := child.Name()
		if err != nil {
			continue
		}
		cwd, _ := child.Cwd() // Some processes do not expose their working directory.
		procs = append(procs, ChildrenInfo{
			Name:       name,
			Cwd:        cwd,
			MemoryInfo: *memory,
			PID:        child.Pid,
		})
	}
	return &procs, nil
}
