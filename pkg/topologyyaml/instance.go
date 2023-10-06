package topologyyaml

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql" // register mysql
	"github.com/jmoiron/sqlx"

	"github.com/lobshunter/dorisctl/pkg/embed"
)

var (
	_ Instance = &FeInstance{}
	_ Instance = &BeInstance{}
)

type Instance interface {
	EnvironmentVars() map[string]string
	PIDFile() string
	StartupScript() string
	StopScript() string
	SystemdEnvironment() string
	SystemdServiceContent() (string, error)
	SystemdServicePath() string
	SystemdServiceName() string
}

type FeInstance struct {
	DeployUser string
	Topo       *Topology
	FESpec
}

func NewFeInstance(topo *Topology, feSpec FESpec) FeInstance {
	return FeInstance{
		DeployUser: topo.Global.DeployUser,
		Topo:       topo,
		FESpec:     feSpec,
	}
}

func (s FeInstance) CheckHealth(ctx context.Context) (bool, error) {
	db, err := s.GetDB(ctx)
	if err != nil {
		return false, err
	}
	defer db.Close()

	err = db.PingContext(ctx)
	return err == nil, nil
}

func (s FeInstance) GetCluserStatus(ctx context.Context) (ClusterStatus, error) {
	clusterStatus := ClusterStatus{}

	db, err := s.GetDB(ctx)
	if err != nil {
		return clusterStatus, err
	}
	db = db.Unsafe()
	defer db.Close()

	clusterStatus.FeMasterHealthy = true
	// show frontends
	feRows, err := db.QueryxContext(ctx, "show frontends")
	if err != nil {
		return clusterStatus, fmt.Errorf("failed to query frontends: %v", err)
	}
	defer feRows.Close()
	for feRows.Next() {
		feStatus := FeStatus{}
		err = feRows.StructScan(&feStatus)
		if err != nil {
			return clusterStatus, fmt.Errorf("failed to scan frontend status: %v", err)
		}

		clusterStatus.Fes = append(clusterStatus.Fes, feStatus)
	}

	beRows, err := db.QueryxContext(ctx, "show backends")
	if err != nil {
		return clusterStatus, fmt.Errorf("failed to query backends: %v", err)
	}
	defer beRows.Close()
	for beRows.Next() {
		beStatus := BeStatus{}
		err = beRows.StructScan(&beStatus)
		if err != nil {
			return clusterStatus, fmt.Errorf("failed to scan backend status: %v", err)
		}

		clusterStatus.Bes = append(clusterStatus.Bes, beStatus)
	}

	return clusterStatus, nil
}

func (s FeInstance) CheckClusterHealth(ctx context.Context) (bool, error) {
	clusterStatus, err := s.GetCluserStatus(ctx)
	if err != nil {
		return false, err
	}

	for _, fe := range clusterStatus.Fes {
		if !fe.Alive {
			return false, nil
		}
	}

	for _, be := range clusterStatus.Bes {
		if !be.Alive {
			return false, nil
		}
	}

	return true, nil
}

func (s FeInstance) GetDB(_ context.Context) (*sqlx.DB, error) {
	// TODO: support username/password
	db, err := sqlx.Open("mysql",
		fmt.Sprintf("root:@tcp(%s:%d)/?charset=utf8mb4&timeout=3s", s.Host, s.FeConfig.QueryPort))
	if err != nil {
		return nil, fmt.Errorf("open mysql failed: %v", err)
	}

	return db, nil
}

// MasterHelper returns master_fe:master_editlog_port(bdbje helper), used by fe startup script
// See: https://doris.apache.org/docs/dev/admin-manual/maint-monitor/metadata-operation/
func (s FeInstance) MasterHelper() string {
	for _, fe := range s.Topo.FEs {
		if fe.IsMaster {
			return fmt.Sprintf("%s:%d", fe.Host, fe.FeConfig.EditLogPort)
		}
	}

	// use first fe as default master fe, but it's not supposed to so
	return fmt.Sprintf("%s:%d", s.Topo.FEs[0].Host, s.Topo.FEs[0].FeConfig.EditLogPort)
}

func (s FeInstance) ConfigDir() string {
	return filepath.Join(s.DeployDir, "conf")
}

func (s FeInstance) PIDFile() string {
	return filepath.Join(s.DeployDir, "bin", "fe.pid")
}

func (s FeInstance) StartupScript() string {
	startupScript := filepath.Join(s.DeployDir, "bin", "start_fe.sh")
	if !s.IsMaster {
		startupScript = fmt.Sprintf("%s --helper %s", startupScript, s.MasterHelper())
	}

	return startupScript
}

func (s FeInstance) StopScript() string {
	return filepath.Join(s.DeployDir, "bin", "stop_fe.sh")
}

func (s FeInstance) EnvironmentVars() map[string]string {
	envs := make(map[string]string)
	if s.InstallJava {
		envs["JAVA_HOME"] = fmt.Sprintf("%s/jdk", s.DeployDir)
	}

	return envs
}

func (s FeInstance) SystemdEnvironment() string {
	envs := []string{}
	for k, v := range s.EnvironmentVars() {
		envs = append(envs, fmt.Sprintf("Environment=%s=%s", k, v))
	}

	return strings.Join(envs, "\n")
}

func (s FeInstance) SystemdServiceContent() (string, error) {
	tpl, err := embed.ReadTemplate(embed.SystemdServicePath)
	if err != nil {
		return "", err
	}

	return executeTemplate(tpl, s)
}

func (s FeInstance) SystemdServicePath() string {
	return fmt.Sprintf("/etc/systemd/system/%s.service", s.SystemdServiceName())
}

func (s FeInstance) SystemdServiceName() string {
	return fmt.Sprintf("doris-fe-%d", s.FeConfig.EditLogPort)
}

type BeInstance struct {
	DeployUser string
	Topo       *Topology
	BESpec
}

func NewBeInstance(topo *Topology, beSpec BESpec) BeInstance {
	return BeInstance{
		DeployUser: topo.Global.DeployUser,
		Topo:       topo,
		BESpec:     beSpec,
	}
}

func (s BeInstance) ConfigDir() string {
	return filepath.Join(s.DeployDir, "conf")
}

func (s BeInstance) PIDFile() string {
	return filepath.Join(s.DeployDir, "bin", "be.pid")
}

func (s BeInstance) StartupScript() string {
	return filepath.Join(s.DeployDir, "bin", "start_be.sh")
}

func (s BeInstance) StopScript() string {
	return filepath.Join(s.DeployDir, "bin", "stop_be.sh")
}

func (s BeInstance) EnvironmentVars() map[string]string {
	envs := make(map[string]string)
	if s.InstallJava {
		envs["JAVA_HOME"] = fmt.Sprintf("%s/jdk", s.DeployDir)
	}

	return envs
}

func (s BeInstance) SystemdEnvironment() string {
	envs := []string{}
	for k, v := range s.EnvironmentVars() {
		envs = append(envs, fmt.Sprintf("Environment=%s=%s", k, v))
	}

	return strings.Join(envs, "\n")
}

func (s BeInstance) SystemdServiceContent() (string, error) {
	tpl, err := embed.ReadTemplate(embed.SystemdServicePath)
	if err != nil {
		return "", err
	}

	return executeTemplate(tpl, s)
}

func (s BeInstance) SystemdServicePath() string {
	return fmt.Sprintf("/etc/systemd/system/%s.service", s.SystemdServiceName())
}

func (s BeInstance) SystemdServiceName() string {
	return fmt.Sprintf("doris-be-%d", s.BeConfig.BePort)
}

func executeTemplate(tpl *template.Template, data interface{}) (string, error) {
	buff := bytes.NewBuffer(nil)
	err := tpl.Execute(buff, data)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}
