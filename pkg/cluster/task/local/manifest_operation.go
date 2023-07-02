package local

import (
	"context"
	"fmt"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/store"
	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
)

var (
	_ task.Task = &SaveManifest{}
	_ task.Task = &RemoveManifest{}
	_ task.Task = &EnsureUniqueManifest{}
)

type EnsureUniqueManifest struct {
	store       store.ManifestStore
	clusterName string
}

func NewEnsureUniqueManifest(store store.ManifestStore, clusterName string) *EnsureUniqueManifest {
	return &EnsureUniqueManifest{
		store:       store,
		clusterName: clusterName,
	}
}

func (t *EnsureUniqueManifest) Name() string {
	return "ensure-unique-manifest"
}

func (t *EnsureUniqueManifest) Execute(_ context.Context) error {
	_, err := t.store.Get(t.clusterName)
	switch err {
	case nil:
		return fmt.Errorf("cluster [%s] already exists", t.clusterName)
	case store.ErrManifestNotFound:
		return nil
	default:
		return err
	}
}

type SaveManifest struct {
	store       store.ManifestStore
	clusterName string
	topo        *topologyyaml.Topology
}

func NewSaveManifest(store store.ManifestStore, clusterName string, topo *topologyyaml.Topology) *SaveManifest {
	return &SaveManifest{
		store:       store,
		clusterName: clusterName,
		topo:        topo,
	}
}

func (t *SaveManifest) Name() string {
	return "save-manifest"
}

func (t *SaveManifest) Execute(_ context.Context) error {
	err := t.store.Save(t.clusterName, *t.topo)
	if err != nil {
		return fmt.Errorf("save manifest %s failed: %v", t.clusterName, err)
	}

	return nil
}

type RemoveManifest struct {
	store       store.ManifestStore
	clusterName string
}

func NewRemoveManifest(store store.ManifestStore, clusterName string) *RemoveManifest {
	return &RemoveManifest{
		store:       store,
		clusterName: clusterName,
	}
}

func (t *RemoveManifest) Name() string {
	return "remove-manifest"
}

func (t *RemoveManifest) Execute(_ context.Context) error {
	err := t.store.Remove(t.clusterName)
	if err != nil {
		return fmt.Errorf("remove manifest %s failed: %v", t.clusterName, err)
	}

	return nil
}
