package config

import (
	"encoding/json"
	"os"
	"sync"
)

const (
	configFile = "data/config.json"
	CoverDir   = "data/cover"
)

type Config struct {
	Port          int    `json:"port"`
	Host          string `json:"host"`
	AllowRegister bool   `json:"allowRegister"`
}

func NewConfig() *Config {
	return &Config{
		Port:          6060,
		Host:          "0.0.0.0",
		AllowRegister: false,
	}
}

func LoadConfig() *Config {
	config := NewConfig()

	data, err := os.ReadFile(configFile)
	if err != nil {
		return config
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return NewConfig()
	}

	return config
}

var (
	configInstance *Config
	once           sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		configInstance = LoadConfig()
	})
	return configInstance
}

func SetConfig(config *Config) {
	configInstance = config
}
