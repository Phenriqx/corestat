package main

import (
	"context"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// App struct
type App struct {
	ctx context.Context
}

type CpuInformation struct {
	CPUPercent []float64
	CPUCores   int
	CPUInfo    []cpu.InfoStat
	CPUFrequency []float64
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetRAM() (map[string]float64, error) {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		slog.Error("Error loading virtual memory", "err", err)
		return nil, err
	}

	vMem := map[string]float64{
		"Total memory": float64(virtualMemory.Total),
		"Used memory": float64(virtualMemory.Used),
		"Cached memory": float64(virtualMemory.Cached),
		"Free memory":  float64(virtualMemory.Free),
		"Percent used": virtualMemory.UsedPercent,
	}

	return vMem, nil
}

func (a *App) FetchDynamicFrequency(cores int) ([]float64, error) {
	frequencies := make([]float64, 0, cores)
	for i := 0; i < cores; i++ {
		freqPath := filepath.Join("/sys/devices/system/cpu", "cpu"+strconv.Itoa(i), "cpufreq", "scaling_cur_freq")
		freq, err := os.ReadFile(freqPath)
		if err != nil {
			frequencies = append(frequencies, 0)
			continue
		}
		
		freqKHz, err := strconv.ParseFloat(strings.TrimSpace(string(freq)), 64)
		if err != nil {
			frequencies = append(frequencies, 0)
			continue
		}

		frequencies = append(frequencies, freqKHz / 1000)
	}
	return frequencies, nil
}

func (a *App) GetCPU() (*CpuInformation, error) {
	interval := time.Duration(1) * time.Second
	usedCpu, err := cpu.Percent(interval, true)
	if err != nil {
		slog.Error("Error loading CPU usage", "err", err)
		return &CpuInformation{}, err
	}
	for i := 0; i < len(usedCpu); i++ {
		usedCpu[i] = math.Round(usedCpu[i])
	}

	cpuCores, err := cpu.Counts(true)
	if err != nil {
		slog.Error("Error loading CPU cores", "err", err)
		return &CpuInformation{}, err
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		slog.Error("Error loading CPU info", "err", err)
		return &CpuInformation{}, err
	}

	cpuFrequency, err := a.FetchDynamicFrequency(cpuCores)
	if err != nil {
		slog.Error("Error loading CPU frequency", "err", err)
		return &CpuInformation{}, err
	}

	return &CpuInformation{
		CPUPercent: usedCpu,
		CPUCores:   cpuCores,
		CPUInfo:   cpuInfo,
		CPUFrequency: cpuFrequency,
	}, nil
}