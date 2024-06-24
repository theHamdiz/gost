package config

import (
	"encoding/json"
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
	defer file.Close()

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
