package doris

import (
	"context"
	"fmt"

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

	clusterStatus := topologyyaml.ClusterStatus{
		FeMasterHealthy: feMasterHealthy,
	}

	if !feMasterHealthy {
		return nil
	}

	db, err := t.feMaster.GetDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to fe-master: %v", err)
	}
	db = db.Unsafe()
	defer db.Close()

	// show frontends
	feRows, err := db.QueryxContext(ctx, "show frontends")
	if err != nil {
		return fmt.Errorf("failed to query frontends: %v", err)
	}
	defer feRows.Close()
	for feRows.Next() {
		feStatus := topologyyaml.FeStatus{}
		err = feRows.StructScan(&feStatus)
		if err != nil {
			return fmt.Errorf("failed to scan frontend status: %v", err)
		}

		clusterStatus.Fes = append(clusterStatus.Fes, feStatus)
	}

	beRows, err := db.QueryxContext(ctx, "show backends")
	if err != nil {
		return fmt.Errorf("failed to query backends: %v", err)
	}
	defer beRows.Close()
	for beRows.Next() {
		beStatus := topologyyaml.BeStatus{}
		err = beRows.StructScan(&beStatus)
		if err != nil {
			return fmt.Errorf("failed to scan backend status: %v", err)
		}

		clusterStatus.Bes = append(clusterStatus.Bes, beStatus)
	}

	t.taskContext.SetClusterStatus(clusterStatus)

	return nil
}
