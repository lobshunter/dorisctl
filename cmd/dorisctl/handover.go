package main

import (
	"github.com/lobshunter/dorisctl/pkg/cluster"
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/config"
	"github.com/lobshunter/dorisctl/pkg/store"
	"github.com/spf13/cobra"
)

func newHandOverCmd() *cobra.Command {
	var clusterName string
	var yes bool
	cmd := cobra.Command{
		Use:               "handover",
		Short:             "Hand over a managed Doris cluster",
		Long:              "Hand over a managed Doris cluster, stop dorisctl from managing it (i.e. remove manifest from dorisctl, but not delete the cluster)",
		Args:              WrapArgsError(cobra.NoArgs),
		ValidArgsFunction: cobra.NoFileCompletions,
		SilenceUsage:      true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				sure := askForConfirmation("Are you sure you want to hand over the cluster? (y/n) ")
				if !sure {
					return nil
				}
			}

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

			tasks := cluster.GetHandOverClusterSteps(taskCtx)
			err = cluster.RunSteps(cmd.Context(), tasks...)
			return err
		},
	}

	cmd.Flags().StringVar(&clusterName, "cluster-name", "tookover", "cluster name")
	_ = cmd.RegisterFlagCompletionFunc("cluster-name", completeClusterName)
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return &cmd
}
