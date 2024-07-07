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
			"cmd/server",
			"cmd/worker",
			"cmd/scripts",
			"log/",
			"storage/",
			"app/cfg",
			"app/handlers",
			"app/middleware",
			"app/types/models",
			"app/types/core",
			"app/router",
			"app/services",
			"app/web/shared/components/header",
			"app/web/shared/components/footer",
			"app/web/shared/components/navigation",
			"app/web/shared/layouts",
			"app/web/shared/pages/signin",
			"app/web/shared/pages/signup",
			"app/web/shared/public",
			"app/web/frontend",
			"app/web/backend",
			"app/assets/static/css",
			"app/assets/static/js",
			"app/assets/static/img",
			"app/web/pages",
			"app/web/errors",
			"app/web/shared",
			"app/db/migrations",
			"app/events",
			"app/api/http/v1",
			"app/api/grpc/v1/proto",
			"app/api/grpc/v1/server",
			"app/api/grpc/v1/client",
			"plugins/auth",
			"public/assets",
		},
	}
}
