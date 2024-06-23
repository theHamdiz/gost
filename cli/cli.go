package cli

import (
	"github.com/spf13/cobra"
	"github.com/theHamdiz/gost/cfg"
)

type GostCli struct {
	Config   *cfg.GostConfig
	Commands *cobra.Command
}
