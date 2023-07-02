package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lobshunter/dorisctl/pkg/version"
)

func newVersionCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "version",
		Short: "Print the version info of dorisctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Version())
		},
	}

	return &cmd
}
