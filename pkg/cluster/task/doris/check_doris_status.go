package doris

import (
	"context"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

var (
	_ task.Task = &CheckDorisNodesStatus{}
)

type CheckDorisNodesStatus struct {
	taskContext *task.Context
	feMaster    topologyyaml.FeInstance
}

func NewCheckDorisClusterStatus(ctx *task.Context, feMaster topologyyaml.FeInstance) *CheckDorisNodesStatus {
	return &CheckDorisNodesStatus{
		taskContext: ctx,
		feMaster:    feMaster,
	}
}

func (t *CheckDorisNodesStatus) Name() string {
	return "check-doris-nodes-status"
}

func (t *CheckDorisNodesStatus) Execute(ctx context.Context) error {
	feMasterHealthy, err := t.feMaster.CheckHealth(ctx)
	if err != nil {
		return err
	}
	if !feMasterHealthy {
		return nil
	}

	clusterStatus, err := t.feMaster.GetCluserStatus(ctx)
	if err != nil {
		return err
	}
	t.taskContext.SetClusterStatus(clusterStatus)

	return nil
}
