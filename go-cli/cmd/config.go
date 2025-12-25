package cmd

import (
	"github.com/spf13/cobra"

	"gitlab.cc-asp.fraunhofer.de/templates/go-cli/internal/app"
)

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect and manage configuration.",
	}

	cmd.AddCommand(newConfigShowCommand())
	cmd.AddCommand(newConfigPathCommand())
	cmd.AddCommand(newConfigPathsCommand())
	cmd.AddCommand(newConfigSchemaCommand())
	cmd.AddCommand(newConfigResetCommand())

	return cmd
}

func newConfigShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Output the effective configuration.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := Context(cmd)
			if err != nil {
				return err
			}
			return app.HandleConfigShow(ctx)
		},
	}
}

func newConfigPathCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the resolved config file path.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := Context(cmd)
			if err != nil {
				return err
			}
			return app.HandleConfigPath(ctx)
		},
	}
}

func newConfigPathsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "paths",
		Short: "Print all resolved paths (config, data, state, cache).",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := Context(cmd)
			if err != nil {
				return err
			}
			return app.HandleConfigPaths(ctx)
		},
	}
}

func newConfigSchemaCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "schema",
		Short: "Print the JSON schema for the config file.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := Context(cmd)
			if err != nil {
				return err
			}
			return app.HandleConfigSchema(ctx)
		},
	}
}

func newConfigResetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Regenerate the default configuration file.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := Context(cmd)
			if err != nil {
				return err
			}
			return app.HandleConfigReset(ctx)
		},
	}
}
