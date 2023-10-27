package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lobshunter/dorisctl/pkg/config"
)

func newSelfCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "self",
		Short: "Shows the information about dorisctl itself",
		Args:  WrapArgsError(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("cache dir: %s\n", config.GlobalConfig.CacheDir)
			fmt.Printf("data dir: %s\n", config.GlobalConfig.DataDir)
			return nil
		},
	}

	return &cmd
}
