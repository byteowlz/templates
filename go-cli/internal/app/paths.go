package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// AppPaths captures the resolved filesystem locations used by the CLI.
type AppPaths struct {
	ConfigFile string
	DataDir    string
	StateDir   string
}

// DiscoverPaths determines the config, data, and state directories for the application.
func DiscoverPaths(app string, override string) (AppPaths, error) {
	configFile, err := resolveConfigFile(app, override)
	if err != nil {
		return AppPaths{}, err
	}

	dataDir, err := defaultDataDir(app)
	if err != nil {
		return AppPaths{}, err
	}

	stateDir, err := defaultStateDir(app)
	if err != nil {
		return AppPaths{}, err
	}

	return AppPaths{
		ConfigFile: configFile,
		DataDir:    dataDir,
		StateDir:   stateDir,
	}, nil
}

// ApplyPathOverrides applies overrides from the loaded config.
func ApplyPathOverrides(paths AppPaths, cfg AppConfig) (AppPaths, error) {
	current := paths

	if cfg.Paths.DataDir != "" {
		value, err := expandPath(cfg.Paths.DataDir)
		if err != nil {
			return AppPaths{}, err
		}
		current.DataDir = value
	}

	if cfg.Paths.StateDir != "" {
		value, err := expandPath(cfg.Paths.StateDir)
		if err != nil {
			return AppPaths{}, err
		}
		current.StateDir = value
	}

	return current, nil
}

// EnsureDirectories creates the data and state directories when necessary.
func EnsureDirectories(paths AppPaths, flags CommonFlags) error {
	if flags.DryRun {
		return nil
	}

	for _, dir := range []string{paths.DataDir, paths.StateDir} {
		if dir == "" {
			continue
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}
	return nil
}

func resolveConfigFile(app string, override string) (string, error) {
	if override != "" {
		path, err := expandPath(override)
		if err != nil {
			return "", err
		}
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			return filepath.Join(path, "config.toml"), nil
		}
		return path, nil
	}

	dir, err := defaultConfigDir(app)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}

func defaultConfigDir(app string) (string, error) {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, app), nil
	}

	if runtime.GOOS == "windows" {
		if dir := os.Getenv("APPDATA"); dir != "" {
			return filepath.Join(dir, app), nil
		}
	}

	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, app), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", app), nil
}

func defaultDataDir(app string) (string, error) {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return filepath.Join(dir, app), nil
	}

	if runtime.GOOS == "windows" {
		if dir := os.Getenv("LOCALAPPDATA"); dir != "" {
			return filepath.Join(dir, app), nil
		}
	}

	if dir, err := os.UserCacheDir(); err == nil && dir != "" && runtime.GOOS == "darwin" {
		return filepath.Join(dir, app), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determine home directory: %w", err)
	}
	return filepath.Join(home, ".local", "share", app), nil
}

func defaultStateDir(app string) (string, error) {
	if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
		return filepath.Join(dir, app), nil
	}

	if runtime.GOOS == "windows" {
		if dir := os.Getenv("LOCALAPPDATA"); dir != "" {
			return filepath.Join(dir, app), nil
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determine home directory: %w", err)
	}
	return filepath.Join(home, ".local", "state", app), nil
}

func defaultCacheDir(app string) (string, error) {
	if dir := os.Getenv("XDG_CACHE_HOME"); dir != "" {
		return filepath.Join(dir, app), nil
	}

	if runtime.GOOS == "windows" {
		if dir := os.Getenv("LOCALAPPDATA"); dir != "" {
			return filepath.Join(dir, app, "cache"), nil
		}
	}

	if dir, err := os.UserCacheDir(); err == nil && dir != "" {
		return filepath.Join(dir, app), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determine home directory: %w", err)
	}
	return filepath.Join(home, ".cache", app), nil
}
