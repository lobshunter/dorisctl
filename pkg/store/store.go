package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencontainers/go-digest"

	"github.com/lobshunter/dorisctl/pkg/utils/filenames"
)

// PackageStore saves downloaded packages locally
type PackageStore struct {
	storePath string
}

func NewPackageStore(storePath string) *PackageStore {
	return &PackageStore{
		storePath: storePath,
	}
}

func (s *PackageStore) Download(_ context.Context, srcURL string, digest digest.Digest) error {
	destDir := filepath.Join(s.storePath, "download", digest.String())
	if cacheExist(destDir) {
		return nil
	}

	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create cache dir %s: %v", destDir, err)
	}

	err = Lock(destDir)
	if err != nil {
		return fmt.Errorf("failed to lock cache dir %s: %v", destDir, err)
	}
	defer UnLock(destDir) //nolint:errcheck

	destFile := filepath.Join(destDir, filenames.CacheDataFile)
	if isLocal(srcURL) {
		err = localCopy(srcURL, destFile)
	} else {
		err = remoteDownload(srcURL, destFile)
	}
	if err != nil {
		return err
	}

	// TODO: run checksum
	urlFile, err := os.Create(filepath.Join(destDir, filenames.CacheURLFile))
	if err != nil {
		return err
	}
	defer urlFile.Close()
	_, err = urlFile.Write([]byte(srcURL))
	if err != nil {
		return err
	}

	digestFile, err := os.Create(filepath.Join(destDir, filenames.CacheDigestFile))
	if err != nil {
		return err
	}
	defer digestFile.Close()
	_, err = digestFile.Write([]byte(digest.String()))
	return err
}

func (s *PackageStore) GetLocalCachePath(digest digest.Digest) (string, error) {
	destDir := filepath.Join(s.storePath, "download", digest.String())
	if !cacheExist(destDir) {
		return "", fmt.Errorf("cache not exist")
	}

	return filepath.Join(destDir, filenames.CacheDataFile), nil
}
