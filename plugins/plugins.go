package plugins

import (
	"errors"
	"fmt"

	"github.com/theHamdiz/gost/plugins/config"
)

type Plugin interface {
	Init() error
	Execute() error
	Shutdown() error
	Name() string
	Version() string
	Dependencies() []string
	AuthorName() string
	AuthorEmail() string
	Website() string
	GitHub() string
}

type GenericPluginManager interface {
	InitPlugins() error
	ExecutePlugins() error
	ShutdownPlugins() error
	RegisterPlugins(plugins []Plugin) error
	RegisterPlugin(plugin Plugin) error
}

type PluginMetadata struct {
	Name         string   `yaml:"name" json:"name"`
	Version      string   `yaml:"version"  json:"version"`
	Dependencies []string `yaml:"dependencies"  json:"dependencies"`
	AuthorName   string   `yaml:"author_name"  json:"author_name"`
	AuthorEmail  string   `yaml:"author_email"  json:"author_email"`
	Website      string   `yaml:"website"  json:"website"`
	GitHub       string   `yaml:"github"  json:"github"`
}

type PluginManager struct {
	plugins     map[string]Plugin
	pluginOrder []string
	config      *config.PluginManagerConfig
}

func (pm *PluginManager) RegisterPlugin(plugin Plugin) error {
	if plugin.Name() == "" {
		return errors.New("plugin name cannot be empty")
	}
	pm.plugins[plugin.Name()] = plugin
	pm.pluginOrder = append(pm.pluginOrder, plugin.Name())
	return nil
}

func (pm *PluginManager) RegisterPlugins(plugins []Plugin) error {
	for _, plugin := range plugins {
		err := pm.RegisterPlugin(plugin)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewPluginManager(config *config.PluginManagerConfig) *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
		config:  config,
	}
}

func (pm *PluginManager) resolveDependencies() ([]Plugin, error) {
	resolved := make(map[string]bool)
	var order []Plugin

	var resolve func(name string) error
	resolve = func(name string) error {
		if resolved[name] {
			return nil
		}

		plugin, exists := pm.plugins[name]
		if !exists {
			return fmt.Errorf("plugin %s not found", name)
		}

		for _, dep := range plugin.Dependencies() {
			if err := resolve(dep); err != nil {
				return err
			}
		}

		resolved[name] = true
		order = append(order, plugin)
		return nil
	}

	for name := range pm.plugins {
		if err := resolve(name); err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (pm *PluginManager) InitPlugins() error {
	for _, name := range pm.pluginOrder {
		plugin := pm.plugins[name]
		if err := plugin.Init(); err != nil {
			return err
		}
	}
	return nil
}

func (pm *PluginManager) ExecutePlugins() error {
	for _, name := range pm.pluginOrder {
		plugin := pm.plugins[name]
		if err := plugin.Execute(); err != nil {
			return err
		}
	}
	return nil
}

func (pm *PluginManager) ShutdownPlugins() error {
	// Reverse the order for shutdown
	for i := len(pm.pluginOrder) - 1; i >= 0; i-- {
		name := pm.pluginOrder[i]
		plugin := pm.plugins[name]
		if err := plugin.Shutdown(); err != nil {
			return err
		}
	}
	return nil
}
