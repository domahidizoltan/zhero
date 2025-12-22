package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Env    EnvConfig    `mapstructure:"env"`
		Log    LogConfig    `mapstructure:"log"`
		DB     DBConfig     `mapstructure:"database"`
		Admin  AdminConfig  `mapstructure:"admin"`
		Public PublicConfig `mapstructure:"public"`
	}

	EnvConfig struct {
		Platform     string `mapstructure:"platform"`
		AbsolutePath string
	}

	LogConfig struct {
		Format string `mapstructure:"format"`
		Level  string `mapstructure:"level"`
		Color  bool   `mapstructure:"color"`
	}

	DBConfig struct {
		SQLite SQLiteConfig `mapstructure:"sqlite"`
	}

	SQLiteConfig struct {
		File string `mapstructure:"file"`
	}

	AdminConfig struct {
		Server ServerConfig `mapstructure:"server"`
		RDF    RdfConfig    `mapstructure:"rdf"`
	}

	PublicConfig struct {
		Server ServerConfig `mapstructure:"server"`
	}

	ServerConfig struct {
		Port int `mapstructure:"port"`
	}

	RdfConfig struct {
		File   string `mapstructure:"file"`
		Source string `mapstructure:"source"`
	}
)

func LoadConfig(absolutePath string) (*Config, error) {
	viper.SetConfigFile(absolutePath + "config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.Env.AbsolutePath = absolutePath
	return &cfg, nil
}
