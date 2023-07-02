package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/lobshunter/dorisctl/pkg/topologyyaml"
	"github.com/lobshunter/dorisctl/pkg/utils/filenames"
)

var (
	_ ManifestStore = &LocalManifestStore{}

	ErrManifestNotFound = errors.New("manifest not found")
)

// ManifestStore saves cluster metadata, such as topology.yaml
type ManifestStore interface { // Use an interface so we can potentially support remote storage in the future
	Save(clusterName string, topo topologyyaml.Topology) error
	Get(clusterName string) (*topologyyaml.Topology, error)
	Remove(clusterName string) error
	// List returns a list of cluster names
	List() ([]string, error)
}

type LocalManifestStore struct {
	StorePath string
}

func NewLocalManifestStore(storePath string) ManifestStore {
	return &LocalManifestStore{
		StorePath: storePath,
	}
}

func (s *LocalManifestStore) Save(clusterName string, topo topologyyaml.Topology) error {
	// TODO: might need filelock?
	fileDir := filepath.Join(s.StorePath, clusterName)
	filename := filepath.Join(fileDir, filenames.TopologyYamlFile)
	tmpFilename := filename + ".tmp"

	err := os.MkdirAll(fileDir, 0755)
	if err != nil {
		return err
	}

	topoData, err := yaml.Marshal(topo)
	if err != nil {
		return err
	}

	err = os.WriteFile(tmpFilename, topoData, 0644)
	if err != nil {
		return err
	}

	return os.Rename(tmpFilename, filename)
}

func (s *LocalManifestStore) Get(clusterName string) (*topologyyaml.Topology, error) {
	fileDir := filepath.Join(s.StorePath, clusterName)
	filename := filepath.Join(fileDir, filenames.TopologyYamlFile)

	topoData, err := os.ReadFile(filename)
	switch {
	case err == nil:
	case os.IsNotExist(err):
		return nil, ErrManifestNotFound
	default:
		return nil, err
	}

	var topo topologyyaml.Topology
	err = yaml.Unmarshal(topoData, &topo)
	if err != nil {
		return nil, err
	}

	return &topo, nil
}

func (s *LocalManifestStore) Remove(clusterName string) error {
	fileDir := filepath.Join(s.StorePath, clusterName)
	return os.RemoveAll(fileDir)
}

func (s *LocalManifestStore) List() ([]string, error) {
	_, err := os.Stat(s.StorePath)
	switch {
	case err == nil:
	case os.IsNotExist(err):
		return nil, nil
	default:
		return nil, fmt.Errorf("stat %s failed: %v", s.StorePath, err)
	}

	clusterNames := make([]string, 0)
	childDepth := strings.Count(s.StorePath, string(os.PathSeparator)) + 1
	err = filepath.WalkDir(s.StorePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		switch {
		case strings.Count(path, string(os.PathSeparator)) < childDepth:
			return nil // parent dir
		case strings.Count(path, string(os.PathSeparator)) == childDepth:
			clusterNames = append(clusterNames, d.Name()) // child dir
		case strings.Count(path, string(os.PathSeparator)) > childDepth:
			return fs.SkipDir // grand child dir, skip
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return clusterNames, nil
}
