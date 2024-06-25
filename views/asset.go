package views

import (
	"fmt"
	"os"
	"path/filepath"
)

// Asset retrieves the content of a file from the current working directory under app/assets/{any_folder}/{any_asset}
func Asset(fileName string) ([]byte, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf(">>Gost>> error getting current working directory: %v", err)
	}

	// Construct the full path to the asset
	assetPath := filepath.Join(cwd, "app", "assets", fileName)

	// Read the file content
	content, err := os.ReadFile(assetPath)
	if err != nil {
		return nil, fmt.Errorf(">>Gost>> error reading file %s: %v", assetPath, err)
	}

	return content, nil
}
