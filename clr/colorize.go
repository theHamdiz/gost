package clr

import (
	"strings"
)

func Colorize(text, color string) string {
	var colorCode string
	switch color {
	case "teal":
		colorCode = "\033[96m" // Bright Cyan
	case "pink":
		colorCode = "\033[95m" // Bright Magenta
	case "black":
		colorCode = "\033[90m" // Bright Black (Gray)
	case "red":
		colorCode = "\033[91m" // Bright Red
	case "green":
		colorCode = "\033[92m" // Bright Green
	default:
		colorCode = "\033[0m" // Reset to default
	}
	if strings.Contains(text, "👉") {
		msg := strings.SplitN(text, "👉", -1)
		return colorCode + msg[0] + "\033[96m" + "👉" + "\033[95m" + msg[1]
	}
	return colorCode + text + "\033[0m"
}
