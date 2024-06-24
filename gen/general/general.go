package general

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/theHamdiz/gost/cleaner"
	"github.com/theHamdiz/gost/config"
	"github.com/theHamdiz/gost/parser"
)

// Generator interface with the Generate method
type Generator interface {
	Generate(data config.ProjectData) error
}

func GenerateFiles(data config.ProjectData, files map[string]func() string) error {
	for path, tmplFunc := range files {
		content, err := parser.ParseTemplateStringAsText(path, tmplFunc(), data)
		if err != nil {
			return fmt.Errorf(">> failed to parse template %s: %w", path, err)
		}

		appNameDowncased := strings.ToLower(data.AppName)
		filePath := filepath.Join(appNameDowncased, path)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("✗ failed to create directory: %w", err)
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		} else {
			if strings.HasSuffix(filePath, ".go") {
				if err = cleaner.SortImports(filePath); err != nil {
					return fmt.Errorf("✗ The file was saved but failed to sort its imports %w", err)
				}
			}
		}
	}
	return nil
}
