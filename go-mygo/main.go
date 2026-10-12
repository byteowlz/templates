package main

import (
	"fmt"
	"os"

	"github.com/byteowlz/{{project_name}}/internal/cli"
	"github.com/byteowlz/{{project_name}}/internal/config"
	"github.com/byteowlz/{{project_name}}/internal/desktop"
	"github.com/byteowlz/{{project_name}}/internal/service"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	options, err := cli.DesktopOptions(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	var startupErr error
	mygo.App.WhenReady(func() {
		cfg, err := config.Load(options, true)
		if err != nil {
			startupErr = err
			mygo.App.Quit()
			return
		}
		svc := service.New(cfg)
		theme, err := desktop.ResolveTheme(svc, cfg)
		if err != nil {
			startupErr = err
			mygo.App.Quit()
			return
		}
		app := desktop.New(svc, theme)
		mygo.NewWindow(mygo.WindowOptions{
			Title: "{{project_name}}", Width: 900, Height: 640, MinWidth: 480, MinHeight: 420,
			Content: ui.View(app.View), // Page() is nil: no hidden webview.
		})
	})
	if err := mygo.App.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if startupErr != nil {
		fmt.Fprintln(os.Stderr, "startup:", startupErr)
		os.Exit(1)
	}
}
