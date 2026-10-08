package desktop

import (
	"context"
	"encoding/json"
	"fmt"

	design "github.com/byteowlz/design-system/impls/mygo-go"
	"github.com/byteowlz/{{project_name}}/internal/service"
	"github.com/egoist/mygo/ui"
)

// App keeps model data only. ui.Element values never escape their build pass.
type App struct {
	Service *service.Service
	Design  design.Resolved
	Reply   service.Envelope
	Status  string
}

func New(svc *service.Service, theme design.Resolved) *App {
	a := &App{Service: svc, Design: theme}
	a.perform("snapshot", "")
	return a
}

func (a *App) perform(operation, input string) {
	a.Reply = a.Service.Inspect(context.Background(), service.Request{Operation: operation, Input: input})
	if !a.Reply.OK {
		a.Status = a.Reply.Error.Code + ": " + a.Reply.Error.Message
		return
	}
	a.Status = "Read-only fixture loaded. No external device or writes."
	if operation == "validate" {
		a.Status = "Snapshot validated. No writes performed."
	}
}

// View is wholly native GPU UI; no HTML, JavaScript, IPC or webview exists.
func (a *App) View(c *ui.Context) {
	theme := a.Design.Theme
	// Preserve OS text-scale preference without changing canonical color slots.
	theme.FontSize *= c.Preferences().TextScale
	c.SetTheme(&theme)
	ui.Column(c).Fill().Padding(24).Gap(16).Children(func() {
		ui.Text(c, "{{project_name}}").FontSize(24).Bold().Role(ui.RoleHeading).Level(1)
		ui.Text(c, "Illustrative native starter — not an approved product composition.").TextColor(a.Design.Roles[design.MutedForeground])
		ui.Row(c).Wrap().Gap(12).Children(func() {
			if ui.Button(c.Key("read"), "Read fixture").MinHeight(44).Clicked() {
				a.perform("snapshot", "")
				c.Announce(a.Status)
			}
			if ui.Button(c.Key("validate"), "Validate snapshot").MinHeight(44).Disabled(!a.Reply.OK).Clicked() {
				data, err := json.Marshal(a.Reply.Result.Snapshot)
				if err != nil {
					a.Status = "Snapshot cannot be encoded."
				} else {
					a.perform("validate", string(data))
				}
				c.Announce(a.Status)
			}
			if ui.Button(c.Key("denied"), "Test denied apply").MinHeight(44).Clicked() {
				a.perform("apply", "")
				c.Announce(a.Status)
			}
		})
		ui.Text(c, a.Status)
		ui.Column(c).Padding(16).Gap(8).Radius(a.Design.Radius.LG).Background(a.Design.Roles[design.SurfaceSunken]).Children(func() {
			if a.Reply.OK && a.Reply.Result.Snapshot != nil {
				ui.Text(c, "Revision: "+a.Reply.Result.Snapshot.Revision)
				for _, item := range a.Reply.Result.Snapshot.Items {
					ui.Text(c.Key(item.ID), fmt.Sprintf("%s — %s", item.ID, item.Label))
				}
			} else {
				ui.Text(c, "No snapshot displayed. Read fixture to recover.")
			}
		})
		ui.Text(c, "Headless equivalent: {{project_name}}ctl snapshot --json").TextColor(a.Design.Roles[design.MutedForeground])
	})
}
