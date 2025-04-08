package main

import (
	"embed"

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

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "System Monitor",
		Width:  1024,
		Height: 768,
		WindowStartState: options.Maximised,
		MinWidth: 800,
		MinHeight: 600,
		Frameless: false,
		StartHidden: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 10, G: 10, B: 10, A: 175},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		OnBeforeClose: app.BeforeClose,
		EnableFraudulentWebsiteDetection: true,

		Windows: &windows.Options{
			WindowIsTranslucent: true,
			DisableWindowIcon: false,
			IsZoomControlEnabled: true,
			ZoomFactor: 1.0,
			DisableFramelessWindowDecorations: false,
			Theme: windows.SystemDefault,
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
				TabFocusesLinks: mac.Enabled,
				FullscreenEnabled: mac.Enabled,				
			},
		},
		})
	if err != nil {
		println("Error:", err.Error())
	}
}
