package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrecedenceAndReadOnly(t *testing.T) {
	t.Chdir(t.TempDir())
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("MYGO_APP_THEME", "light")
	t.Setenv("MYGO_APP_RADIUS", "7")
	t.Setenv("MYGO_APP_OMARCHY_FILE", "environment")
	path, _ := GlobalPath()
	if err := Init(path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("theme='dark'\nradius=1\nomarchy_file='global'"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("config.toml", []byte("radius=2\nomarchy_file='local'"), 0600); err != nil {
		t.Fatal(err)
	}
	explicit := filepath.Join(t.TempDir(), "explicit.toml")
	if err := os.WriteFile(explicit, []byte("radius=3\n"), 0600); err != nil {
		t.Fatal(err)
	}
	c, err := Load(Options{Path: explicit}, false)
	if err != nil || c.Radius != 3 || c.Theme != "light" || c.OmarchyFile != "environment" {
		t.Fatalf("%+v %v", c, err)
	}
	radius := 4
	theme := "dark"
	file := "flag"
	c, err = Load(Options{Path: explicit, Radius: &radius, Theme: &theme, OmarchyFile: &file}, false)
	if err != nil || c.Radius != 4 || c.Theme != "dark" || c.OmarchyFile != "flag" {
		t.Fatalf("%+v %v", c, err)
	}
	// Remove environment to prove local overrides global, then global fallback.
	t.Setenv("MYGO_APP_RADIUS", "")
	os.Unsetenv("MYGO_APP_RADIUS")
	c, err = Load(Options{}, false)
	if err != nil || c.Radius != 2 {
		t.Fatalf("%+v %v", c, err)
	}
	os.Remove("config.toml")
	c, err = Load(Options{}, false)
	if err != nil || c.Radius != 1 {
		t.Fatalf("%+v %v", c, err)
	}
	before, _ := os.ReadFile(path)
	if err := Init(path); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(path)
	if string(before) != string(after) {
		t.Fatal("init overwrote config")
	}
}

func TestFirstRunAndInvalidConfig(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	for _, key := range []string{"MYGO_APP_THEME", "MYGO_APP_RADIUS", "MYGO_APP_OMARCHY_FILE"} {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}
	path, _ := GlobalPath()
	if _, err := Load(Options{}, false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("read-only load wrote config")
	}
	if _, err := Load(Options{}, true); err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(path)
	if err != nil || st.Mode().Perm() != 0600 {
		t.Fatalf("private default: %v %v", st, err)
	}
	if err := os.WriteFile("config.toml", []byte("unknown=1"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(Options{}, false); err == nil {
		t.Fatal("accepted unknown key")
	}
	os.Remove("config.toml")
	if _, err := Load(Options{Path: "missing.toml"}, false); err == nil {
		t.Fatal("ignored missing explicit config")
	}
	t.Setenv("MYGO_APP_RADIUS", "NaN")
	if _, err := Load(Options{}, false); err == nil {
		t.Fatal("accepted invalid env")
	}
}
