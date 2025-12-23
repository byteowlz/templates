package cmd

import (
	"github.com/spf13/cobra"

	"gitlab.cc-asp.fraunhofer.de/templates/go-cli/internal/app"
)

func newRunCommand() *cobra.Command {
	opts := app.RunOptions{
		Task: "default",
	}

	cmd := &cobra.Command{
		Use:   "run [TASK]",
		Short: "Execute the CLI's primary behavior.",
		Long:  "Runs the template's core workflow. Specify an optional task name and override the active profile if desired.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Task = args[0]
			}

			ctx, err := Context(cmd)
			if err != nil {
				return err
			}

			return app.HandleRun(ctx, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Profile, "profile", "", "Override the profile to run under.")

	return cmd
}
