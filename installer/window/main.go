package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// uerr := win.CheckAndElevate()
	// if uerr != nil {
	// 	fmt.Printf("Error: %v\n", uerr)
	// 	return
	// }
	app := NewApp()
	err := wails.Run(&options.App{
		Title:         Conf.AppName,
		Width:         Conf.Width,
		Height:        Conf.Height,
		DisableResize: Conf.Max,
		Frameless:     true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			BackdropType:         windows.Mica,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
