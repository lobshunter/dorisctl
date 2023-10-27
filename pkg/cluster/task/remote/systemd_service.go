package remote

import (
	"context"
	"fmt"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &StartService{}
	_ task.Task = &StopService{}
)

type StartService struct {
	sshClient  *ssh.SSHClient
	instance   topologyyaml.Instance
	useSystemd bool
}

func NewStartService(sshClient *ssh.SSHClient, instance topologyyaml.Instance, useSystemd bool) *StartService {
	return &StartService{
		sshClient:  sshClient,
		instance:   instance,
		useSystemd: useSystemd,
	}
}

func (t *StartService) Name() string {
	return "start-service"
}

func (t *StartService) Execute(ctx context.Context) error {
	var commands []string
	if t.useSystemd {
		commands = []string{
			"systemctl enable " + t.instance.SystemdServiceName(),
			"systemctl start " + t.instance.SystemdServiceName(),
		}
	} else {
		envs := t.instance.EnvironmentVars()
		envString := ""
		for k, v := range envs {
			envString += fmt.Sprintf("%s=%s ", k, v)
		}

		commands = []string{
			// kill -0 checks if the process is running
			fmt.Sprintf("kill -0 $(cat %s) 2>/dev/null || %s %s --daemon", t.instance.PIDFile(), envString, t.instance.StartupScript()),
		}
	}

	return t.sshClient.Run(ctx, commands...)
}

type StopService struct {
	sshClient  *ssh.SSHClient
	instance   topologyyaml.Instance
	useSystemd bool
}

func NewStopService(sshClient *ssh.SSHClient, instance topologyyaml.Instance, useSystemd bool) *StopService {
	return &StopService{
		sshClient:  sshClient,
		instance:   instance,
		useSystemd: useSystemd,
	}
}

func (t *StopService) Name() string {
	return "stop-service"
}

func (t *StopService) Execute(ctx context.Context) error {
	var commands []string
	if t.useSystemd {
		commands = []string{
			"systemctl stop " + t.instance.SystemdServiceName(),
			"systemctl disable " + t.instance.SystemdServiceName(),
		}
	} else {
		commands = []string{
			fmt.Sprintf("kill -0 $(cat %s) 2>/dev/null && %s || true", t.instance.PIDFile(), t.instance.StopScript()),
		}
	}

	return t.sshClient.Run(ctx, commands...)
}
