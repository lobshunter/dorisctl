package remote

import (
	"context"
	"fmt"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &ResetInstanceHost{}
)

type ResetInstanceUniqueHost struct {
	sshClient  *ssh.SSHClient
	deployUser string
}

func NewResetInstanceUniqueHost(sshClient *ssh.SSHClient, deployUser string) *ResetInstanceUniqueHost {
	return &ResetInstanceUniqueHost{
		sshClient:  sshClient,
		deployUser: deployUser,
	}
}

func (t *ResetInstanceUniqueHost) Name() string {
	return "reset-instance-unique-host"
}

func (t *ResetInstanceUniqueHost) Execute(ctx context.Context) error {
	if t.deployUser != "root" {
		// ignore delete user error
		cmd := fmt.Sprintf("userdel -r %s || true", t.deployUser)
		return t.sshClient.Run(ctx, cmd)
	}

	return nil
}

type ResetInstanceHost struct {
	sshClient          *ssh.SSHClient
	deployUser         string
	deployDir          string
	SystemdServicePath string
}

func NewResetInstanceHost(sshClient *ssh.SSHClient, deployUser string, deployDir string, systemdServicePath string) *ResetInstanceHost {
	return &ResetInstanceHost{
		sshClient:          sshClient,
		deployUser:         deployUser,
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
