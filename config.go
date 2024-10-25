package main

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	NumBoids   int   `yaml:"NumBoids"`
	NumWorkers int   `yaml:"NumWorkers"`
	RandomSeed int64 `yaml:"RandomSeed"`
}

func ParseConfigFile(configFilePath string) (Config, error) {
	slog.Debug("start loading of config file", "configFilePath", configFilePath)

	slog.Debug("open and read config file")
	config := Config{}
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		slog.Error("error reading config file", "configFilePath", configFilePath, "error", err)
		return config, err
	}

	slog.Debug("unmarshal config file buffer to Config struct")
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		slog.Error("error unmarshalling config file", "configFilePath", configFilePath, "error", err)
		return config, err
	}

	return config, nil
}
