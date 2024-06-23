package dirs

import (
	"log"
	"os"
	"path/filepath"
)

import "github.com/theHamdiz/gost/gen/config"

type DbGenerator struct {
	Dirs []string
}

func (g *DbGenerator) Generate(data config.ProjectData) error {
	projectDir := data.ProjectDir
	for _, dir := range g.Dirs {
		path := filepath.Join(projectDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("Error creating directory %s: %v", path, err)
		}
	}
	return nil
}

func NewApiGenerator() *DbGenerator {
	return &DbGenerator{
		Dirs: []string{
			"cmd/app",
			"cmd/scripts",
			"app/cfg",
			"app/handlers",
			"app/middleware",
			"app/types/models",
			"app/types/core",
			"app/router",
			"app/services",
			"app/views",
			"app/assets/static/css",
			"app/assets/static/js",
			"app/assets/static/img",
			"app/views/layouts",
			"app/views/components/header",
			"app/views/components/footer",
			"app/views/components/navigation",
			"app/views/pages",
			"app/views/errors",
			"app/db/migrations",
			"app/events",
			"app/api/v1",
			"plugins/auth",
			"public/assets",
		},
	}
}
