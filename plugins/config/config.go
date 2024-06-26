package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type PluginConfig struct {
	Settings map[string]interface{}
}

func LoadConfig(filePath string) (*PluginConfig, error) {
	config := &PluginConfig{
		Settings: make(map[string]interface{}),
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config.Settings); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *PluginConfig) GetPluginConfig(pluginName string) map[string]interface{} {
	if cfg, ok := c.Settings[pluginName].(map[string]interface{}); ok {
		return cfg
	}
	return nil
}

type PluginManagerConfig struct {
	PluginsDir string
}
