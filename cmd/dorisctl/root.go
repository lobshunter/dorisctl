package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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

var rootCmd = cobra.Command{
	Use:   "dorisctl",
	Short: "Dorisctl: a management tool for apache doris",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.MakeGlobalConfig()
	},
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

// askForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user.
// copied from: https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
