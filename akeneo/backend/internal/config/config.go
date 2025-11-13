package config

import (
	"fmt"

	viperLoader "akeneo-event-config/internal/config/static/viper"

	"github.com/spf13/viper"
)

const configContext = "app_config"

type Config struct {
	Akeneo AkeneoConfig
	Server ServerConfig
	CORS   CORSConfig
}

type AkeneoConfig struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type CORSConfig struct {
	AllowedOrigins string
}

func Load() (*Config, error) {
	// Load configuration using viper
	loader := viperLoader.NewViperConfig()
	if err := loader.LoadConfiguration(configContext); err != nil {
		return nil, err
	}

	// Get the viper instance for this context
	v := viper.Get(configContext).(viper.Viper)

	config := &Config{
		Akeneo: AkeneoConfig{
			BaseURL:      v.GetString("akeneo.base_url"),
			ClientID:     v.GetString("akeneo.client_id"),
			ClientSecret: v.GetString("akeneo.client_secret"),
			Username:     v.GetString("akeneo.username"),
			Password:     v.GetString("akeneo.password"),
		},
		Server: ServerConfig{
			Port:    v.GetString("server.port"),
			GinMode: v.GetString("server.gin_mode"),
		},
		CORS: CORSConfig{
			AllowedOrigins: v.GetString("cors.allowed_origins"),
		},
	}

	// Set defaults if not provided
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.GinMode == "" {
		config.Server.GinMode = "debug"
	}
	if config.CORS.AllowedOrigins == "" {
		config.CORS.AllowedOrigins = "http://localhost:3000"
	}

	if err := validate(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validate(cfg *Config) error {
	if cfg.Akeneo.BaseURL == "" {
		return fmt.Errorf("akeneo.base_url is required")
	}
	if cfg.Akeneo.ClientID == "" {
		return fmt.Errorf("akeneo.client_id is required")
	}
	if cfg.Akeneo.ClientSecret == "" {
		return fmt.Errorf("akeneo.client_secret is required")
	}
	if cfg.Akeneo.Username == "" {
		return fmt.Errorf("akeneo.username is required")
	}
	if cfg.Akeneo.Password == "" {
		return fmt.Errorf("akeneo.password is required")
	}
	return nil
}
