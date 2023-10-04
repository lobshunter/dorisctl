package task

import (
	"context"
	"fmt"

	"github.com/lobshunter/dorisctl/pkg/store"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

type Task interface {
	Name() string
	Execute(ctx context.Context) error
}

// TODO: perhaps wrap Context into interface, and disallow updates explicitly
type Context struct {
	ClusterName   string
	TopoYaml      topologyyaml.Topology
	SSHClients    map[string]*ssh.SSHClient
	PackageStore  *store.PackageStore
	ManifestStore store.ManifestStore

	ClusterStatus topologyyaml.ClusterStatus
}

func BuildContext(clusterName string, topo topologyyaml.Topology,
	packageStore *store.PackageStore, manifestStore store.ManifestStore) (*Context, error) {
	sshClients := make(map[string]*ssh.SSHClient)
	for _, fe := range topo.FEs {
		if _, exist := sshClients[fe.Host]; exist {
			continue
		}

		sshClient, err := ssh.NewSSHClient(fmt.Sprintf("%s:%d", fe.Host, fe.SSHPort), topo.Global.DeployUser, topo.Global.SSHPrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("create sshclient on %s: %v", fe.Host, err)
		}

		sshClients[fe.Host] = sshClient
	}
	for _, be := range topo.BEs {
		if _, exist := sshClients[be.Host]; exist {
			continue
		}

		sshClient, err := ssh.NewSSHClient(fmt.Sprintf("%s:%d", be.Host, be.SSHPort), topo.Global.DeployUser, topo.Global.SSHPrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("create sshclient on %s: %v", be.Host, err)
		}

		sshClients[be.Host] = sshClient
	}

	return &Context{
		ClusterName:   clusterName,
		TopoYaml:      topo,
		SSHClients:    sshClients,
		PackageStore:  packageStore,
		ManifestStore: manifestStore,
	}, nil
}

func (c *Context) SetClusterStatus(status topologyyaml.ClusterStatus) {
	// TODO: read/write lock
	c.ClusterStatus = status
}
