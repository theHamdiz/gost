package dirs

import (
	"log"
	"os"
	"path/filepath"
)

type Generator struct {
	Dirs []string
}

func (g *Generator) Generate(projectDir string) error {
	for _, dir := range g.Dirs {
		path := filepath.Join(projectDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("Error creating directory %s: %v", path, err)
		}
	}
	return nil
}

func NewDirsGenerator() *Generator {
	return &Generator{
		Dirs: []string{
			"cmd/app",
			"log/",
			"storage/",
			"cmd/scripts",
			"app/cfg",
			"app/handlers",
			"app/middleware",
			"app/types/models",
			"app/types/core",
			"app/router",
			"app/services",
			"app/ui",
			"app/assets/static/css",
			"app/assets/static/js",
			"app/assets/static/img",
			"app/ui/layouts",
			"app/ui/components/header",
			"app/ui/components/footer",
			"app/ui/components/navigation",
			"app/ui/pages",
			"app/ui/errors",
			"app/ui/shared",
			"app/db/migrations",
			"app/events",
			"app/api/v1",
			"plugins/auth",
			"public/assets",
		},
	}
}
