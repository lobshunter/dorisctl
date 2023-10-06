package main

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/lobshunter/dorisctl/pkg/cluster"
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/config"
	"github.com/lobshunter/dorisctl/pkg/store"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

func newTakeOverCmd() *cobra.Command {
	var clusterName string

	var (
		feHosts     string
		feDeployDir string
		beHosts     string
		beDeployDir string

		deployUser string
		sshPort    int
		sshKeyPath string
	)

	cmd := cobra.Command{
		Use:   "takeover",
		Short: "Take over a manually deployed(non-systemd) Doris cluster, and make it managed by dorisctl",
		Args:  WrapArgsError(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageStore := store.NewPackageStore(config.GlobalConfig.CacheDir)
			manifestStore := store.NewLocalManifestStore(config.GlobalConfig.DataDir)

			topo := topologyyaml.BuildTopoFromManualDeploy(topologyyaml.ManualDeployInfo{
				DeployUser: deployUser,
				SSHKeyPath: sshKeyPath,
				SSHPort:    sshPort,

				FeHosts:     strings.Split(feHosts, ","),
				BeHosts:     strings.Split(beHosts, ","),
				FeDeployDir: feDeployDir,
				BeDeployDir: beDeployDir,
			})

			taskCtx, err := task.BuildContext(clusterName, *topo, packageStore, manifestStore)
			if err != nil {
				return err
			}

			tasks := cluster.GetTakeOverClusterSteps(taskCtx)
			err = cluster.RunSteps(cmd.Context(), tasks...)
			return err
		},
	}

	cmd.Flags().StringVar(&clusterName, "cluster-name", "tookover", "cluster name")

	// TODO: maybe don't require hosts, but retrieve them via Fe query port(just need a running fe master)
	cmd.Flags().StringVar(&feHosts, "fe-hosts", "", "--fe-hosts <host1,host2,...>")
	_ = cmd.MarkFlagRequired("fe-hosts")
	cmd.Flags().StringVar(&beHosts, "be-hosts", "", "--be-hosts <host1,host2,...>")
	_ = cmd.MarkFlagRequired("be-hosts")
	cmd.Flags().StringVar(&feDeployDir, "fe-deploy-dir", "", "--fe-deploy-dir <dir>")
	_ = cmd.MarkFlagRequired("fe-deploy-dir")
	cmd.Flags().StringVar(&beDeployDir, "be-deploy-dir", "", "--be-deploy-dir <dir>")
	_ = cmd.MarkFlagRequired("be-deploy-dir")
	cmd.Flags().StringVar(&deployUser, "deploy-user", "root", "--deploy-user <user>")
	cmd.Flags().IntVar(&sshPort, "ssh-port", 22, "--ssh-port <port>")
	cmd.Flags().StringVar(&sshKeyPath, "ssh-key-path", "", "--ssh-key-path <path>")

	return &cmd
}
