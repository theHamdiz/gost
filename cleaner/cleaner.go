package cleaner

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func SortImports(filePath string) error {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Convert content to string
	contentStr := string(content)

	// Find the import section
	startIdx := strings.Index(contentStr, "import (")
	if startIdx == -1 {
		// No Import Section Found, Just return nil.
		return nil
	}
	endIdx := strings.Index(contentStr[startIdx:], ")")
	if endIdx == -1 {
		return fmt.Errorf(">>Gost>>no closing parenthesis for import section found")
	}
	endIdx += startIdx

	// Extract the import section
	importSection := contentStr[startIdx+len("import (") : endIdx]
	importLines := strings.Split(importSection, "\n")

	// Trim spaces and filter out empty lines
	var imports []string
	for _, line := range importLines {
		line = strings.TrimSpace(line)
		if line != "" {
			imports = append(imports, line)
		}
	}

	// Sort the imports
	sort.Strings(imports)

	// Reconstruct the import section
	sortedImports := "import (\n    " + strings.Join(imports, "\n    ") + "\n)"

	// Replace the original import section with the sorted one
	newContent := contentStr[:startIdx] + sortedImports + contentStr[endIdx+1:]

	// Write the sorted content back to the file
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
