package config

import (
	"os"

	"github.com/ricoberger/httpmonitor/pkg/target"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Targets []target.Config `yaml:"targets"`
}

func New(file string) (*Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config

	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
