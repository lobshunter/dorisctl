package doris

import (
	"context"
	"fmt"
	"time"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

var (
	_ task.Task = &WaitFeMasterReady{}
)

type WaitFeMasterReady struct {
	feMaster       topologyyaml.FeInstance
	defaultTimeout time.Duration
}

func NewWaitFeMasterReady(feMaster topologyyaml.FeInstance) *WaitFeMasterReady {
	return &WaitFeMasterReady{
		feMaster:       feMaster,
		defaultTimeout: 60 * time.Second,
	}
}

func (t *WaitFeMasterReady) Name() string {
	return "wait-fe-master-ready"
}

func (t *WaitFeMasterReady) Execute(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, t.defaultTimeout)
	defer cancel()

	checkInterval := 3 * time.Second
	ticker := time.NewTicker(checkInterval)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait fe-master ready timeout")
		case <-ticker.C:
			healthy, err := t.feMaster.CheckHealth(ctx)
			if err != nil {
				return err
			}
			if healthy {
				return nil
			}
		}
	}
}
