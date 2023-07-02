package remote

import (
	"bytes"
	"context"
	"path/filepath"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
	"github.com/lobshunter/dorisctl/pkg/utils/filenames"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &InitFeConfig{}
	_ task.Task = &InitBeConfig{}
)

type InitFeConfig struct {
	sshClient *ssh.SSHClient
	fe        topologyyaml.FeInstance
	deployDir string
}

func NewInitFeConfig(sshClient *ssh.SSHClient, fe topologyyaml.FeInstance, deployDir string) *InitFeConfig {
	return &InitFeConfig{
		sshClient: sshClient,
		fe:        fe,
		deployDir: deployDir,
	}
}

func (t *InitFeConfig) Name() string {
	return "init-fe-config"
}

func (t *InitFeConfig) Execute(ctx context.Context) error {
	err := t.sshClient.Run(ctx, "mkdir -p "+t.fe.ConfigDir())
	if err != nil {
		return err
	}
	err = t.sshClient.WriteFile(bytes.NewBufferString(t.fe.Config),
		filepath.Join(t.fe.ConfigDir(), filenames.DorisFEConfFile))
	if err != nil {
		return err
	}
	startupScript, err := t.fe.StarupScript()
	if err != nil {
		return err
	}
	err = t.sshClient.WriteFile(bytes.NewBufferString(startupScript), t.fe.StarupScriptPath())
	if err != nil {
		return err
	}
	systemdUnit, err := t.fe.SystemdServiceContent()
	if err != nil {
		return err
	}
	err = t.sshClient.WriteFile(bytes.NewBufferString(systemdUnit), t.fe.SystemdServicePath())
	if err != nil {
		return err
	}

	err = t.sshClient.Run(ctx, "chown -R "+t.fe.DeployUser+" "+t.fe.DeployDir)
	if err != nil {
		return err
	}

	return nil
}

type InitBeConfig struct {
	sshClient *ssh.SSHClient
	be        topologyyaml.BeInstance
	deployDir string
}

func NewInitBeConfig(sshClient *ssh.SSHClient, be topologyyaml.BeInstance, deployDir string) *InitBeConfig {
	return &InitBeConfig{
		sshClient: sshClient,
		be:        be,
		deployDir: deployDir,
	}
}

func (t *InitBeConfig) Name() string {
	return "init-be-config"
}

func (t *InitBeConfig) Execute(_ context.Context) error {
	err := t.sshClient.Run(context.Background(), "mkdir -p "+t.be.ConfigDir())
	if err != nil {
		return err
	}
	err = t.sshClient.WriteFile(bytes.NewBufferString(t.be.Config),
		filepath.Join(t.be.ConfigDir(), filenames.DorisBEConfFile))
	if err != nil {
		return err
	}
	startupScript, err := t.be.StarupScript()
	if err != nil {
		return err
	}
	err = t.sshClient.WriteFile(bytes.NewBufferString(startupScript), t.be.StarupScriptPath())
	if err != nil {
		return err
	}

	systemdUnit, err := t.be.SystemdServiceContent()
	if err != nil {
		return err
	}
	err = t.sshClient.WriteFile(bytes.NewBufferString(systemdUnit), t.be.SystemdServicePath())
	if err != nil {
		return err
	}

	err = t.sshClient.Run(context.Background(), "chown -R "+t.be.DeployUser+" "+t.be.DeployDir)
	if err != nil {
		return err
	}

	return nil
}
