package gen

import "github.com/theHamdiz/gost/gen/config"

type MiddlewareGenerator struct {
	Files []string
}

func (g *MiddlewareGenerator) Generate(data config.ProjectData) error {
	return nil
}

func NewApiGenerator() *MiddlewareGenerator {
	return &MiddlewareGenerator{}
}
