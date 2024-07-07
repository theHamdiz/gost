package scripts

import (
	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type GenScriptsPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenScriptsPlugin) Init() error {
	// Initialize Files
	g.Files = map[string]func() string{
		"cmd/scripts/setup.sh": func() string {
			return ``
		},
		"cmd/scripts/migrate.sh": func() string {
			return ``
		},
		"cmd/scripts/deploy.sh": func() string {
			return ``
		},
		"cmd/scripts/backup.sh": func() string {
			return ``
		},
		"cmd/scripts/cleanup.sh": func() string {
			return ``
		},
		"cmd/scripts/gendocs.sh": func() string {
			return ``
		},
	}

	return nil
}

func (g *GenScriptsPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenScriptsPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenScriptsPlugin) Name() string {
	return "GenScriptsPlugin"
}

func (g *GenScriptsPlugin) Version() string {
	return "1.0.0"
}

func (g *GenScriptsPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenScriptsPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenScriptsPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenScriptsPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenScriptsPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/services"
}

func (g *GenScriptsPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenScriptsPlugin(data config.ProjectData) *GenScriptsPlugin {
	return &GenScriptsPlugin{
		Data: data,
	}
}
