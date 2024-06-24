package gen

import (
	"github.com/theHamdiz/gost/config"
	"github.com/theHamdiz/gost/gen/cfg"
	"github.com/theHamdiz/gost/gen/db"
	"github.com/theHamdiz/gost/gen/events"
	"github.com/theHamdiz/gost/gen/files"
	"github.com/theHamdiz/gost/gen/general"
	"github.com/theHamdiz/gost/gen/handlers"
	"github.com/theHamdiz/gost/gen/middleware"
	"github.com/theHamdiz/gost/gen/plugins"
	"github.com/theHamdiz/gost/gen/router"
	"github.com/theHamdiz/gost/gen/services"
	"github.com/theHamdiz/gost/gen/types"
	"github.com/theHamdiz/gost/gen/ui"
	"github.com/theHamdiz/gost/gen/views"
)

func ExecuteGeneration(data config.ProjectData) error {
	generators := []general.Generator{
		cfg.NewGenerator(),
		db.NewGenerator(),
		events.NewGenerator(),
		files.NewGenerator(),
		handlers.NewGenerator(),
		middleware.NewGenerator(),
		plugins.NewGenerator(),
		router.NewGenerator(),
		services.NewGenerator(),
		views.NewGenerator(),
		files.NewGenerator(),
		types.NewGenerator(),
		ui.NewGenerator(),
	}

	for _, generator := range generators {
		if err := generator.Generate(data); err != nil {
			return err
		}
	}
	return nil
}
