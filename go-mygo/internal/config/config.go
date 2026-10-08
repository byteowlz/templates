package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const AppName = "{{project_name}}"

// Config is the portable identity/data configuration, not domain authority.
type Config struct {
	Theme       string `json:"theme" toml:"theme"`
	Radius      int    `json:"radius" toml:"radius"`
	OmarchyFile string `json:"omarchy_file" toml:"omarchy_file"`
}

type Options struct {
	Path        string
	Theme       *string
	Radius      *int
	OmarchyFile *string
}

func Default() Config { return Config{Theme: "dark"} }

func GlobalPath() (string, error) {
	root := os.Getenv("XDG_CONFIG_HOME")
	if root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		root = filepath.Join(home, ".config")
	}
	if !filepath.IsAbs(root) {
		return "", errors.New("XDG_CONFIG_HOME must be absolute")
	}
	return filepath.Join(root, AppName, "config.toml"), nil
}

// Load overlays defaults/global/local/environment/explicit-file/flags in order.
// Only desktop first-run passes initialize=true; read-only CLI audits never write.
func Load(o Options, initialize bool) (Config, error) {
	cfg := Default()
	global, err := GlobalPath()
	if err != nil {
		return cfg, err
	}
	found := false
	for _, path := range []string{global, "config.toml"} {
		exists, err := overlay(&cfg, path, false)
		if err != nil {
			return cfg, err
		}
		found = found || exists
	}
	if err := environment(&cfg); err != nil {
		return cfg, err
	}
	if o.Path != "" {
		if _, err := overlay(&cfg, o.Path, true); err != nil {
			return cfg, err
		}
		found = true
	}
	if o.Theme != nil {
		cfg.Theme = *o.Theme
	}
	if o.Radius != nil {
		cfg.Radius = *o.Radius
	}
	if o.OmarchyFile != nil {
		cfg.OmarchyFile = *o.OmarchyFile
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	if initialize && !found {
		if err := Init(global); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}

func overlay(cfg *Config, path string, required bool) (bool, error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) && !required {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("config: %w", err)
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, 65537))
	if err != nil {
		return false, err
	}
	if len(data) > 65536 {
		return false, errors.New("config exceeds 64 KiB")
	}
	if err := toml.NewDecoder(strings.NewReader(string(data))).DisallowUnknownFields().Decode(cfg); err != nil {
		return false, fmt.Errorf("invalid TOML config: %w", err)
	}
	return true, nil
}

func environment(cfg *Config) error {
	if value, ok := os.LookupEnv("MYGO_APP_THEME"); ok {
		cfg.Theme = value
	}
	if value, ok := os.LookupEnv("MYGO_APP_RADIUS"); ok {
		radius, err := strconv.Atoi(value)
		if err != nil {
			return errors.New("MYGO_APP_RADIUS must be an integer")
		}
		cfg.Radius = radius
	}
	if value, ok := os.LookupEnv("MYGO_APP_OMARCHY_FILE"); ok {
		cfg.OmarchyFile = value
	}
	return nil
}

func (c Config) Validate() error {
	if c.Theme != "dark" && c.Theme != "light" && c.Theme != "omarchy" {
		return errors.New("theme must be dark, light or omarchy")
	}
	if c.Radius < 0 || c.Radius > 32 {
		return errors.New("radius must be 0..32 pixels")
	}
	return nil
}

// Init exclusively creates a private default file; never overwrites existing data.
func Init(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if errors.Is(err, os.ErrExist) {
		return nil
	}
	if err != nil {
		return err
	}
	_, writeErr := f.WriteString("theme = \"dark\"\nradius = 0\nomarchy_file = \"\"\n")
	closeErr := f.Close()
	return errors.Join(writeErr, closeErr)
}
