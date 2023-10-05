package step

import (
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/cluster/task/doris"
	"github.com/lobshunter/dorisctl/pkg/cluster/task/remote"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
	"golang.org/x/exp/slices"
)

func NewStartCluster(ctx *task.Context) Step {
	tasks := make([]task.Task, 0, len(ctx.TopoYaml.FEs)+len(ctx.TopoYaml.BEs))

	for _, fe := range ctx.TopoYaml.FEs {
		sshClient := ctx.SSHClients[fe.Host]
		feInstance := topologyyaml.NewFeInstance(&ctx.TopoYaml, fe)
		tasks = append(tasks, remote.NewStartService(sshClient, feInstance, ctx.TopoYaml.UseSystemd()))
	}

	for _, be := range ctx.TopoYaml.BEs {
		sshClient := ctx.SSHClients[be.Host]
		beInstance := topologyyaml.NewBeInstance(&ctx.TopoYaml, be)
		tasks = append(tasks, remote.NewStartService(sshClient, beInstance, ctx.TopoYaml.UseSystemd()))
	}

	return NewParallel(tasks...)
}

func NewAddBes(ctx *task.Context) Step {
	feMasterIdx := slices.IndexFunc(ctx.TopoYaml.FEs, func(fe topologyyaml.FESpec) bool {
		return fe.IsMaster
	})
	feMaster := topologyyaml.NewFeInstance(&ctx.TopoYaml, ctx.TopoYaml.FEs[feMasterIdx])

	return NewSerial(
		doris.NewWaitFeMasterReady(feMaster),
		doris.NewCheckDorisClusterStatus(ctx, feMaster),
		doris.NewAddBes(ctx, feMaster, ctx.TopoYaml.BEs),
	)
}

func NewCheckClusterStatus(ctx *task.Context) Step {
	feMasterIdx := slices.IndexFunc(ctx.TopoYaml.FEs, func(fe topologyyaml.FESpec) bool {
		return fe.IsMaster
	})
	feMaster := topologyyaml.NewFeInstance(&ctx.TopoYaml, ctx.TopoYaml.FEs[feMasterIdx])

	task := doris.NewCheckDorisClusterStatus(ctx, feMaster)

	return NewSerial(task)
}
