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

func (s FeInstance) StarupScript() string {
	startupScript := filepath.Join(s.DeployDir, "bin", "start_fe.sh")
	if !s.IsMaster {
		startupScript = fmt.Sprintf("%s --helper %s", startupScript, s.MasterHelper())
	}

	return startupScript
}

func (s FeInstance) SystemdEnvironment() string {
	envs := []string{}
	if s.InstallJava {
		envs = append(envs,
			fmt.Sprintf("Environment=JAVA_HOME=%s/jdk", s.DeployDir))
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

func (s BeInstance) StarupScript() string {
	return filepath.Join(s.DeployDir, "bin", "start_be.sh")
}

func (s BeInstance) SystemdEnvironment() string {
	envs := []string{}
	if s.InstallJava {
		envs = append(envs,
			fmt.Sprintf("Environment=JAVA_HOME=%s/jdk", s.DeployDir))
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
