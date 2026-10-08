package service

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

func (s *Service) appearance() Envelope {
	s.mu.RLock()
	cfg := s.cfg
	s.mu.RUnlock()
	appearance := &Appearance{Theme: cfg.Theme, Radius: cfg.Radius, Colors: map[string]string{}}
	if cfg.Theme != "omarchy" {
		return success(Result{Appearance: appearance})
	}
	path := cfg.OmarchyFile
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return Fail("THEME_UNAVAILABLE", "cannot locate Omarchy data")
		}
		path = filepath.Join(home, ".local", "state", "omarchy", "current", "theme", "colors.toml")
	}
	file, err := os.Open(path)
	if err != nil {
		return Fail("THEME_UNAVAILABLE", "cannot read Omarchy colors.toml")
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil || !stat.Mode().IsRegular() || stat.Size() > 65536 {
		return Fail("INVALID_THEME", "theme must be a regular data file under 64 KiB")
	}
	data, err := io.ReadAll(io.LimitReader(file, 65537))
	if err != nil || len(data) > 65536 {
		return Fail("INVALID_THEME", "theme exceeds 64 KiB or could not be read")
	}
	if err := toml.Unmarshal(data, &appearance.Colors); err != nil {
		return Fail("INVALID_THEME", "expected flat TOML string data")
	}
	return success(Result{Appearance: appearance})
}
