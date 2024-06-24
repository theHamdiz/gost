package plugins

import (
	"fmt"

	"github.com/theHamdiz/gost/plugins/config"
	"github.com/theHamdiz/gost/services"
)

type PluginManager struct {
	plugins map[string]Plugin
	config  Config
	sc      *ServiceContainer
}

func NewPluginManager(config config.PluginConfig, sc *services.ServiceContainer) *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
		config:  config,
		sc:      sc,
	}
}

func (pm *PluginManager) RegisterPlugin(name string, plugin Plugin) {
	pm.plugins[name] = plugin
}

func (pm *PluginManager) InitPlugins() error {
	for name, plugin := range pm.plugins {
		cfg := pm.config.GetPluginConfig(name)
		if err := plugin.Init(cfg, pm.sc); err != nil {
			return fmt.Errorf("failed to initialize plugin %s: %v", name, err)
		}
	}
	return nil
}

func (pm *PluginManager) ExecutePlugins() error {
	for name, plugin := range pm.plugins {
		if err := plugin.Execute(); err != nil {
			return fmt.Errorf("failed to execute plugin %s: %v", name, err)
		}
	}
	return nil
}

func (pm *PluginManager) ShutdownPlugins() error {
	for name, plugin := range pm.plugins {
		if err := plugin.Shutdown(); err != nil {
			return fmt.Errorf("failed to shutdown plugin %s: %v", name, err)
		}
	}
	return nil
}
