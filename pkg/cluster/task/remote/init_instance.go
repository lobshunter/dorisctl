package remote

import (
	"context"
	"strings"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &InitBe{}
)

type InitBe struct {
	sshClient  *ssh.SSHClient
	deployUser string
	beSpec     topologyyaml.BeInstance
}

func NewInitBe(sshClient *ssh.SSHClient, deployUser string, be topologyyaml.BeInstance) *InitBe {
	return &InitBe{
		sshClient:  sshClient,
		deployUser: deployUser,
		beSpec:     be,
	}
}

func (t *InitBe) Name() string {
	return "init-be"
}

func (t *InitBe) Execute(ctx context.Context) error {
	storageRootPaths := strings.Split(t.beSpec.BeConfig.StorageRootPath, ";")
	commands := make([]string, 0, len(storageRootPaths)*2)
	for _, storageRootPath := range storageRootPaths {
		commands = append(commands, "mkdir -p "+storageRootPath)
		commands = append(commands, "chown -R "+t.deployUser+" "+storageRootPath)
	}

	return t.sshClient.Run(ctx, commands...)
}
