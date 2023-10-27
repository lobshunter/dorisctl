package step

import (
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/cluster/task/remote"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

func NewExecCmd(ctx *task.Context, hostsExec map[string]*topologyyaml.ExecOutput, command string) Step {
	tasks := make([]task.Task, 0, len(hostsExec))
	for host, output := range hostsExec {
		tasks = append(tasks, remote.NewExecCmd(ctx.SSHClients[host], command, output))
	}

	return NewParallel(tasks...)
}
