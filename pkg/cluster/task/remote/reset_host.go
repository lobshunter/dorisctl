package remote

import (
	"context"
	"strings"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
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

type ResetBeHost struct {
	sshClient  *ssh.SSHClient
	beInstance topologyyaml.BeInstance
}

func NewResetBeHost(sshClient *ssh.SSHClient, beInstance topologyyaml.BeInstance) *ResetBeHost {
	return &ResetBeHost{
		sshClient:  sshClient,
		beInstance: beInstance,
	}
}

func (t *ResetBeHost) Name() string {
	return "reset-be-host"
}

func (t *ResetBeHost) Execute(ctx context.Context) error {
	storageRootPaths := strings.Split(t.beInstance.BeConfig.StorageRootPath, ";")
	commands := []string{
		"rm -rf " + strings.Join(storageRootPaths, " "),
	}

	return t.sshClient.Run(ctx, commands...)
}
