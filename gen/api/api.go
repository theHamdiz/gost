package gen

import (
	"github.com/theHamdiz/gost/config"
)

type ApiGenerator struct {
	Files []string
}

func (g *ApiGenerator) Generate(data config.ProjectData) error {
	return nil
}

func NewApiGenerator() *ApiGenerator {
	return &ApiGenerator{}
}
