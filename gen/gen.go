package gen

import (
	"github.com/theHamdiz/gost/gen/config"
)

// Generator interface with the Generate method
type Generator interface {
	Generate(data config.ProjectData) error
}

generators := []Generator{
	
}