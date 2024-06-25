package loader

import (
	"io/ioutil"

	"github.com/theHamdiz/gost/plugins"
	"gopkg.in/yaml.v3"
)

func LoadPluginMetadata(metaFile string) (*plugins.PluginMetadata, error) {
	data, err := ioutil.ReadFile(metaFile)
	if err != nil {
		return nil, err
	}

	var meta plugins.PluginMetadata
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	return &meta, nil
}
