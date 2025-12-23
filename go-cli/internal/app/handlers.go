package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// RunOptions configure the run command behaviour.
type RunOptions struct {
	Task    string
	Profile string
}

// InitOptions configure the init command behaviour.
type InitOptions struct {
	Force bool
}

// HandleRun executes the run command.
func HandleRun(ctx *RuntimeContext, opts RunOptions) error {
	effective := ctx.Config.WithProfileOverride(opts.Profile)
	runCfg := effective.RunConfig()

	if ctx.Common.Parallelism != nil {
		value := *ctx.Common.Parallelism
		runCfg.Runtime.Parallelism = &value
	}

	if runCfg.Runtime.Parallelism == nil {
		value := defaultParallelism()
		runCfg.Runtime.Parallelism = &value
	}

	if ctx.Common.TimeoutSeconds != nil {
		value := *ctx.Common.TimeoutSeconds
		runCfg.Runtime.TimeoutSeconds = &value
	}

	parallelism := 0
	if runCfg.Runtime.Parallelism != nil {
		parallelism = *runCfg.Runtime.Parallelism
	}

	timeout := 0
	if runCfg.Runtime.TimeoutSeconds != nil {
		timeout = *runCfg.Runtime.TimeoutSeconds
	}

	ctx.Logger.Info("running task %s with profile %s", opts.Task, runCfg.Profile)

	result := map[string]any{
		"task":        opts.Task,
		"profile":     runCfg.Profile,
		"parallelism": parallelism,
		"timeout":     timeout,
	}

	switch {
	case ctx.Common.JSON:
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case ctx.Common.YAML:
		data, err := yaml.Marshal(result)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
	default:
		fmt.Printf("Running task %q with profile %q (parallelism: %d, timeout: %ds)\n", opts.Task, runCfg.Profile, parallelism, timeout)
	}

	return nil
}

// HandleInit creates the config if necessary.
func HandleInit(ctx *RuntimeContext, opts InitOptions) error {
	path := ctx.Paths.ConfigFile
	if _, err := os.Stat(path); err == nil && !(opts.Force || ctx.Common.AssumeYes) {
		return fmt.Errorf("config already exists at %s (use --force to overwrite)", path)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if ctx.Common.DryRun {
		ctx.Logger.Info("dry-run: would write default config to %s", path)
		return nil
	}

	if err := writeDefaultConfig(path); err != nil {
		return err
	}

	ctx.Logger.Info("wrote default config to %s", path)
	return nil
}

// HandleConfigShow prints the effective configuration.
func HandleConfigShow(ctx *RuntimeContext) error {
	switch {
	case ctx.Common.JSON:
		data, err := json.MarshalIndent(ctx.Config, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case ctx.Common.YAML:
		data, err := yaml.Marshal(ctx.Config)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
	default:
		fmt.Printf("%+v\n", ctx.Config)
	}
	return nil
}

// HandleConfigPath prints the config path.
func HandleConfigPath(ctx *RuntimeContext) error {
	fmt.Println(ctx.Paths.ConfigFile)
	return nil
}

// HandleConfigReset rewrites the default config file.
func HandleConfigReset(ctx *RuntimeContext) error {
	if ctx.Common.DryRun {
		ctx.Logger.Info("dry-run: would reset config at %s", ctx.Paths.ConfigFile)
		return nil
	}

	if err := writeDefaultConfig(ctx.Paths.ConfigFile); err != nil {
		return err
	}

	ctx.Logger.Info("reset config at %s", ctx.Paths.ConfigFile)
	return nil
}
