package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/lobshunter/dorisctl/pkg/cluster"
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/config"
	"github.com/lobshunter/dorisctl/pkg/store"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

func newStatusCmd() *cobra.Command {
	var clusterName string
	cmd := cobra.Command{
		Use:          "status",
		Short:        "Check status of a Doris cluster",
		Args:         WrapArgsError(cobra.NoArgs),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			tasks := cluster.GetCheckClusterStatusSteps(taskCtx)
			err = cluster.RunSteps(cmd.Context(), tasks...)
			if err != nil {
				return err
			}

			printClusterStatus(taskCtx.ClusterStatus)
			return nil
		},
	}

	cmd.Flags().StringVar(&clusterName, "cluster-name", "default", "cluster name")

	return &cmd
}

func printClusterStatus(status topologyyaml.ClusterStatus) {
	if !status.FeMasterHealthy {
		fmt.Println("FeMaster is not healthy")
		return
	}

	fmt.Println("\nFrontends:")
	w := tablewriter.NewWriter(os.Stdout)
	w.SetAlignment(tablewriter.ALIGN_CENTER)
	w.SetHeader([]string{"Host", "QueryPort", "IsMaster", "Alive", "Version"})
	for _, fe := range status.Fes {
		w.Append([]string{
			fe.Host,
			fmt.Sprintf("%d", fe.QueryPort),
			fmt.Sprintf("%v", fe.IsMaster),
			fmt.Sprintf("%v", fe.Alive),
			fe.Version})
	}
	w.Render()

	fmt.Println("\nBackends:")
	beWriter := tablewriter.NewWriter(os.Stdout)
	w.SetAlignment(tablewriter.ALIGN_CENTER)
	beWriter.SetHeader([]string{"Host", "Alive", "AvailCapacity", "TotalCapacity", "UsedPct", "Version"})
	for _, be := range status.Bes {
		beWriter.Append([]string{
			be.Host,
			fmt.Sprintf("%v", be.Alive),
			be.AvailCapacity,
			be.TotalCapacity,
			be.UsedPct,
			be.Version})
	}
	beWriter.Render()
}
