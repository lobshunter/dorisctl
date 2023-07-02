package remote

import (
	"context"
	"fmt"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &InitInstanceUniqueHost{}
	_ task.Task = &InitInstanceHost{}
)

type InitInstanceUniqueHost struct {
	sshClient  *ssh.SSHClient
	deployUser string
}

func NewInitInstanceUniqueHost(sshClient *ssh.SSHClient, deployUser string) *InitInstanceUniqueHost {
	return &InitInstanceUniqueHost{
		sshClient:  sshClient,
		deployUser: deployUser,
	}
}

func (t *InitInstanceUniqueHost) Name() string {
	return "init-instance-host"
}

func (t *InitInstanceUniqueHost) Execute(ctx context.Context) error {
	commands := []string{
		fmt.Sprintf("id -u %s || useradd -m %s", t.deployUser, t.deployUser),
		"mkdir -p /home/" + t.deployUser + "/.ssh",
		"chmod 700 /home/" + t.deployUser + "/.ssh",
		// add ssh pub key to authorized_keys if not exists
		"grep " + t.sshClient.AuthorizedKey() + " /home/" + t.deployUser + "/.ssh/authorized_keys || " +
			"echo " + t.sshClient.AuthorizedKey() + " >> /home/" + t.deployUser + "/.ssh/authorized_keys",
	}

	return t.sshClient.Run(ctx, commands...)
}

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
