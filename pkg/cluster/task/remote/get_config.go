package remote

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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
	output = trimTrailSpaces(output)

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
	output = trimTrailSpaces(output)

	*t.output = output
	return nil
}

// trimTrailSpaces trims spaces at the end of each line.
// it's essential for "gopkg.in/yaml.v3".Marshal to marshal multiline string to human readable format.
// see: https://github.com/go-yaml/yaml/blob/496545a6307b2a7d7a710fd516e5e16e8ab62dbc/emitterc.go#L1392-L1394
func trimTrailSpaces(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}

	return strings.Join(lines, "\n")
}
