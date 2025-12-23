package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func newCompletionsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completions [shell]",
		Short: "Generate shell completion scripts.",
		Long:  "Generate shell completion scripts for supported shells (bash, zsh, fish, powershell).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := args[0]
			root := cmd.Root()
			var err error

			switch shell {
			case "bash":
				err = root.GenBashCompletion(os.Stdout)
			case "zsh":
				err = root.GenZshCompletion(os.Stdout)
			case "fish":
				err = root.GenFishCompletion(os.Stdout, true)
			case "powershell":
				err = root.GenPowerShellCompletionWithDesc(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell %q", shell)
			}

			if err != nil && err != io.EOF {
				return err
			}
			return nil
		},
	}
}
