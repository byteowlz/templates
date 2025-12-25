package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// configSchemaJSON contains the JSON schema for the config file.
// This is a literal copy of examples/config.schema.json.
const configSchemaJSON = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://github.com/byteowlz/go-cli/schemas/config.schema.json",
  "title": "go-cli configuration",
  "description": "Configuration schema for go-cli",
  "type": "object",
  "properties": {
    "$schema": {
      "type": "string",
      "description": "JSON Schema reference for editor support"
    },
    "profile": {
      "type": "string",
      "description": "Active configuration profile",
      "default": "default"
    },
    "logging": {
      "type": "object",
      "description": "Logging configuration",
      "properties": {
        "level": {
          "type": "string",
          "description": "Log level",
          "enum": ["error", "warn", "info", "debug", "trace"],
          "default": "info"
        },
        "file": {
          "type": "string",
          "description": "Optional path for log file output. Supports ~ and environment variables."
        }
      },
      "additionalProperties": false
    },
    "runtime": {
      "type": "object",
      "description": "Runtime configuration",
      "properties": {
        "parallelism": {
          "type": "integer",
          "description": "Worker pool size. Defaults to logical CPU count.",
          "minimum": 1
        },
        "timeout": {
          "type": "integer",
          "description": "Timeout in seconds for long-running operations",
          "default": 60,
          "minimum": 1
        },
        "fail_fast": {
          "type": "boolean",
          "description": "Stop on first error",
          "default": true
        }
      },
      "additionalProperties": false
    },
    "paths": {
      "type": "object",
      "description": "Custom paths for data and state",
      "properties": {
        "data_dir": {
          "type": "string",
          "description": "Directory for persistent data. Supports ~ and environment variables."
        },
        "state_dir": {
          "type": "string",
          "description": "Directory for state files. Supports ~ and environment variables."
        }
      },
      "additionalProperties": false
    }
  },
  "additionalProperties": false
}
`

// AppConfig represents the template's configuration schema.
type AppConfig struct {
	Profile string        `mapstructure:"profile" json:"profile" yaml:"profile"`
	Logging LoggingConfig `mapstructure:"logging" json:"logging" yaml:"logging"`
	Runtime RuntimeConfig `mapstructure:"runtime" json:"runtime" yaml:"runtime"`
	Paths   PathsConfig   `mapstructure:"paths" json:"paths" yaml:"paths"`
}

// LoggingConfig controls log output.
type LoggingConfig struct {
	Level string `mapstructure:"level" json:"level" yaml:"level"`
	File  string `mapstructure:"file" json:"file" yaml:"file"`
}

// RuntimeConfig contains runtime tuning parameters.
type RuntimeConfig struct {
	Parallelism    *int `mapstructure:"parallelism" json:"parallelism,omitempty" yaml:"parallelism,omitempty"`
	TimeoutSeconds *int `mapstructure:"timeout" json:"timeout,omitempty" yaml:"timeout,omitempty"`
	FailFast       bool `mapstructure:"fail_fast" json:"fail_fast" yaml:"fail_fast"`
}

// PathsConfig lets users override data/state locations.
type PathsConfig struct {
	DataDir  string `mapstructure:"data_dir" json:"data_dir,omitempty" yaml:"data_dir,omitempty"`
	StateDir string `mapstructure:"state_dir" json:"state_dir,omitempty" yaml:"state_dir,omitempty"`
}

// RunConfig is the subset of AppConfig used by `run`.
type RunConfig struct {
	Profile string        `json:"profile" yaml:"profile"`
	Runtime RuntimeConfig `json:"runtime" yaml:"runtime"`
}

// LoadOrInitConfig ensures the config file exists (unless dry-run) and loads it.
func LoadOrInitConfig(paths AppPaths, flags CommonFlags) (AppConfig, error) {
	if _, err := os.Stat(paths.ConfigFile); os.IsNotExist(err) {
		if flags.DryRun {
			fmt.Fprintf(os.Stderr, "dry-run: would create default config at %s\n", paths.ConfigFile)
		} else {
			if err := writeDefaultConfig(paths.ConfigFile); err != nil {
				return AppConfig{}, err
			}
		}
	} else if err != nil {
		return AppConfig{}, fmt.Errorf("failed to stat config file: %w", err)
	}

	cfg := defaultConfig()

	v := viper.New()
	v.SetConfigFile(paths.ConfigFile)
	v.SetConfigType("toml")
	v.SetEnvPrefix(EnvPrefix())
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	v.AutomaticEnv()

	v.SetDefault("profile", cfg.Profile)
	v.SetDefault("logging.level", cfg.Logging.Level)
	v.SetDefault("runtime.timeout", 60)
	v.SetDefault("runtime.fail_fast", true)

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) && !(flags.DryRun && os.IsNotExist(err)) {
			return AppConfig{}, err
		}
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return AppConfig{}, fmt.Errorf("decode config: %w", err)
	}

	if cfg.Logging.File != "" {
		expanded, err := expandPath(cfg.Logging.File)
		if err != nil {
			return AppConfig{}, fmt.Errorf("expand log file path: %w", err)
		}
		cfg.Logging.File = expanded
	}

	if cfg.Runtime.TimeoutSeconds == nil {
		defaultTimeout := 60
		cfg.Runtime.TimeoutSeconds = &defaultTimeout
	}

	return cfg, nil
}

// WithProfileOverride returns a shallow copy with the profile overridden.
func (cfg AppConfig) WithProfileOverride(profile string) AppConfig {
	if profile != "" {
		cfg.Profile = profile
	}
	return cfg
}

func (cfg AppConfig) RunConfig() RunConfig {
	return RunConfig{
		Profile: cfg.Profile,
		Runtime: cfg.Runtime,
	}
}

func writeDefaultConfig(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	body := strings.Builder{}
	body.WriteString(defaultConfigHeader(path))
	body.WriteString(defaultConfigBody())
	if err := os.WriteFile(path, []byte(body.String()), 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func defaultConfigHeader(path string) string {
	return fmt.Sprintf("# Configuration for %s\n# File: %s\n\n", appName, path)
}

func defaultConfigBody() string {
	return `profile = "default"

[logging]
# Valid levels: error, warn, info, debug, trace
level = "info"
# Optional path for log file output; supports ~ and environment variables.
# file = "~/Library/Logs/` + appName + `.log"

[runtime]
# Override the worker pool size; defaults to logical CPU count when unset.
# parallelism = 8
# Timeout in seconds for long-running operations.
timeout = 60
fail_fast = true

[paths]
# Uncomment to move persistent data/state to custom directories.
# data_dir = "$XDG_DATA_HOME/` + appName + `"
# state_dir = "$XDG_STATE_HOME/` + appName + `"
`
}

func defaultConfig() AppConfig {
	defaultTimeout := 60
	return AppConfig{
		Profile: "default",
		Logging: LoggingConfig{
			Level: "info",
		},
		Runtime: RuntimeConfig{
			TimeoutSeconds: &defaultTimeout,
			FailFast:       true,
		},
	}
}

// TimeoutDuration returns the configured timeout as a time.Duration.
func (cfg RuntimeConfig) TimeoutDuration() time.Duration {
	if cfg.TimeoutSeconds == nil {
		return 0
	}
	return time.Duration(*cfg.TimeoutSeconds) * time.Second
}
