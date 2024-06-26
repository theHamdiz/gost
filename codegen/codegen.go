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
	"github.com/theHamdiz/gost/codegen/temmlatplates"
	"github.com/theHamdiz/gost/codegen/typespes"
	"github.com/theHamdiz/gost/codegen/ui
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
		files.NewGenFilesPlugin(data),
		handlers.NewGenHandlersPlugin(data),
		middleware.NewGenMiddlewarePlugin(data),
		genPlugins.NewGenPluginsPlugin(data),
		router.NewGenRouterPlugin(data),
		services.NewGenServicesPlugin(data),
		templates.NewGenViewsPlugin(data),
		types.NewGenTypesPlugin(data),
		ui.NewGenUiPlugin(data),
	}

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

	defer func() {
		if err := pm.ShutdownPlugins(); err != nil {
			fmt.Println("Error during plugin shutdown:", err)
		}
	}()

	err = pm.ExecutePlugins()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
