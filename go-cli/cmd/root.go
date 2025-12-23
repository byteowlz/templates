package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"gitlab.cc-asp.fraunhofer.de/templates/go-cli/internal/app"
)

var (
	rootCmd      *cobra.Command
	commonFlags  app.CommonFlags
	timeoutFlag  int
	parallelFlag int
)

func init() {
	rootCmd = &cobra.Command{
		Use:           "go-cli",
		Short:         "Opinionated starting point for cross-platform Go CLIs.",
		Long:          "go-cli is a batteries-included template demonstrating structured commands, config loading, logging, and shell completion generation.",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if _, ok := app.FromContext(cmd.Context()); ok {
				return nil
			}

			flags := commonFlags
			if f := cmd.Flags().Lookup("timeout"); f != nil && f.Changed {
				flags.TimeoutSeconds = &timeoutFlag
			}
			if f := cmd.Flags().Lookup("parallel"); f != nil && f.Changed {
				flags.Parallelism = &parallelFlag
			}

			if flags.JSON && flags.YAML {
				return fmt.Errorf("--json and --yaml cannot be used together")
			}

			if err := flags.ValidateColor(); err != nil {
				return err
			}

			rtx, err := app.NewRuntimeContext(cmd.Context(), flags)
			if err != nil {
				return err
			}

			cmd.SetContext(rtx.Context)
			return nil
		},
	}

	pflags := rootCmd.PersistentFlags()
	pflags.StringVar(&commonFlags.ConfigPath, "config", "", "Override the config file path.")
	pflags.BoolVarP(&commonFlags.Quiet, "quiet", "q", false, "Reduce output to only errors.")
	pflags.CountVarP(&commonFlags.Verbose, "verbose", "v", "Increase logging verbosity (stackable).")
	pflags.BoolVar(&commonFlags.Debug, "debug", false, "Enable debug logging (equivalent to -vv).")
	pflags.BoolVar(&commonFlags.Trace, "trace", false, "Enable trace logging (overrides other levels).")
	pflags.BoolVar(&commonFlags.JSON, "json", false, "Output machine-readable JSON.")
	pflags.BoolVar(&commonFlags.YAML, "yaml", false, "Output machine-readable YAML.")
	pflags.BoolVar(&commonFlags.NoColor, "no-color", false, "Disable ANSI colors in output.")
	pflags.StringVar(&commonFlags.Color, "color", "auto", "Color output policy: auto, always, or never.")
	pflags.BoolVar(&commonFlags.DryRun, "dry-run", false, "Do not change anything on disk.")
	pflags.BoolVarP(&commonFlags.AssumeYes, "yes", "y", false, "Assume yes for interactive prompts (alias for --force).")
	pflags.BoolVar(&commonFlags.NoProgress, "no-progress", false, "Disable progress indicators.")
	pflags.BoolVar(&commonFlags.Diagnostics, "diagnostics", false, "Emit additional diagnostics for troubleshooting.")
	pflags.IntVar(&timeoutFlag, "timeout", 0, "Maximum seconds to allow an operation to run.")
	pflags.IntVar(&parallelFlag, "parallel", 0, "Override the degree of parallelism.")

	rootCmd.AddCommand(newRunCommand())
	rootCmd.AddCommand(newInitCommand())
	rootCmd.AddCommand(newConfigCommand())
	rootCmd.AddCommand(newCompletionsCommand())
}

// Execute runs the CLI.
func Execute() error {
	return rootCmd.Execute()
}

// Context extracts the runtime context from a command.
func Context(cmd *cobra.Command) (*app.RuntimeContext, error) {
	rtx, ok := app.FromContext(cmd.Context())
	if !ok {
		return nil, fmt.Errorf("internal error: runtime context not initialized")
	}
	return rtx, nil
}
