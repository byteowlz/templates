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

// resolveBase applies the XDG Base Directory Specification: the variable is
// used only when it holds an absolute path (unset, empty and relative values are
// ignored), else $HOME/<unixRel>. Windows falls back to its known folder instead,
// matching the rust-cli template. macOS uses the XDG defaults, not
// ~/Library/Application Support.
func resolveBase(xdg, home, winDir string, isWindows bool, unixRel string) (string, error) {
	if xdg != "" && filepath.IsAbs(xdg) {
		return xdg, nil
	}
	if isWindows {
		if winDir == "" {
			return "", fmt.Errorf("unable to determine base directory")
		}
		return winDir, nil
	}
	if home == "" {
		return "", fmt.Errorf("unable to determine home directory")
	}
	return filepath.Join(home, unixRel), nil
}

func baseDir(xdgVar, unixRel, winVar string) (string, error) {
	home, _ := os.UserHomeDir()
	base, err := resolveBase(os.Getenv(xdgVar), home, os.Getenv(winVar), runtime.GOOS == "windows", unixRel)
	if err != nil {
		return "", fmt.Errorf("%w (%s)", err, xdgVar)
	}
	return base, nil
}

func appDir(app, xdgVar, unixRel, winVar string) (string, error) {
	base, err := baseDir(xdgVar, unixRel, winVar)
	if err != nil {
		return "", err
	}
	return filepath.Join(base, app), nil
}

func defaultConfigDir(app string) (string, error) {
	return appDir(app, "XDG_CONFIG_HOME", ".config", "APPDATA")
}

func defaultDataDir(app string) (string, error) {
	return appDir(app, "XDG_DATA_HOME", filepath.Join(".local", "share"), "APPDATA")
}

func defaultStateDir(app string) (string, error) {
	return appDir(app, "XDG_STATE_HOME", filepath.Join(".local", "state"), "LOCALAPPDATA")
}

func defaultCacheDir(app string) (string, error) {
	return appDir(app, "XDG_CACHE_HOME", ".cache", "LOCALAPPDATA")
}
