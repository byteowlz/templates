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
