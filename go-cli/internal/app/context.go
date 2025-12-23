package app

import "context"

const appName = "go-cli"

// ContextKey is used to store the runtime context in a context.Context.
type ContextKey struct{}

// RuntimeContext holds shared state for command handlers.
type RuntimeContext struct {
	context.Context
	Common      CommonFlags
	Paths       AppPaths
	Config      AppConfig
	Logger      Logger
	LogSettings LogSettings
}

// NewRuntimeContext builds a runtime context from CLI flags and the current environment.
func NewRuntimeContext(parent context.Context, flags CommonFlags) (*RuntimeContext, error) {
	if parent == nil {
		parent = context.Background()
	}

	paths, err := DiscoverPaths(appName, flags.ConfigPath)
	if err != nil {
		return nil, err
	}

	cfg, err := LoadOrInitConfig(paths, flags)
	if err != nil {
		return nil, err
	}

	effPaths, err := ApplyPathOverrides(paths, cfg)
	if err != nil {
		return nil, err
	}

	if err := EnsureDirectories(effPaths, flags); err != nil {
		return nil, err
	}

	logSettings, err := ResolveLogSettings(flags, cfg)
	if err != nil {
		return nil, err
	}
	logger := ConfigureLogger(logSettings)

	rtx := &RuntimeContext{
		Context:     parent,
		Common:      flags,
		Paths:       effPaths,
		Config:      cfg,
		Logger:      logger,
		LogSettings: logSettings,
	}

	rtx.Context = context.WithValue(parent, ContextKey{}, rtx)
	rtx.Logger.Debug("resolved paths: config=%s data=%s state=%s", rtx.Paths.ConfigFile, rtx.Paths.DataDir, rtx.Paths.StateDir)

	return rtx, nil
}

// FromContext extracts the runtime context instance from ctx, if present.
func FromContext(ctx context.Context) (*RuntimeContext, bool) {
	if ctx == nil {
		return nil, false
	}
	if rtx, ok := ctx.Value(ContextKey{}).(*RuntimeContext); ok && rtx != nil {
		return rtx, true
	}
	return nil, false
}

// EnvPrefix returns the environment variable prefix for configuration overrides.
func EnvPrefix() string {
	return toEnvPrefix(appName)
}

// Close releases resources held by the runtime context.
func (rtx *RuntimeContext) Close() error {
	if rtx == nil {
		return nil
	}
	return rtx.Logger.Close()
}
