package doris

import (
	"context"
	"fmt"
	"time"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

var (
	_ task.Task = &WaitClusterHealthy{}
)

type WaitClusterHealthy struct {
	feMaster       topologyyaml.FeInstance
	defaultTimeout time.Duration
}

func NewWaitClusterHealthy(feMaster topologyyaml.FeInstance) *WaitClusterHealthy {
	return &WaitClusterHealthy{
		feMaster:       feMaster,
		defaultTimeout: 60 * time.Second,
	}
}

func (t *WaitClusterHealthy) Name() string {
	return "wait-cluster-healthy"
}

func (t *WaitClusterHealthy) Execute(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, t.defaultTimeout)
	defer cancel()

	checkInterval := 3 * time.Second
	ticker := time.NewTicker(checkInterval)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait cluster healthy ready timeout")
		case <-ticker.C:
			healthy, err := t.feMaster.CheckClusterHealth(ctx)
			if err != nil {
				return err
			}
			if healthy {
				return nil
			}
		}
	}
}
