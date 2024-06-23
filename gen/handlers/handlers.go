package gen

import "github.com/theHamdiz/gost/gen/config"

type HandlersGenerator struct {
	Files []string
}

func (g *HandlersGenerator) Generate(data config.ProjectData) error {
	return nil
}

func NewApiGenerator() *HandlersGenerator {
	return &HandlersGenerator{}
}
