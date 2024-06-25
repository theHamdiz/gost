package codegen

import (
	"fmt"

	"github.com/theHamdiz/gost/codegen/api"
	"github.com/theHamdiz/gost/codegen/cfg"
	"github.com/theHamdiz/gost/codegen/db"
	"github.com/theHamdiz/gost/codegen/events"
	"github.com/theHamdiz/gost/codegen/files"
	"github.com/theHamdiz/gost/codegen/handlers"
	"github.com/theHamdiz/gost/codegen/middleware"
	genPlugins "github.com/theHamdiz/gost/codegen/plugins"
	"github.com/theHamdiz/gost/codegen/router"
	"github.com/theHamdiz/gost/codegen/services"
	"github.com/theHamdiz/gost/codegen/types"
	"github.com/theHamdiz/gost/codegen/ui"
	"github.com/theHamdiz/gost/codegen/views"
	"github.com/theHamdiz/gost/config"
	"github.com/theHamdiz/gost/plugins"
	pmCfg "github.com/theHamdiz/gost/plugins/config"
)

func ExecuteGeneration(data config.ProjectData) error {
	generators := []plugins.Plugin{
		api.NewGenApiPlugin(data),
		cfg.NewGenConfPlugin(data),
		db.NewGenDbPlugin(data),
		events.NewGenEventsPlugin(data),
		files.NewGenerator(),
		handlers.NewGenerator(),
		middleware.NewGenerator(),
		genPlugins.NewGenerator(),
		router.NewGenerator(),
		services.NewGenerator(),
		views.NewGenerator(),
		files.NewGenerator(),
		types.NewGenerator(),
		ui.NewGenerator(),
	}

	//for _, generator := range generators {
	//	if err := generator.Generate(data); err != nil {
	//		return err
	//	}
	//}

	pmConfig := &pmCfg.PluginManagerConfig{
		PluginsDir: "plugins",
	}
	pm := plugins.NewPluginManager(pmConfig)
	err := pm.RegisterPlugins(generators)
	if err != nil {
		return err
	}
	err = pm.InitPlugins()
	if err != nil {
		fmt.Println(err)
	}

	err = pm.ExecutePlugins()
	if err != nil {
		fmt.Println(err)
	}

	err = pm.ShutdownPlugins()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
