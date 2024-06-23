package main

import (
	"fmt"
	"os"

	"github.com/theHamdiz/gost/cfg"
	"github.com/theHamdiz/gost/cli"
	"github.com/theHamdiz/gost/helpers"
)

func main() {
	helpers.Config = &cfg.GostConfig{}
	c := cli.GostCli{
		Config: helpers.Config,
	}
	if len(os.Args) == 1 {
		helpers.BuildConfig(c.Config)
	} else {
		rootCmd := helpers.AddCommands(*c.Config)
		c.Commands = rootCmd
		if err := c.Commands.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
