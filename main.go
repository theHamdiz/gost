package main

import (
	"fmt"

	"os"

	"github.com/theHamdiz/gost/gen"
)

var config gen.GostConfig

func main() {
	if len(os.Args) == 1 {
		config := &gen.GostConfig{}
		buildConfig(config)
	} else {
		rootCmd := addCommands(config)
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
