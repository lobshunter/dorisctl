package remote

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencontainers/go-digest"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/store"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
	"github.com/lobshunter/dorisctl/pkg/utils/filenames"
	"github.com/lobshunter/dorisctl/pkg/utils/ssh"
)

var (
	_ task.Task = &InstallPackage{}
)

type InstallPackage struct {
	sshClient     *ssh.SSHClient
	store         *store.PackageStore
	packageDigest digest.Digest
	deployDir     string
	packageType   topologyyaml.ComponentType
}

func NewInstallPackage(sshClient *ssh.SSHClient, store *store.PackageStore, packageDigest digest.Digest, deployDir string, packageType topologyyaml.ComponentType) *InstallPackage {
	return &InstallPackage{
		sshClient:     sshClient,
		store:         store,
		packageDigest: packageDigest,
		deployDir:     deployDir,
		packageType:   packageType,
	}
}

func (t *InstallPackage) Name() string {
	return "install-package"
}

func (t *InstallPackage) Execute(ctx context.Context) error {
	localPackagePath, err := t.store.GetLocalCachePath(t.packageDigest)
	if err != nil {
		return err
	}
	localFile, err := os.Open(localPackagePath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	var packagePath string
	var untarCmds []string // TODO: unify doris package compression format, or detect its format
	switch t.packageType {
	case topologyyaml.BE:
		packagePath = filepath.Join(t.deployDir, filenames.DorisBEPackage)
		tmpPath := filepath.Join(t.deployDir, "tmp")
		// NOTE: it's extremely annoying that doris package doesn't have a unified structure
		// it can be either(both be and fe are in the same package):
		// apache-doris-1.2.4.1-bin-aarch64
		// ├── be
		// │   └── be contents...
		// └── fe
		//     └── fe contents...
		// or (be and fe are in different packages):
		// apache-doris-be-1.2.3-bin-arm
		// └── be contents...
		//
		// so we need to handle both cases
		untarCmds = append(untarCmds,
			fmt.Sprintf("mkdir -p %s", tmpPath),
			fmt.Sprintf("tar --no-same-owner --strip-components 1 -xf %s -C %s", packagePath, tmpPath),
			fmt.Sprintf("if [ ! -d %s/be ]; then mv %s/* %s; fi", tmpPath, tmpPath, t.deployDir),
			fmt.Sprintf("if [ -d %s/be ]; then mv %s/be/* %s; fi", tmpPath, tmpPath, t.deployDir),
			fmt.Sprintf("rm -rf %s", tmpPath),
		)
	case topologyyaml.FE:
		packagePath = filepath.Join(t.deployDir, filenames.DorisFEPackage)
		tmpPath := filepath.Join(t.deployDir, "tmp")
		untarCmds = append(untarCmds,
			fmt.Sprintf("mkdir -p %s", tmpPath),
			fmt.Sprintf("tar --no-same-owner --strip-components 1 -xf %s -C %s", packagePath, tmpPath),
			fmt.Sprintf("if [ ! -d %s/fe ]; then mv %s/* %s; fi", tmpPath, tmpPath, t.deployDir),
			fmt.Sprintf("if [ -d %s/fe ]; then mv %s/fe/* %s; fi", tmpPath, tmpPath, t.deployDir),
			fmt.Sprintf("rm -rf %s", tmpPath),
		)
	case topologyyaml.JDK:
		packagePath = filepath.Join(t.deployDir, filenames.JDKPackage)
		untarCmds = append(untarCmds, fmt.Sprintf("tar --no-same-owner --strip-components 1 -xzf %s -C %s", packagePath, t.deployDir))
	}

	err = t.sshClient.Run(ctx, fmt.Sprintf("mkdir -p %s", t.deployDir))
	if err != nil {
		return err
	}
	err = t.sshClient.WriteFile(localFile, packagePath)
	if err != nil {
		return err
	}

	return t.sshClient.Run(ctx, untarCmds...)
}
