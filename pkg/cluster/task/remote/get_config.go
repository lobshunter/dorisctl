package remote

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
	"github.com/lobshunter/dorisctl/pkg/utils/filenames"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &GetRunningFeConfig{}
)

type GetRunningFeConfig struct {
	sshClient *ssh.SSHClient
	fe        topologyyaml.FeInstance
	output    *string
}

func NewGetRunningFeConfig(sshClient *ssh.SSHClient, fe topologyyaml.FeInstance, output *string) *GetRunningFeConfig {
	return &GetRunningFeConfig{
		sshClient: sshClient,
		fe:        fe,
		output:    output,
	}
}

func (t *GetRunningFeConfig) Name() string {
	return "get-running-fe-config"
}

func (t *GetRunningFeConfig) Execute(ctx context.Context) error {
	command := fmt.Sprintf("cat %s", filepath.Join(t.fe.ConfigDir(), filenames.DorisFEConfFile))

	output, err := t.sshClient.Exec(ctx, command)
	if err != nil {
		return err
	}

	*t.output = output
	return nil
}

type GetRunningBeConfig struct {
	sshClient *ssh.SSHClient
	be        topologyyaml.BeInstance
	output    *string
}

func NewGetRunningBeConfig(sshClient *ssh.SSHClient, be topologyyaml.BeInstance, output *string) *GetRunningBeConfig {
	return &GetRunningBeConfig{
		sshClient: sshClient,
		be:        be,
		output:    output,
	}
}

func (t *GetRunningBeConfig) Name() string {
	return "get-running-be-config"
}

func (t *GetRunningBeConfig) Execute(ctx context.Context) error {
	command := fmt.Sprintf("cat %s", filepath.Join(t.be.ConfigDir(), filenames.DorisBEConfFile))

	output, err := t.sshClient.Exec(ctx, command)
	if err != nil {
		return err
	}

	*t.output = output
	return nil
}
