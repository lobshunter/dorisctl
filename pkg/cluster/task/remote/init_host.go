package remote

import (
	"context"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &InitInstanceHost{}
)

type InitInstanceHost struct {
	sshClient  *ssh.SSHClient
	deployUser string
	deployDir  string
}

func NewInitInstanceHost(sshClient *ssh.SSHClient, deployUser string, deployDir string) *InitInstanceHost {
	return &InitInstanceHost{
		sshClient:  sshClient,
		deployUser: deployUser,
		deployDir:  deployDir,
	}
}

func (t *InitInstanceHost) Name() string {
	return "init-instance-host-dir"
}

func (t *InitInstanceHost) Execute(ctx context.Context) error {
	commands := []string{
		"mkdir -p " + t.deployDir,
		"chown -R " + t.deployUser + " " + t.deployDir,
	}

	return t.sshClient.Run(ctx, commands...)
}
