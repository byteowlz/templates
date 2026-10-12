// Native-proof opens a synthetic native window, records content captures and exits.
// Content captures re-render on CPU; use an OS screenshot plus frame stats for GPU evidence.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	design "github.com/byteowlz/design-system/impls/mygo-go"
	"github.com/byteowlz/{{project_name}}/internal/config"
	"github.com/byteowlz/{{project_name}}/internal/desktop"
	"github.com/byteowlz/{{project_name}}/internal/service"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: native-proof OUTPUT_DIR")
		os.Exit(2)
	}
	output, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := os.MkdirAll(output, 0700); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	done := make(chan error, 1)
	mygo.App.WhenReady(func() {
		cfg := config.Default()
		svc := service.New(cfg)
		dark, err := desktop.ResolveTheme(svc, cfg)
		if err != nil {
			done <- err
			mygo.App.Quit()
			return
		}
		cfg.Theme = "light"
		light, err := desktop.ResolveTheme(svc, cfg)
		if err != nil {
			done <- err
			mygo.App.Quit()
			return
		}
		app := desktop.New(svc, dark)
		win := mygo.NewWindow(mygo.WindowOptions{Title: "MyGo Native Fixture Proof", Width: 900, Height: 640, Content: ui.View(app.View)})
		if win.Page() != nil {
			done <- fmt.Errorf("unexpected webview")
			mygo.App.Quit()
			return
		}
		go func() {
			err := capture(win, app, light, output)
			done <- err
			mygo.App.Quit()
		}()
	})
	if err := mygo.App.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	select {
	case err := <-done:
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "native proof closed before completion")
		os.Exit(1)
	}
}

func capture(win *mygo.Window, app *desktop.App, light design.Resolved, output string) error {
	time.Sleep(4 * time.Second)
	if err := capturePNG(win, filepath.Join(output, "native-content-dark.png")); err != nil {
		return err
	}
	updated := make(chan struct{})
	win.Update(func() { app.Design = light; close(updated) })
	select {
	case <-updated:
	case <-time.After(5 * time.Second):
		return fmt.Errorf("native update timed out")
	}
	time.Sleep(2 * time.Second)
	if err := capturePNG(win, filepath.Join(output, "native-content-light.png")); err != nil {
		return err
	}
	metadata := map[string]any{"mode": "native ui.View", "webview": false, "capture": "CapturePage native CPU content re-render", "gpu_claim": "inspect frame-stats log and OS screenshot separately", "fixture": "synthetic read-only"}
	raw, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(output, "metadata.json"), raw, 0600)
}

func capturePNG(win *mygo.Window, path string) error {
	picture, err := win.CapturePage()
	if err != nil {
		return err
	}
	return os.WriteFile(path, picture, 0600)
}
