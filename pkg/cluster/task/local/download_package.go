package local

import (
	"context"
	"fmt"

	"github.com/opencontainers/go-digest"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
	"github.com/lobshunter/dorisctl/pkg/store"
)

var (
	_ task.Task = &DownloadPackage{}
)

type DownloadPackage struct {
	store  *store.PackageStore
	url    string
	digest digest.Digest
}

func NewDownloadPackage(store *store.PackageStore, packageURL string, digest digest.Digest) *DownloadPackage {
	return &DownloadPackage{
		store:  store,
		url:    packageURL,
		digest: digest,
	}
}

func (t *DownloadPackage) Name() string {
	return "download-package"
}

func (t *DownloadPackage) Execute(ctx context.Context) error {
	err := t.store.Download(ctx, t.url, t.digest)
	if err != nil {
		return fmt.Errorf("download package %s failed: %v", t.url, err)
	}

	return nil
}
