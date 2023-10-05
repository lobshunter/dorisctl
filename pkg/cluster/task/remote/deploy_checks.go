package remote

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &CheckMaxMapCount{}
)

type CheckMaxMapCount struct {
	sshClient *ssh.SSHClient
}

func NewCheckMaxMapCount(sshClient *ssh.SSHClient) *CheckMaxMapCount {
	return &CheckMaxMapCount{
		sshClient: sshClient,
	}
}

func (t *CheckMaxMapCount) Name() string {
	return "check-max-map-count"
}

func (t *CheckMaxMapCount) Execute(ctx context.Context) error {
	stdout, err := t.sshClient.Exec(ctx, "uname -s")
	if err != nil {
		return err
	}

	if stdout == "Darwin" {
		return nil
	}

	maxMapCountString, err := t.sshClient.Exec(ctx, "cat /proc/sys/vm/max_map_count")
	if err != nil {
		return err
	}

	maxMapCount, err := strconv.Atoi(maxMapCountString)
	if err != nil {
		return err
	}

	if maxMapCount < 2000000 {
		return fmt.Errorf("max_map_count=%d which is too small, please run 'sysctl -w vm.max_map_count=2000000'", maxMapCount)
	}

	return nil
}
