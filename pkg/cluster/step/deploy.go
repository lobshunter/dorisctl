package step

import (
	"path/filepath"

	"github.com/opencontainers/go-digest"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/cluster/task/local"
	"github.com/lobshunter/dorisctl/pkg/cluster/task/remote"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
	"github.com/lobshunter/dorisctl/pkg/utils/filenames"
)

func NewUniqueCluster(ctx *task.Context) Step {
	return NewSerial(local.NewEnsureUniqueManifest(ctx.ManifestStore, ctx.ClusterName))
}

func NewDownloadPackage(ctx *task.Context) Step {
	pkgs := make(map[string]digest.Digest)
	for _, be := range ctx.TopoYaml.BEs {
		pkgs[be.PackageURL] = be.Digest
		if be.InstallJava {
			pkgs[be.JavaPackageURL] = be.JavaDigest
		}
	}
	for _, fe := range ctx.TopoYaml.FEs {
		pkgs[fe.PackageURL] = fe.Digest
		if fe.InstallJava {
			pkgs[fe.JavaPackageURL] = fe.JavaDigest
		}
	}

	tasks := make([]task.Task, 0, len(pkgs))
	for pkg, digest := range pkgs {
		tasks = append(tasks, local.NewDownloadPackage(ctx.PackageStore, pkg, digest))
	}

	return NewParallel(tasks...)
}

func NewInitHost(ctx *task.Context) Step {
	initUniqueHostTasks := make([]task.Task, 0)
	uniqueHosts := make(map[string]struct{})
	for _, fe := range ctx.TopoYaml.FEs {
		uniqueHosts[fe.Host] = struct{}{}
	}
	for _, be := range ctx.TopoYaml.BEs {
		uniqueHosts[be.Host] = struct{}{}
	}
	for host := range uniqueHosts {
		sshClient := ctx.SSHClients[host]
		initUniqueHostTasks = append(initUniqueHostTasks, remote.NewInitInstanceUniqueHost(sshClient, ctx.TopoYaml.Global.DeployUser))
	}

	initInstanceHostTasks := make([]task.Task, 0, len(ctx.TopoYaml.FEs)+len(ctx.TopoYaml.BEs))
	// maybe split into two tasks: component dedicated & component agnositic
	for _, fe := range ctx.TopoYaml.FEs {
		sshClient := ctx.SSHClients[fe.Host]
		initInstanceHostTasks = append(initInstanceHostTasks, remote.NewInitInstanceHost(sshClient, ctx.TopoYaml.Global.DeployUser, fe.DeployDir))
	}
	for _, be := range ctx.TopoYaml.BEs {
		sshClient := ctx.SSHClients[be.Host]
		initInstanceHostTasks = append(initInstanceHostTasks,
			remote.NewInitInstanceHost(sshClient, ctx.TopoYaml.Global.DeployUser, be.DeployDir))
	}

	initBeTasks := make([]task.Task, 0, len(ctx.TopoYaml.BEs))
	for _, be := range ctx.TopoYaml.BEs {
		sshClient := ctx.SSHClients[be.Host]
		initBeTasks = append(initBeTasks,
			remote.NewInitBe(sshClient, ctx.TopoYaml.Global.DeployUser,
				topologyyaml.NewBeInstance(&ctx.TopoYaml, be)))
	}

	return NewSerial(
		NewParallel(initUniqueHostTasks...),
		NewParallel(initInstanceHostTasks...),
		NewParallel(initBeTasks...),
	)
}

func NewInstallPackage(ctx *task.Context) Step {
	tasks := make([]task.Task, 0, len(ctx.TopoYaml.FEs)+len(ctx.TopoYaml.BEs))

	for _, fe := range ctx.TopoYaml.FEs {
		sshClient := ctx.SSHClients[fe.Host]
		feDir := fe.DeployDir
		javaDir := filepath.Join(feDir, filenames.JDKDir)
		if fe.InstallJava {
			tasks = append(tasks, NewSerial(
				remote.NewInstallPackage(sshClient, ctx.PackageStore, fe.Digest, feDir, topologyyaml.FE),
				remote.NewInstallPackage(sshClient, ctx.PackageStore, fe.JavaDigest, javaDir, topologyyaml.JDK),
			))
		} else {
			tasks = append(tasks, NewSerial(
				remote.NewInstallPackage(sshClient, ctx.PackageStore, fe.Digest, feDir, topologyyaml.FE),
			))
		}

	}
	for _, be := range ctx.TopoYaml.BEs {
		sshClient := ctx.SSHClients[be.Host]
		beDir := be.DeployDir
		javaDir := filepath.Join(beDir, filenames.JDKDir)
		if be.InstallJava {
			tasks = append(tasks, NewSerial(
				remote.NewInstallPackage(sshClient, ctx.PackageStore, be.Digest, beDir, topologyyaml.BE),
				remote.NewInstallPackage(sshClient, ctx.PackageStore, be.JavaDigest, javaDir, topologyyaml.JDK),
			))
		} else {
			tasks = append(tasks, NewSerial(
				remote.NewInstallPackage(sshClient, ctx.PackageStore, be.Digest, beDir, topologyyaml.BE),
			))
		}
	}

	return NewParallel(tasks...)
}

func NewInitConfig(ctx *task.Context) Step {
	tasks := make([]task.Task, 0, len(ctx.TopoYaml.FEs)+len(ctx.TopoYaml.BEs))

	for _, fe := range ctx.TopoYaml.FEs {
		sshClient := ctx.SSHClients[fe.Host]
		tasks = append(tasks, remote.NewInitFeConfig(sshClient, topologyyaml.NewFeInstance(&ctx.TopoYaml, fe), fe.DeployDir))
	}

	for _, be := range ctx.TopoYaml.BEs {
		sshClient := ctx.SSHClients[be.Host]
		tasks = append(tasks, remote.NewInitBeConfig(sshClient, topologyyaml.NewBeInstance(&ctx.TopoYaml, be), be.DeployDir))
	}

	return NewParallel(tasks...)
}

func NewSaveManifest(ctx *task.Context) Step {
	task := local.NewSaveManifest(ctx.ManifestStore, ctx.ClusterName, &ctx.TopoYaml)
	return NewSerial(task)
}
