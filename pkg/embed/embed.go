package embed

import (
	"embed"
	"text/template"
)

//go:embed systemd
var embededFiles embed.FS

const (
	SystemdServicePath = "systemd/component.service"
)

func ReadTemplate(path string) (*template.Template, error) {
	return template.ParseFS(embededFiles, path)
}
