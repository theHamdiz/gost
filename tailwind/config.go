package tailwind

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AppendToTailwindConfig appends new entries to a specified section in the tailwind.config.js file.
func AppendToTailwindConfig(filePath, section, newEntry string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf(">>Gost>> failed to open file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	var outputLines []string
	scanner := bufio.NewScanner(file)
	inSection := false
	sectionFound := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Check if we are in the section
		if strings.HasPrefix(trimmedLine, section+":") {
			inSection = true
			sectionFound = true
		}

		// Check if we are leaving the section
		if inSection && strings.HasPrefix(trimmedLine, "}") {
			inSection = false
		}

		// Add the new entry before leaving the section
		if inSection && trimmedLine == "]," {
			outputLines = append(outputLines, "\t\t"+newEntry+",")
		}

		outputLines = append(outputLines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf(">>Gost>> error reading file: %w", err)
	}

	// If section was not found, create it
	if !sectionFound {
		outputLines = append(outputLines, fmt.Sprintf("\t%s: [", section), "\t\t"+newEntry+",", "\t],")
	}

	// Write the modified content back to the file
	file, err = os.Create(filePath)
	if err != nil {
		return fmt.Errorf(">>Gost>> failed to create file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	for _, line := range outputLines {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func GetPluginConfigForComponentFramework(componentFramework string) string {
	switch strings.ToLower(componentFramework) {
	case "daisyui":
		return "require('daisyui'),"
	case "prelineui":
		return "require('preline/plugin'),"
	case "flowbite":
		return "require('flowbite/plugin'),"
	case "tw-elements":
		return "require(\"tw-elements/plugin.cjs\")"
	default:
		return ""
	}
}

func GetContentConfigForComponentFramework(componentFramework string) string {
	switch strings.ToLower(componentFramework) {
	case "daisyui":
		return ""
	case "prelineui":
		return "preline"
	case "flowbite":
		return "flowbite"
	case "tw-elements":
		return "./node_modules/tw-elements/js/**/*.js"
	default:
		return ""
	}
}
