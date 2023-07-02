package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/lobshunter/dorisctl/pkg/cluster"
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/config"
	"github.com/lobshunter/dorisctl/pkg/store"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

func newDeplpyCmd() *cobra.Command {
	var clusterName string
	cmd := cobra.Command{
		Use:          "deploy <topology.yaml>",
		Short:        "Deploy a Doris cluster",
		Args:         WrapArgsError(cobra.ExactArgs(1)),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: support http sourced topo.yaml
			topo, err := topologyyaml.LoadFromFile(args[0])
			if err != nil {
				return err
			}
			err = topologyyaml.Validate(topo)
			if err != nil {
				return err
			}

			packageStore := store.NewPackageStore(config.GlobalConfig.CacheDir)
			manifestStore := store.NewLocalManifestStore(config.GlobalConfig.DataDir)
			taskCtx, err := task.BuildContext(clusterName, *topo, packageStore, manifestStore)
			if err != nil {
				return err
			}

			tasks := cluster.GetDeployClusterSteps(taskCtx)
			err = cluster.RunSteps(context.TODO(), tasks...)
			return err
		},
	}

	cmd.Flags().StringVar(&clusterName, "cluster-name", "default", "cluster name")

	return &cmd
}
