package tailwind

import "strings"

func GetInstallationCommandFor(compnentFramework string) string {
	switch strings.ToLower(compnentFramework) {
	case "daisyui":
		return "npm i -D daisyui@latest"
	case "prelineui":
		return "npm i preline"
	case "tw-elements":
		return "npm install tw-elements"
	case "flowbite":
		return "npm install flowbite"
	default:
		return ""
	}
}
