package main

import (
	"github.com/spf13/cobra"

	"github.com/lobshunter/dorisctl/pkg/cluster"
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/config"
	"github.com/lobshunter/dorisctl/pkg/store"
)

func newDestroyCmd() *cobra.Command {
	var yes bool
	cmd := cobra.Command{
		Use:               "destroy",
		Short:             "Destroy a Doris cluster",
		Args:              WrapArgsError(cobra.ExactArgs(1)),
		ValidArgsFunction: completeClusterName,
		SilenceUsage:      true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				sure := askForConfirmation("Are you sure you want to destroy the cluster? (y/n) ")
				if !sure {
					return nil
				}
			}

			clusterName := args[0]
			packageStore := store.NewPackageStore(config.GlobalConfig.CacheDir)
			manifestStore := store.NewLocalManifestStore(config.GlobalConfig.DataDir)
			topo, err := manifestStore.Get(clusterName)
			if err != nil {
				return err
			}

			taskCtx, err := task.BuildContext(clusterName, *topo, packageStore, manifestStore)
			if err != nil {
				return err
			}

			tasks := cluster.GetDestroyClusterSteps(taskCtx)
			err = cluster.RunSteps(cmd.Context(), tasks...)
			return err
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return &cmd
}
