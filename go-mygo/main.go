package main

import (
	"fmt"
	"os"

	"github.com/byteowlz/{{project_name}}/internal/cli"
	"github.com/byteowlz/{{project_name}}/internal/config"
	"github.com/byteowlz/{{project_name}}/internal/desktop"
	"github.com/byteowlz/{{project_name}}/internal/service"
	"github.com/egoist/mygo"
)

// Build/generate bind without config I/O. Startup I/O belongs in WhenReady.
func main() {
	options, err := cli.DesktopOptions(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	svc := service.New(config.Default())
	mygo.Bind(desktop.Desktop{Service: svc})
	mygo.App.WhenReady(func() {
		cfg, err := config.Load(options, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, "configuration:", err)
			mygo.App.Quit()
			return
		}
		svc.Configure(cfg)
		mygo.NewWindow(mygo.WindowOptions{
			Title: "{{project_name}}", Width: 900, Height: 640,
			MinWidth: 480, MinHeight: 360, URL: "/",
		})
	})
	if err := mygo.App.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
