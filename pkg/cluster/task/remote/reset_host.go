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
	useSystemd         bool
	SystemdServicePath string
}

func NewResetInstanceHost(sshClient *ssh.SSHClient, deployDir string, useSystemd bool, systemdServicePath string) *ResetInstanceHost {
	return &ResetInstanceHost{
		sshClient:          sshClient,
		deployDir:          deployDir,
		useSystemd:         useSystemd,
		SystemdServicePath: systemdServicePath,
	}
}

func (t *ResetInstanceHost) Name() string {
	return "reset-instance-host"
}

func (t *ResetInstanceHost) Execute(ctx context.Context) error {
	commands := []string{
		"rm -rf " + t.deployDir,
	}

	if t.useSystemd {
		commands = append(commands, "rm -f "+t.SystemdServicePath)
	}

	return t.sshClient.Run(ctx, commands...)
}
