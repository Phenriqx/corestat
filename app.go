package main

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"log/slog"
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetRAM() (map[string]float64, error) {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		slog.Error("Error loading virtual memory", "err", err)
		return nil, err
	}

	fmt.Printf("Total memory: %v\nFree memory: %v\n, Used memory: %v", virtualMemory.Total, virtualMemory.Free, virtualMemory.UsedPercent)
	vMem := map[string]float64{
		"Total memory": float64(virtualMemory.Total),
		"Used memory":  float64(virtualMemory.Free),
		"Percent used": virtualMemory.UsedPercent,
	}

	return vMem, nil
}
