package helpers

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
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
