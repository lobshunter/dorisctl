package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/lobshunter/dorisctl/pkg/cluster"
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/config"
	"github.com/lobshunter/dorisctl/pkg/store"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

func newExecCmd() *cobra.Command {
	var clusterName string
	var hosts string
	var be, fe bool

	cmd := cobra.Command{
		Use:               "exec",
		Short:             "Execute a command on hosts of a cluster",
		Args:              WrapArgsError(cobra.ExactArgs(1)),
		ValidArgsFunction: cobra.NoFileCompletions,
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

			command := args[0]
			execHosts := getExecHosts(topo, be, fe, strings.Split(hosts, ","))
			execs := map[string]*topologyyaml.ExecOutput{}
			for _, host := range execHosts {
				execs[host] = &topologyyaml.ExecOutput{}
			}

			tasks := cluster.GetExecSteps(taskCtx, execs, command)
			err = cluster.RunSteps(cmd.Context(), tasks...)
			if err != nil {
				return err
			}

			printExecResult(execs)

			return nil
		},
	}

	cmd.Flags().StringVar(&clusterName, "cluster", "default", "cluster name")
	_ = cmd.RegisterFlagCompletionFunc("cluster", completeClusterName)
	cmd.Flags().BoolVar(&be, "be", false, "execute command on BE hosts")
	cmd.Flags().BoolVar(&fe, "fe", false, "execute command on FE hosts")
	cmd.Flags().StringVar(&hosts, "hosts", "", "hosts to execute command on, --hosts=<host1,host2,...>")

	return &cmd
}

func getExecHosts(topo *topologyyaml.Topology, be, fe bool, literalHosts []string) []string {
	allClusterHosts := make([]string, 0, len(topo.BEs)+len(topo.FEs))
	for _, beSpec := range topo.BEs {
		allClusterHosts = append(allClusterHosts, beSpec.Host)
	}
	for _, feSpec := range topo.FEs {
		allClusterHosts = append(allClusterHosts, feSpec.Host)
	}

	var hosts []string
	for _, h := range literalHosts {
		if slices.Contains(allClusterHosts, h) {
			hosts = append(hosts, h)
		}
	}
	if be {
		for _, beSpec := range topo.BEs {
			hosts = append(hosts, beSpec.Host)
		}
	}
	if fe {
		for _, feSpec := range topo.FEs {
			hosts = append(hosts, feSpec.Host)
		}
	}

	slices.Sort(hosts)
	return slices.Compact(hosts)
}

func printExecResult(execs map[string]*topologyyaml.ExecOutput) {
	for host, exec := range execs {
		if exec.ExitCode == 0 {
			color.Green(host)
			fmt.Printf("stdout: %s\n", exec.Stdout)
		} else {
			color.Red("%s %d", host, exec.ExitCode)
			fmt.Printf("stdout: %s\n", exec.Stdout)
			fmt.Printf("stderr: %s\n", exec.Stderr)
		}
	}
}
