package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Admin AdminConfig `mapstructure:"admin"`
	}

	AdminConfig struct {
		Server AdminServerConfig `mapstructure:"server"`
		RDF    RdfConfig         `mapstructure:"rdf"`
	}

	AdminServerConfig struct {
		Port int `mapstructure:"port"`
	}

	RdfConfig struct {
		File   string `mapstructure:"file"`
		Source string `mapstructure:"source"`
	}
)

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
