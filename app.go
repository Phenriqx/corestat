package main

import (
	"context"
	"log/slog"

	"github.com/shirou/gopsutil/v3/mem"
)

// App struct
type App struct {
	ctx context.Context
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

func (a *App) GetCPU() {
	
}