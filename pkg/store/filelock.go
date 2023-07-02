package store

import (
	"os"
	"path/filepath"
)

func Lock(filename string) error {
	lockfile := filepath.Clean(filename) + ".lock"
	f, err := os.OpenFile(lockfile, os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}

	return f.Close()
}

func UnLock(filename string) error {
	lockfile := filepath.Clean(filename) + ".lock"
	return os.Remove(lockfile)
}
