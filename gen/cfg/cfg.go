package gen

import "github.com/theHamdiz/gost/gen/config"

type CfgGenerator struct {
	Files []string
}

func (g *CfgGenerator) Generate(data config.ProjectData) error {
	return nil
}

func NewApiGenerator() *CfgGenerator {
	return &CfgGenerator{}
}
