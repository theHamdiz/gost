package dirs

import (
	"log"
	"os"
	"path/filepath"
)

var dirs = []string{
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
	"plugins/auth",
	"public/assets",
}

func Generate(projectDir, appName string) error {
	for _, dir := range dirs {
		path := filepath.Join(projectDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("Error creating directory %s: %v", path, err)
		}
	}
	return nil
}
