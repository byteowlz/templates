package app

import (
	"path/filepath"
	"testing"
)

func TestResolveBaseFollowsXDGSpec(t *testing.T) {
	home := filepath.FromSlash("/home/ada")
	custom, _ := filepath.Abs(filepath.FromSlash("/custom/config"))
	cases := []struct {
		name      string
		xdg       string
		isWindows bool
		want      string
	}{
		{"absolute value is used", custom, false, custom},
		{"unset uses default", "", false, filepath.Join(home, ".config")},
		{"relative value is ignored", "relative/config", false, filepath.Join(home, ".config")},
		{"tilde is not expanded", "~/config", false, filepath.Join(home, ".config")},
		{"absolute value wins on windows", custom, true, custom},
		{"windows falls back to known folder", "relative", true, filepath.FromSlash("/appdata")},
	}
	for _, tc := range cases {
		got, err := resolveBase(tc.xdg, home, filepath.FromSlash("/appdata"), tc.isWindows, ".config")
		if err != nil || got != tc.want {
			t.Errorf("%s: got %q, %v; want %q", tc.name, got, err, tc.want)
		}
	}
	if _, err := resolveBase("", "", "", false, ".config"); err == nil {
		t.Error("missing home must be an error")
	}
}

func TestDefaultDirsIgnoreInvalidEnvironment(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_DATA_HOME", "relative")
	t.Setenv("XDG_STATE_HOME", "")
	data, err := defaultDataDir("demo")
	if err != nil || data != filepath.Join(home, ".local", "share", "demo") {
		t.Errorf("data dir: got %q, %v", data, err)
	}
	state, err := defaultStateDir("demo")
	if err != nil || state != filepath.Join(home, ".local", "state", "demo") {
		t.Errorf("state dir: got %q, %v", state, err)
	}
}
