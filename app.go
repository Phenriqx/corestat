package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Phenriqx/corestat/helpers"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// App struct
type App struct {
	ctx context.Context
}

type HostInformation struct {
}

type CpuInformation struct {
	CPUPercent     []float64
	CPUCores       int
	CPUInfo        []cpu.InfoStat
	CPUFrequency   []float64
	CPUTemperature map[string]string
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
		"Total memory":  float64(virtualMemory.Total),
		"Used memory":   float64(virtualMemory.Used),
		"Cached memory": float64(virtualMemory.Cached),
		"Free memory":   float64(virtualMemory.Free),
		"Percent used":  virtualMemory.UsedPercent,
	}

	return vMem, nil
}

// LINUX ONLY
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

		frequencies = append(frequencies, freqKHz/1000)
	}
	return frequencies, nil
}

func (a *App) GetCPU() (*helpers.CPUInformation, error) {
	interval := time.Duration(1) * time.Second
	usedCpu, err := cpu.Percent(interval, true)
	if err != nil {
		slog.Error("Error loading CPU usage", "err", err)
		return nil, err
	}
	for i := 0; i < len(usedCpu); i++ {
		usedCpu[i] = math.Round(usedCpu[i])
	}

	cpuCores, err := cpu.Counts(true)
	if err != nil {
		slog.Error("Error loading CPU cores", "err", err)
		return nil, err
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		slog.Error("Error loading CPU info", "err", err)
		return nil, err
	}

	cpuFrequency, err := a.FetchDynamicFrequency(cpuCores)
	if err != nil {
		slog.Error("Error loading CPU frequency", "err", err)
		return nil, err
	}

	cpuTemperatures, err := a.FetchTemperature()
	if err != nil {
		slog.Error("Error loading CPU temperature", "err", err)
		return nil, err
	}

	return &helpers.CPUInformation{
		CPUPercent:     usedCpu,
		CPUCores:       cpuCores,
		CPUInfo:        cpuInfo,
		CPUFrequency:   cpuFrequency,
		CPUTemperature: cpuTemperatures,
	}, nil
}

func (a *App) FetchTemperature() (map[string]string, error) {
	temperatures := make(map[string]string)
	cmd := exec.Command("sensors")

	output, err := cmd.Output()
	if err != nil {
		slog.Error("Error fetching temperature", "err", err)
		return make(map[string]string), err
	}

	re := regexp.MustCompile(`(Core \d+|Package id \d+):\s+\+([\d.]+)°C`)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	index := 0
	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); matches != nil {
			core := fmt.Sprintf("%d", index)
			temp := matches[2]
			temperatures[core] = temp
			index++
		}
	}
	return temperatures, nil
}

func (a *App) GetDiskUsage() (map[string]disk.UsageStat, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		slog.Error("Error fetching disk partitions", "err", err)
		return nil, err
	}

	diskInfo := make(map[string]disk.UsageStat)
	for _, partition := range partitions {
		if helpers.FilterPartitions(partition.Fstype) {
			slog.Info("Skipping filtered partition", "fstype", partition.Fstype)
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			slog.Error("Error fetching disk usage",
				"mountpoint", partition.Mountpoint,
				"err", err)
			continue // Skip this partition instead of returning error
		}
		if usage.Total == 0 {
			slog.Info("Skipping partition with total size 0", "mountpoint", partition.Mountpoint)
			continue
		}

		slog.Info("Disk usage", "mountpoint", usage.Path, "total", usage.Total, "free", usage.Free)
		diskInfo[usage.Path] = *usage
	}

	if len(diskInfo) == 0 {
		slog.Error("No valid disk information found")
		return nil, fmt.Errorf("no valid disk information found")
	}

	return diskInfo, nil
}

func (a *App) GetHostInfo() (*helpers.HostInformation, error) {
	majorInfo, err := host.Info()
	if err != nil {
		slog.Error("Error loading host info", "err", err)
		return nil, err
	}
	uptime := helpers.ParseTime(int(majorInfo.Uptime))

	return &helpers.HostInformation{
		MajorInfo: majorInfo,
		Uptime:    uptime,
	}, nil
}

func (a *App) GetProcesses() (*helpers.ProcessInformation, error) {
	procs, err := process.Processes()
	if err != nil {
		slog.Error("Error loading processes: ", "err", err)
		return nil, err
	}

	var infos []helpers.ProcessInfo

	for _, proc := range procs {
		parentProcess, err := proc.Parent()
		if err != nil {
			slog.Error("error getting parent process", "err", err)
			continue
		}

		name, err := parentProcess.Name()
		if err != nil {
			slog.Error("Error loading process name: ", "err", err)
			continue
		}
		cwd, err := parentProcess.Cwd()
		if err != nil {
			slog.Error("Error loading process cwd: ", "err", err)
			continue
		}
		hostUser, err := parentProcess.Username()
		if err != nil {
			slog.Error("Error loading process username: ", "err", err)
			continue
		}
		cpuPercent, err := parentProcess.CPUPercent()
		if err != nil {
			slog.Error("Error loading process CPU percent: ", "err", err)
			continue
		}
		PID, err := parentProcess.Ppid()
		if err != nil {
			slog.Error("Error loading process PID: ", "err", err)
			continue
		}
		memory, err := parentProcess.MemoryInfo()
		if err != nil {
			slog.Error("Error loading process memory info: ", "err", err)
			continue
		}
		threads, err := parentProcess.NumThreads()
		if err != nil {
			slog.Error("Error loading process threads: ", "err", err)
			continue
		}

		procChildren, err := helpers.GetProcessChildren(parentProcess)
		if err != nil {
			slog.Error("Error getting process's children.", "err", err)
			continue
		} else if len(procChildren) == 0 {
			procChildren = nil
		}

		procChild, err := helpers.ProcessChild(procChildren)
		if err != nil {
			slog.Error("Error getting process's children information.", "err", err)
			continue
		}

		info := helpers.ProcessInfo{
			Name:       name,
			Cwd:        cwd,
			HostUser:   hostUser,
			CPUPercent: cpuPercent,
			PID:        PID,
			MemoryInfo: *memory,
			Threads:    threads,
			Children:   *procChild,
		}
		infos = append(infos, info)
	}

	result := helpers.ProcessInformation(infos)
	if len(result) == 0 {
		slog.Error("No valid process information found")
		return nil, fmt.Errorf("no valid process information found")
	}

	return &result, nil
}

func (a *App) SigKillProcess(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		slog.Error("Error creating process", "err", err)
		return fmt.Errorf("error creating process: %v", err)
	}

	result, err := proc.IsRunning()
	if !result {
		slog.Error("Process is not running", "pid", pid, "err", err)
		return fmt.Errorf("process with PID %d is not running", pid)
	}

	if err := proc.Kill(); err != nil {
		slog.Error("Error killing process", "err", err)
		return fmt.Errorf("error killing process: %v", err)
	}

	return nil
}

func (a *App) SigTerminateProcess(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		slog.Error("Error creating process", "err", err)
		return fmt.Errorf("error creating process: %v", err)
	}

	result, err := proc.IsRunning()
	if !result {
		slog.Error("Process is not running", "pid", pid, "err", err)
		return fmt.Errorf("process with PID %d is not running", pid)
	}

	if err := proc.Terminate(); err != nil {
		slog.Error("Error terminating process", "err", err)
		return fmt.Errorf("error terminating process: %v", err)
	}

	return nil
}

func (a *App) GetRemainingBattery() (int, error) {
	path := fmt.Sprintf("/sys/class/power_supply/BAT0/capacity")
	capacity, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Error reading from file", "path", path, "err", err)
		return 1, err
	}
	battery, err := strconv.Atoi(strings.Trim(string(capacity), "\n"))
	if err != nil {
		slog.Error("Error parsing int", "err", err)
		return 1, err
	}

	return battery, nil 
}