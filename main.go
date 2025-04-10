package main

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Logging Configuration
	if envErr := godotenv.Load(); envErr != nil {
		fmt.Printf("Error loading .env file: %v\n", envErr)
		return 
	}

	LOG_FILE := os.Getenv("LOG_FILE")
	handlerOptions := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	file, fileErr := os.OpenFile(LOG_FILE, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if fileErr != nil {
		fmt.Printf("Error opening log file: %v\n", fileErr)
		return 
	}

	defer file.Close()
	logger := slog.New(slog.NewJSONHandler(file, &handlerOptions))
	slog.SetDefault(logger)

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "System Monitor",
		Width:            1024,
		Height:           768,
		WindowStartState: options.Maximised,
		MinWidth:         80,
		MinHeight:        60,
		Frameless:        false,
		StartHidden:      false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 10, G: 10, B: 10, A: 175},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		OnBeforeClose:                    app.BeforeClose,
		EnableFraudulentWebsiteDetection: true,

		Windows: &windows.Options{
			WindowIsTranslucent:               true,
			DisableWindowIcon:                 false,
			IsZoomControlEnabled:              true,
			ZoomFactor:                        1.0,
			DisableFramelessWindowDecorations: false,
			Theme:                             windows.SystemDefault,
		},

		Linux: &linux.Options{
			WindowIsTranslucent: true,
		},

		Mac: &mac.Options{
			WindowIsTranslucent: true,
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
			},
			Appearance: mac.NSAppearanceNameDarkAqua,
			Preferences: &mac.Preferences{
				TabFocusesLinks:   mac.Enabled,
				FullscreenEnabled: mac.Enabled,
			},
		},
	})
	if err != nil {
		slog.Error("Error running application", "err", err)
		println("Error:", err.Error())
	}
}