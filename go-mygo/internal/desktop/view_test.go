package desktop

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/byteowlz/{{project_name}}/internal/cli"
	"github.com/byteowlz/{{project_name}}/internal/config"
	"github.com/byteowlz/{{project_name}}/internal/service"
	"github.com/egoist/mygo/ui"
)

func TestNativeUIAndCLIParity(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	for _, key := range []string{"MYGO_APP_THEME", "MYGO_APP_RADIUS", "MYGO_APP_OMARCHY_FILE"} {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}
	svc := service.New(config.Default())
	theme, err := ResolveTheme(svc, config.Default())
	if err != nil {
		t.Fatal(err)
	}
	app := New(svc, theme)
	tester := ui.NewTester(app.View, 900, 640)
	for _, tc := range []struct {
		button, command, input string
		exit                   int
	}{
		{"Read fixture", "snapshot", "", 0},
		{"Validate snapshot", "validate", `{"revision":"fixture-v1","items":[{"id":"sample-a","label":"Synthetic sample A"},{"id":"sample-b","label":"Synthetic sample B"}]}`, 0},
		{"Test denied apply", "apply", "", 1},
		{"Read fixture", "snapshot", "", 0},
	} {
		if err := tester.Click(tc.button); err != nil {
			t.Fatal(err)
		}
		var out, errOut bytes.Buffer
		exit := cli.Run([]string{tc.command, "--json"}, bytes.NewBufferString(tc.input), &out, &errOut)
		var reply service.Envelope
		if err := json.Unmarshal(out.Bytes(), &reply); err != nil {
			t.Fatal(err)
		}
		if exit != tc.exit || !reflect.DeepEqual(reply, app.Reply) {
			t.Fatalf("UI/CLI parity %s: %d %+v != %+v", tc.command, exit, reply, app.Reply)
		}
		if !tester.HasText(app.Status) {
			t.Fatalf("noncolor status missing: want %q, shown %q", app.Status, tester.Texts())
		}
	}
	tester.Key(0, ui.KeyTab)
	if !tester.Focused("Validate snapshot") {
		t.Fatal("native keyboard focus missing")
	}
	tester.Key(0, ui.KeyEnter)
	if !app.Reply.OK {
		t.Fatal("keyboard read failed")
	}
	if len(tester.Announcements()) == 0 {
		t.Fatal("native status not announced")
	}
	tester.SetSize(480, 640)
	if !tester.HasText("Synthetic sample A") && !tester.HasText("sample-a — Synthetic sample A") {
		t.Fatal("fixture missing at narrow size")
	}
	if dir := os.Getenv("MYGO_HEADLESS_ARTIFACT_DIR"); dir != "" {
		os.MkdirAll(dir, 0700)
		// Headless software rendering is labeled separately from actual native capture.
		writePNG(t, filepath.Join(dir, "native-headless-narrow.png"), tester.Image())
	}
}

func TestNativeModesAndInvalidTheme(t *testing.T) {
	for _, mode := range []string{"dark", "light"} {
		cfg := config.Config{Theme: mode, Radius: 8}
		svc := service.New(cfg)
		theme, err := ResolveTheme(svc, cfg)
		if err != nil {
			t.Fatal(err)
		}
		app := New(svc, theme)
		tester := ui.NewTester(app.View, 900, 640)
		if !tester.HasText("Read fixture") {
			t.Fatal("native view missing")
		}
		if dir := os.Getenv("MYGO_HEADLESS_ARTIFACT_DIR"); dir != "" {
			os.MkdirAll(dir, 0700)
			writePNG(t, filepath.Join(dir, "native-headless-"+mode+".png"), tester.Image())
		}
	}
	cfg := config.Config{Theme: "omarchy", OmarchyFile: filepath.Join(t.TempDir(), "missing")}
	if _, err := ResolveTheme(service.New(cfg), cfg); err == nil {
		t.Fatal("invalid theme silently fell back")
	}
}
