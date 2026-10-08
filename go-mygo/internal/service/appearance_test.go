package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/byteowlz/{{project_name}}/internal/config"
)

func TestThemeDataOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "colors.toml")
	for _, text := range []string{"broken = [", "background=123", strings.Repeat("x", 65537)} {
		os.WriteFile(path, []byte(text), 0600)
		reply := New(config.Config{Theme: "omarchy", OmarchyFile: path}).Execute(context.Background(), Request{Operation: "appearance"})
		if reply.OK || reply.Error.Code != "INVALID_THEME" {
			t.Fatal(reply)
		}
	}
	payload := "background = \"#123456\"\n"
	os.WriteFile(path, []byte(payload), 0600)
	reply := New(config.Config{Theme: "omarchy", OmarchyFile: path}).Execute(context.Background(), Request{Operation: "appearance"})
	if !reply.OK || reply.Result.Appearance.Colors["background"] != "#123456" {
		t.Fatal(reply)
	}
	after, _ := os.ReadFile(path)
	if string(after) != payload {
		t.Fatal("theme data changed")
	}
	missing := New(config.Config{Theme: "omarchy", OmarchyFile: filepath.Join(dir, "missing")}).Execute(context.Background(), Request{Operation: "appearance"})
	if missing.OK || missing.Error.Code != "THEME_UNAVAILABLE" {
		t.Fatal(missing)
	}
}
