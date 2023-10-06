package step

import (
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/cluster/task/local"
	"github.com/lobshunter/dorisctl/pkg/cluster/task/remote"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

// CompleteTopology retrives necessary info from cluster(via ssh and sql)
func NewCompleteTopology(ctx *task.Context) Step {
	// this step is really dirty:
	// - topo isn't a deep copy, so it's not clean to modify it
	// - a new SetTopoYaml task is added just to update the topo in taskCtx...

	topo := ctx.TopoYaml

	tasks := make([]task.Task, 0, len(ctx.TopoYaml.FEs)+len(ctx.TopoYaml.BEs))
	for i, fe := range ctx.TopoYaml.FEs {
		tasks = append(tasks, remote.NewGetRunningFeConfig(ctx.SSHClients[fe.Host], topologyyaml.NewFeInstance(&ctx.TopoYaml, fe), &topo.FEs[i].Config))
	}
	for i, be := range ctx.TopoYaml.BEs {
		tasks = append(tasks, remote.NewGetRunningBeConfig(ctx.SSHClients[be.Host], topologyyaml.NewBeInstance(&ctx.TopoYaml, be), &topo.BEs[i].Config))
	}

	tasks = append(tasks, local.NewSetTopoYaml(ctx, &topo))

	return NewParallel(tasks...)
}
