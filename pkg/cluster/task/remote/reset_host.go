package remote

import (
	"context"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &ResetInstanceHost{}
)

type ResetInstanceHost struct {
	sshClient          *ssh.SSHClient
	deployDir          string
	SystemdServicePath string
}

func NewResetInstanceHost(sshClient *ssh.SSHClient, deployDir string, systemdServicePath string) *ResetInstanceHost {
	return &ResetInstanceHost{
		sshClient:          sshClient,
		deployDir:          deployDir,
		SystemdServicePath: systemdServicePath,
	}
}

func (t *ResetInstanceHost) Name() string {
	return "reset-instance-host"
}

func (t *ResetInstanceHost) Execute(ctx context.Context) error {
	commands := []string{
		"rm -rf " + t.deployDir,
		"rm " + t.SystemdServicePath,
	}

	return t.sshClient.Run(ctx, commands...)
}
