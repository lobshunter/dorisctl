package cluster

import (
	"context"
	"fmt"

	"github.com/lobshunter/dorisctl/pkg/cluster/step"
	"github.com/lobshunter/dorisctl/pkg/cluster/task"
)

func RunSteps(ctx context.Context, steps ...step.Step) error {
	var err error
	for _, t := range steps {
		err = t.Execute(ctx)
		if err != nil {
			return fmt.Errorf("task %s failed: %v", t.Name(), err)
		}
	}

	return nil
}

func GetDeployClusterSteps(ctx *task.Context) []step.Step {
	return []step.Step{
		// TODO: maybe lock manifest before finish
		step.NewStepDisplay("Ensure Unique Cluster ", step.NewUniqueCluster(ctx)),
		step.NewStepDisplay("Check Hosts Config ", step.NewCheckHostsConfig(ctx)),
		step.NewStepDisplay("Download Packages ", step.NewDownloadPackage(ctx)),
		step.NewStepDisplay("Init Hosts ", step.NewInitHost(ctx)),
		step.NewStepDisplay("Install Packages ", step.NewInstallPackage(ctx)),
		step.NewStepDisplay("Init Configs ", step.NewInitConfig(ctx)),
		step.NewStepDisplay("Save Manifest ", step.NewSaveManifest(ctx)),
	}
}

func GetStartClusterSteps(ctx *task.Context) []step.Step {
	return []step.Step{
		step.NewStepDisplay("Start Cluster ", step.NewStartCluster(ctx)),
		step.NewStepDisplay("Add Bes ", step.NewAddBes(ctx)),
	}
}

func GetStopClusterSteps(ctx *task.Context) []step.Step {
	return []step.Step{
		step.NewStepDisplay("Stop Cluster ", step.NewStopCluster(ctx)),
	}
}

func GetDestroyClusterSteps(ctx *task.Context) []step.Step {
	stopSteps := GetStopClusterSteps(ctx)
	destroySteps := []step.Step{
		step.NewStepDisplay("Reset Hosts ", step.NewResetHost(ctx)),
		step.NewStepDisplay("Remove Manifest ", step.NewRemoveManifest(ctx)),
	}
	return append(stopSteps, destroySteps...)
}

func GetCheckClusterStatusSteps(ctx *task.Context) []step.Step {
	return []step.Step{
		step.NewStepDisplay("Check Cluster Status ", step.NewCheckClusterStatus(ctx)),
	}
}
