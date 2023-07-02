package remote

import (
	"context"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &StartService{}
	_ task.Task = &StopService{}
)

type StartService struct {
	sshClient   *ssh.SSHClient
	serviceName string
}

func NewStartService(sshClient *ssh.SSHClient, serviceName string) *StartService {
	return &StartService{
		sshClient:   sshClient,
		serviceName: serviceName,
	}
}

func (t *StartService) Name() string {
	return "start-service"
}

func (t *StartService) Execute(ctx context.Context) error {
	commands := []string{
		"systemctl enable " + t.serviceName,
		"systemctl start " + t.serviceName,
	}
	return t.sshClient.Run(ctx, commands...)
}

type StopService struct {
	sshClient   *ssh.SSHClient
	serviceName string
}

func NewStopService(sshClient *ssh.SSHClient, serviceName string) *StopService {
	return &StopService{
		sshClient:   sshClient,
		serviceName: serviceName,
	}
}

func (t *StopService) Name() string {
	return "stop-service"
}

func (t *StopService) Execute(ctx context.Context) error {
	commands := []string{
		"systemctl stop " + t.serviceName,
		"systemctl disable " + t.serviceName,
	}
	return t.sshClient.Run(ctx, commands...)
}
