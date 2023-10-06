package local

import (
	"context"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

var (
	_ task.Task = &SetTopoYaml{}
)

type SetTopoYaml struct {
	taskCtx  *task.Context
	TopoYaml *topologyyaml.Topology
}

func NewSetTopoYaml(taskCtx *task.Context, topoYaml *topologyyaml.Topology) *SetTopoYaml {
	return &SetTopoYaml{
		taskCtx:  taskCtx,
		TopoYaml: topoYaml,
	}
}

func (t *SetTopoYaml) Name() string {
	return "set-topo-yaml"
}

func (t *SetTopoYaml) Execute(_ context.Context) error {
	t.taskCtx.SetTopoYaml(*t.TopoYaml)
	return nil
}
