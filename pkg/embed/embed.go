package embed

import (
	"embed"
	"text/template"
)

//go:embed scripts systemd
var embededFiles embed.FS

const (
	FeScriptPath       = "scripts/doris_fe.sh"
	BeScriptPath       = "scripts/doris_be.sh"
	SystemdServicePath = "systemd/component.service"
)

func ReadTemplate(path string) (*template.Template, error) {
	return template.ParseFS(embededFiles, path)
}
