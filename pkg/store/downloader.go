package store

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/lobshunter/dorisctl/pkg/utils/filenames"
)

func cacheExist(contentDir string) bool {
	_, err := os.Stat(filepath.Join(contentDir, filenames.CacheDigestFile))
	return !os.IsNotExist(err)
}

func isLocal(s string) bool {
	return !strings.Contains(s, "://") || strings.HasPrefix(s, "file://")
}

func localCopy(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func remoteDownload(src, dst string) error {
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	resp, err := http.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(dstFile, resp.Body)
	return err
}
