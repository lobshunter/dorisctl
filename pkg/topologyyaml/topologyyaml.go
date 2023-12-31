package topologyyaml

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/opencontainers/go-digest"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
)

type ComponentType string

const (
	FE  ComponentType = "fe"
	BE  ComponentType = "be"
	JDK ComponentType = "jdk"
)

type Topology struct {
	Global GlobalSpec `yaml:"global"`
	// TODO: maybe some common config for all fe/bes, to shorten the topology yaml
	FEs []FESpec `yaml:"fes"`
	BEs []BESpec `yaml:"bes"`
}

func (t *Topology) UseSystemd() bool {
	if t.Global.UseSystemd == nil {
		return true
	}
	return *t.Global.UseSystemd
}

type GlobalSpec struct {
	DeployUser        string `yaml:"deploy_user"`
	SSHPrivateKeyPath string `yaml:"ssh_private_key_path"`
	UseSystemd        *bool  `yaml:"use_systemd"`
}

type FESpec struct {
	ComponentSpec  `yaml:",inline"`
	FeConfig       FeConfig      `yaml:"-"`
	IsMaster       bool          `yaml:"is_master"`
	InstallJava    bool          `yaml:"install_java"`
	JavaPackageURL string        `yaml:"java_package_url"`
	JavaDigest     digest.Digest `yaml:"java_digest"`
	PackageURL     string        `yaml:"package_url"`
	Digest         digest.Digest `yaml:"digest"`
}

type FeConfig struct {
	HTTPPort    int `ini:"http_port"`
	RPCPort     int `ini:"rpc_port"`
	QueryPort   int `ini:"query_port"`
	EditLogPort int `ini:"edit_log_port"`
}

type BESpec struct {
	ComponentSpec  `yaml:",inline"`
	BeConfig       BeConfig      `yaml:"-"`
	InstallJava    bool          `yaml:"install_java"`
	JavaPackageURL string        `yaml:"java_package_url"`
	JavaDigest     digest.Digest `yaml:"java_digest"`
	PackageURL     string        `yaml:"package_url"`
	Digest         digest.Digest `yaml:"digest"`
}

type BeConfig struct {
	BePort               int    `ini:"be_port"`
	WebServerPort        int    `ini:"webserver_port"`
	HeartbeatServicePort int    `ini:"heartbeat_service_port"`
	StorageRootPath      string `ini:"storage_root_path"`
}

type ComponentSpec struct {
	Host      string `yaml:"host"`
	SSHPort   int    `yaml:"ssh_port"`
	DeployDir string `yaml:"deploy_dir"`
	Config    string `yaml:"config"`
}

var _ yaml.Unmarshaler = (*FESpec)(nil)
var _ yaml.Unmarshaler = (*BESpec)(nil)

func (s *FESpec) UnmarshalYAML(value *yaml.Node) error {
	type feSpec FESpec
	spec := feSpec{}
	err := value.Decode(&spec)
	if err != nil {
		return err
	}

	*s = FESpec(spec)

	// cutsom unmarshal to parse config(in .ini format)
	if s.Config == "" {
		return nil
	}

	configData := bytes.NewBufferString(s.Config)
	iniData, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, configData)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	err = iniData.MapTo(&s.FeConfig)
	if err != nil {
		return fmt.Errorf("failed to map config to struct: %w", err)
	}

	return nil
}

func (s *BESpec) UnmarshalYAML(value *yaml.Node) error {
	type beSpec BESpec
	spec := beSpec{}
	err := value.Decode(&spec)
	if err != nil {
		return err
	}

	*s = BESpec(spec)

	// cutsom unmarshal to parse config(in .ini format)
	if s.Config == "" {
		return nil
	}

	configData := bytes.NewBufferString(s.Config)
	iniData, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, configData)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	err = iniData.MapTo(&s.BeConfig)
	if err != nil {
		return fmt.Errorf("failed to map config to struct: %w", err)
	}

	return nil
}

func LoadFromFile(filename string) (*Topology, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(f)
}

func Load(r io.Reader) (*Topology, error) {
	t := &Topology{}
	err := yaml.NewDecoder(r).Decode(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

type ManualDeployInfo struct {
	DeployUser string
	SSHPort    int
	SSHKeyPath string

	FeMaster    string
	FeHosts     []string
	BeHosts     []string
	FeDeployDir string
	BeDeployDir string
}

// BuildTopoFromManualDeploy transform deploy information into (incomplete) Topology struct
// It then be used to complete and generate a topology yaml file
func BuildTopoFromManualDeploy(deployInfo ManualDeployInfo) *Topology {
	useSystemd := false
	topo := &Topology{
		Global: GlobalSpec{
			DeployUser:        deployInfo.DeployUser,
			SSHPrivateKeyPath: deployInfo.SSHKeyPath,
			UseSystemd:        &useSystemd,
		},
		FEs: make([]FESpec, 0, len(deployInfo.FeHosts)),
		BEs: make([]BESpec, 0, len(deployInfo.BeHosts)),
	}

	for _, host := range deployInfo.FeHosts {
		topo.FEs = append(topo.FEs, FESpec{
			IsMaster: host == deployInfo.FeMaster,
			ComponentSpec: ComponentSpec{
				Host:      host,
				SSHPort:   deployInfo.SSHPort,
				DeployDir: deployInfo.FeDeployDir,
			},
		})
	}

	for _, host := range deployInfo.BeHosts {
		topo.BEs = append(topo.BEs, BESpec{
			ComponentSpec: ComponentSpec{
				Host:      host,
				SSHPort:   deployInfo.SSHPort,
				DeployDir: deployInfo.BeDeployDir,
			},
		})
	}

	return topo
}
