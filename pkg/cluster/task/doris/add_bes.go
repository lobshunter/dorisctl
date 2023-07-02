package doris

import (
	"context"
	"fmt"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
	"golang.org/x/exp/slices"
)

var (
	_ task.Task = &AddBes{}
)

type AddBes struct {
	taskCtx  *task.Context
	feMaster topologyyaml.FeInstance
	bes      []topologyyaml.BESpec
}

func NewAddBes(taskCtx *task.Context, feMaster topologyyaml.FeInstance, bes []topologyyaml.BESpec) *AddBes {
	return &AddBes{
		taskCtx:  taskCtx,
		feMaster: feMaster,
		bes:      bes,
	}
}

func (t *AddBes) Name() string {
	return "add-bes"
}

func (t *AddBes) Execute(ctx context.Context) error {
	beStatus := t.taskCtx.ClusterStatus.Bes

	db, err := t.feMaster.GetDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, be := range t.bes {
		// skip if already exists
		if slices.ContainsFunc(beStatus, func(beStat topologyyaml.BeStatus) bool {
			return beStat.Host == be.Host
		}) {
			continue
		}

		// add backend
		stmt := fmt.Sprintf(`ALTER SYSTEM ADD BACKEND "%s:%d"`, be.Host, be.BeConfig.HeartbeatServicePort)
		_, err = db.ExecContext(ctx, stmt)
		if err != nil {
			return fmt.Errorf("failed to add backend, sql %s, err: %v", stmt, err)
		}
	}

	return nil
}
