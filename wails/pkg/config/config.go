package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config is a struct to hold your configuration values.
// Define your config fields here.
type Credentials struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
type Repository struct {
	Path string `mapstructure:"path"`
}
type Integration struct {
	Path        string `mapstructure:"path"`
	MaxFileSize int64  `mapstructure:"max_file_size"`
}

type ApplicationConfig struct {
	Credentials Credentials `mapstructure:"credentials"` // Nested struct
	Repository  Repository  `mapstructure:"repository"`
	Integration Integration `mapstructure:"integration"`
}

// LoadConfig loads config.toml from the current directory.
func LoadConfig() (*ApplicationConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/home/rnm/Dev/git-repo-test")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg ApplicationConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
