package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/lobshunter/dorisctl/pkg/utils/filenames"
)

var (
	GlobalConfig Config
)

type Config struct {
	CacheDir string
	DataDir  string
}

func MakeGlobalConfig() error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	cacheDir = filepath.Join(cacheDir, filenames.Dorisctl)

	dataDir, err := UserDataDir()
	if err != nil {
		return err
	}
	dataDir = filepath.Join(dataDir, filenames.Dorisctl)

	GlobalConfig = Config{
		CacheDir: cacheDir,
		DataDir:  dataDir,
	}
	return nil
}

func UserDataDir() (string, error) {
	dir := os.Getenv("XDG_DATA_HOME")
	if dir == "" {
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("neither $XDG_CACHE_HOME nor $HOME are defined")
		}
	}
	dir = filepath.Join(dir, ".local/share")
	return dir, nil
}
