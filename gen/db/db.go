package gen

import "github.com/theHamdiz/gost/gen/config"

type DbGenerator struct {
	Files []string
}

func (g *DbGenerator) Generate(data config.ProjectData) error {
	return nil
}

func NewApiGenerator() *DbGenerator {
	return &DbGenerator{}
}
