package step

import (
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/cluster/task/remote"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

func NewStopCluster(ctx *task.Context) Step {
	tasks := make([]task.Task, 0, len(ctx.TopoYaml.FEs)+len(ctx.TopoYaml.BEs))

	for _, fe := range ctx.TopoYaml.FEs {
		sshClient := ctx.SSHClients[fe.Host]
		feInstance := topologyyaml.NewFeInstance(&ctx.TopoYaml, fe)
		tasks = append(tasks, remote.NewStopService(sshClient, feInstance.SystemdServiceName()))
	}

	for _, be := range ctx.TopoYaml.BEs {
		sshClient := ctx.SSHClients[be.Host]
		beInstance := topologyyaml.NewBeInstance(&ctx.TopoYaml, be)
		tasks = append(tasks, remote.NewStopService(sshClient, beInstance.SystemdServiceName()))
	}

	return NewParallel(tasks...)
}
