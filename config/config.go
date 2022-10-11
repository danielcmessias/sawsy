package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Theme ThemeConfig `yaml:"theme"`
}

type ThemeConfig struct {
	ShowIcons bool `yaml:"showIcons"`
}

func ReadConfig() (Config, error) {
	config := getDefaultConfig()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}
	path := filepath.Join(homeDir, ".sawsy.yml")

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
