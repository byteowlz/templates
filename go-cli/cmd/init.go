package cmd

import (
	"github.com/spf13/cobra"

	"gitlab.cc-asp.fraunhofer.de/templates/go-cli/internal/app"
)

func newInitCommand() *cobra.Command {
	opts := app.InitOptions{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create config directories and default files.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := Context(cmd)
			if err != nil {
				return err
			}
			return app.HandleInit(ctx, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Recreate configuration even if it already exists.")

	return cmd
}
