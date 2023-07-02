package step

import (
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/cluster/task/local"
	"github.com/lobshunter/dorisctl/pkg/cluster/task/remote"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

func NewResetHost(ctx *task.Context) Step {
	tasks := make([]task.Task, 0, len(ctx.TopoYaml.FEs)+len(ctx.TopoYaml.BEs))

	uniqueHosts := make(map[string]struct{})
	for _, fe := range ctx.TopoYaml.FEs {
		uniqueHosts[fe.Host] = struct{}{}
	}
	for _, be := range ctx.TopoYaml.BEs {
		uniqueHosts[be.Host] = struct{}{}
	}
	for host := range uniqueHosts {
		sshClient := ctx.SSHClients[host]
		tasks = append(tasks, remote.NewResetInstanceUniqueHost(sshClient, ctx.TopoYaml.Global.DeployUser))
	}

	for _, fe := range ctx.TopoYaml.FEs {
		sshClient := ctx.SSHClients[fe.Host]
		feInstance := topologyyaml.NewFeInstance(&ctx.TopoYaml, fe)
		tasks = append(tasks,
			remote.NewResetInstanceHost(sshClient, ctx.TopoYaml.Global.DeployUser, fe.DeployDir, feInstance.SystemdServicePath()))
	}
	for _, be := range ctx.TopoYaml.BEs {
		sshClient := ctx.SSHClients[be.Host]
		beInstance := topologyyaml.NewBeInstance(&ctx.TopoYaml, be)
		tasks = append(tasks,
			remote.NewResetInstanceHost(sshClient, ctx.TopoYaml.Global.DeployUser, be.DeployDir, beInstance.SystemdServicePath()))
	}

	return NewParallel(tasks...)
}

func NewRemoveManifest(ctx *task.Context) Step {
	task := local.NewRemoveManifest(ctx.ManifestStore, ctx.ClusterName)
	return NewSerial(task)
}
