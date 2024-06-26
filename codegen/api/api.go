package api

import (
	"github.com/theHamdiz/gost/config"
)

type GenApiPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenApiPlugin) Init() error {
	// Any initialization logic for the plugin
	return nil
}

func (g *GenApiPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenApiPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenApiPlugin) Name() string {
	return "GenApiPlugin"
}

func (g *GenApiPlugin) Version() string {
	return "1.0.0"
}

func (g *GenApiPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenApiPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenApiPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenApiPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenApiPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/api"
}

func (g *GenApiPlugin) Generate(data config.ProjectData) error {
	// Logic for generating API code
	return nil
}

func NewGenApiPlugin(data config.ProjectData) *GenApiPlugin {
	return &GenApiPlugin{
		Files: map[string]func() string{
			"app/api/v1/api.go": func() string { return `package api` },
		},
		Data: data,
	}
}
