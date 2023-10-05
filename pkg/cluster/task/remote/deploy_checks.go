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
	_ task.Task = &CheckNoFile{}
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
		return fmt.Errorf("max_map_count=%d on host %s is too small, please run 'sysctl -w vm.max_map_count=2000000'", maxMapCount, t.sshClient.RemoveAddr())
	}

	return nil
}

type CheckNoFile struct {
	sshClient *ssh.SSHClient
}

func NewCheckNoFile(sshClient *ssh.SSHClient) *CheckNoFile {
	return &CheckNoFile{
		sshClient: sshClient,
	}
}

func (t *CheckNoFile) Name() string {
	return "check-no-file"
}

func (t *CheckNoFile) Execute(ctx context.Context) error {
	noFileString, err := t.sshClient.Exec(ctx, "ulimit -n")
	if err != nil {
		return err
	}

	noFile, err := strconv.Atoi(noFileString)
	if err != nil {
		return err
	}

	if noFile < 65536 {
		return fmt.Errorf("ulimit -n=%d on host %s is too small, please run 'ulimit -n 65536'", noFile, t.sshClient.RemoveAddr())
	}

	return nil
}
