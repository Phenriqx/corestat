package helpers

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/process"
)

type HostInformation struct {
	MajorInfo *host.InfoStat
	Uptime    string
}

type CPUInformation struct {
	CPUPercent     []float64
	CPUCores       int
	CPUInfo        []cpu.InfoStat
	CPUFrequency   []float64
	CPUTemperature map[string]string
}

type ProcessInfo struct {
	Name       string
	Cwd        string
	HostUser   string
	CPUPercent float64
	PID        int32
	MemoryInfo process.MemoryInfoStat
	Threads    int32
}

type ProcessInformation []ProcessInfo
