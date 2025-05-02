package helpers

import (
	"fmt"
	"strings"
	"time"
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
