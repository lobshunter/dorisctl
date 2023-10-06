package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lobshunter/dorisctl/pkg/config"
)

func init() {
	rootCmd.AddCommand(
		newVersionCmd(),
		newDeplpyCmd(),
		newStartCmd(),
		newDestroyCmd(),
		newListCmd(),
		newStopCmd(),
		newStatusCmd(),
		newTakeOverCmd(),
		newHandOverCmd(),
	)
}

// WrapArgsError annotates cobra args error with some context, so the error message is more user-friendly
func WrapArgsError(argFn cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		err := argFn(cmd, args)
		if err == nil {
			return nil
		}

		return fmt.Errorf("%q %s.\nSee '%s --help'.\n\nUsage:  %s\n\n%s",
			cmd.CommandPath(), err.Error(),
			cmd.CommandPath(),
			cmd.UseLine(), cmd.Short,
		)
	}
}

var rootCmd = cobra.Command{
	Use:   "dorisctl",
	Short: "Dorisctl: a management tool for apache doris",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.MakeGlobalConfig()
	},
}
