package config

import (
	"errors"
	"log/slog"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	WindowWidth                    int32   `yaml:"WindowWidth"`
	WindowHeight                   int32   `yaml:"WindowHeight"`
	NumBoids                       int     `yaml:"NumBoids"`
	BoidVelocity                   float32 `yaml:"BoidVelocity"`
	BoidVision                     float32 `yaml:"BoidVision"`
	BoidSeparationOptimalProximity float32 `yaml:"BoidSeparationOptimalProximity"`
	BoidSeparationCoefficient      float32 `yaml:"BoidSeparationCoefficient"`
	BoidAlignmentCoefficient       float32 `yaml:"BoidAlignmentCoefficient"`
	BoidCohesionCoefficient        float32 `yaml:"BoidCohesionCoefficient"`
	NumWorkers                     int     `yaml:"NumWorkers"`
	RandomSeed                     uint64  `yaml:"RandomSeed"`
}

func ParseConfigFile(configFilePath string) (Config, error) {
	slog.Info("loading config file", "configFilePath", configFilePath)

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

	// --------------------------------------------------------------------------------
	// Check logical values of config data

	if config.WindowWidth <= 0 {
		slog.Error("window width must be positive", "WindowWidth", config.WindowWidth)
		return config, errors.New("window width must be positive")
	}

	if config.WindowHeight <= 0 {
		slog.Error("window height must be positive", "WindowHeight", config.WindowHeight)
		return config, errors.New("window height must be positive")
	}

	if config.NumWorkers <= 0 {
		slog.Error("num workers must be positive", "NumWorkers", config.NumWorkers)
		return config, errors.New("num workers must be positive")
	}

	if config.NumBoids <= 0 {
		slog.Error("num boids must be positive", "NumBoids", config.NumBoids)
		return config, errors.New("num boids must be positive")
	}

	if config.BoidVelocity <= 0 {
		slog.Error("boid velocity must be positive", "BoidVelocity", config.BoidVelocity)
		return config, errors.New("boid velocity must be positive")
	}

	if config.BoidVision <= 0 {
		slog.Error("boid vision must be positive", "BoidVision", config.BoidVision)
		return config, errors.New("boid vision must be positive")
	}

	if config.BoidSeparationOptimalProximity < 0 || config.BoidSeparationOptimalProximity > 1 {
		slog.Error("boid separation optimal proximity must be between 0 and 1", "BoidSeparationOptimalProximity", config.BoidSeparationOptimalProximity)
		return config, errors.New("boid separation optimal proximity must be between 0 and 1")
	}

	// --------------------------------------------------------------------------------

	slog.Debug("loaded config data", "config", config)

	return config, nil
}
