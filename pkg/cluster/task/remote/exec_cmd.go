package remote

import (
	"context"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &ExecCmd{}
)

type ExecCmd struct {
	sshClient *ssh.SSHClient
	command   string

	output *topologyyaml.ExecOutput
}

func NewExecCmd(sshClient *ssh.SSHClient, command string, output *topologyyaml.ExecOutput) *ExecCmd {
	return &ExecCmd{
		sshClient: sshClient,
		command:   command,
		output:    output,
	}
}

func (t *ExecCmd) Name() string {
	return "ExecCmd"
}

func (t *ExecCmd) Execute(ctx context.Context) error {
	stdout, stderr, exitCode, err := t.sshClient.ExecVerbose(ctx, t.command)
	if err != nil {
		return err
	}

	t.output.Stdout = stdout
	t.output.Stderr = stderr
	t.output.ExitCode = exitCode
	return nil
}
