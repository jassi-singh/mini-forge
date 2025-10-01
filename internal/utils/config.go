package utils

import (
	"errors"
	"fmt"

	"os"

	"gopkg.in/yaml.v2"
)

// Config holds the application configuration values

type Config struct {
	Port      int `yaml:"PORT"`
	RangeSize int `yaml:"RANGE_SIZE"`
}

// LoadConfig reads configuration values from config.yml
func LoadConfig(path string) (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	fmt.Println(configPath)
	if configPath == "" {
		return nil, errors.New("CONFIG_PATH env must be set")
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var appConfig Config
	err = yaml.Unmarshal([]byte(os.ExpandEnv(string(file))), &appConfig)

	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}
