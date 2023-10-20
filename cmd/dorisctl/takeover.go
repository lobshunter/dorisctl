package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/lobshunter/dorisctl/pkg/cluster"
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/config"
	"github.com/lobshunter/dorisctl/pkg/store"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

func newTakeOverCmd() *cobra.Command {
	var clusterName string

	var (
		feMaster    string
		feHosts     string
		feDeployDir string
		beHosts     string
		beDeployDir string

		deployUser string
		sshPort    int
		sshKeyPath string
	)

	cmd := cobra.Command{
		Use:               "takeover",
		Short:             "Take over a manually deployed(non-systemd) Doris cluster, and make it managed by dorisctl",
		Args:              WrapArgsError(cobra.NoArgs),
		SilenceUsage:      true,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			packageStore := store.NewPackageStore(config.GlobalConfig.CacheDir)
			manifestStore := store.NewLocalManifestStore(config.GlobalConfig.DataDir)

			deployInfo := topologyyaml.ManualDeployInfo{
				DeployUser: deployUser,
				SSHKeyPath: sshKeyPath,
				SSHPort:    sshPort,

				FeMaster:    feMaster,
				FeHosts:     strings.Split(feHosts, ","),
				BeHosts:     strings.Split(beHosts, ","),
				FeDeployDir: feDeployDir,
				BeDeployDir: beDeployDir,
			}
			if !slices.Contains(deployInfo.FeHosts, deployInfo.FeMaster) {
				return fmt.Errorf("fe master %s not in fe hosts %v", deployInfo.FeMaster, deployInfo.FeHosts)
			}

			topo := topologyyaml.BuildTopoFromManualDeploy(deployInfo)
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
	cmd.Flags().StringVar(&feMaster, "fe-master", "", "--fe-master <host>")
	_ = cmd.MarkFlagRequired("fe-master")
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

	home, _ := os.UserHomeDir()
	cmd.Flags().StringVar(&sshKeyPath, "ssh-key-path", filepath.Join(home, ".ssh/id_rsa"), "--ssh-key-path <path>")

	return &cmd
}
