package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/lobshunter/dorisctl/pkg/config"
	"github.com/lobshunter/dorisctl/pkg/store"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

func newListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:          "list",
		Short:        "List all Doris clusters",
		Args:         WrapArgsError(cobra.NoArgs),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestStore := store.NewLocalManifestStore(config.GlobalConfig.DataDir)
			clusters, err := manifestStore.List()
			if err != nil {
				return fmt.Errorf("failed to list clusters %v", err)
			}

			clusterMap := make(map[string]*topologyyaml.Topology)
			for _, cluster := range clusters {
				topo, err := manifestStore.Get(cluster)
				if err != nil {
					return fmt.Errorf("failed to get cluster %s %v", cluster, err)
				}
				clusterMap[cluster] = topo
			}

			printClusterList(clusterMap)
			return nil
		},
	}

	return &cmd
}

func printClusterList(clusters map[string]*topologyyaml.Topology) {
	w := tablewriter.NewWriter(os.Stdout)
	w.SetAlignment(tablewriter.ALIGN_CENTER)
	w.SetHeader([]string{"Cluster Name", "FeMaster", "Fe Quert Port"})
	for name, topo := range clusters {
		for _, fe := range topo.FEs {
			if !fe.IsMaster {
				continue
			}

			w.Append([]string{name, fe.Host, fmt.Sprintf("%d", fe.FeConfig.QueryPort)})
		}
	}
	w.Render()
}
