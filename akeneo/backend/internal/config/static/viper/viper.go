package viper

import (
	kit_config "akeneo-event-config/internal/config/static"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type viperConfig struct{}

// NewViperConfig fetch configurations.
func NewViperConfig() kit_config.ConfigurationLoader {
	return &viperConfig{}
}

// LoadConfiguration load the setup for the configuration object.
func (vp *viperConfig) LoadConfiguration(context string) error {
	if viper.IsSet(context) {
		return nil
	}
	viperContext := *viper.New()

	cwd, _ := os.Getwd() //nolint:errcheck // fallback to current directory

	// Support both ENVIRONMENT and ENV variables
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env == "" {
		env = strings.ToLower(os.Getenv("ENV"))
	}
	if env == "" {
		env = "dev"
	}

	servicePath := strings.ToLower(os.Getenv("CONFIG_PATH"))

	// Set the file name of the configurations file
	configPath := fmt.Sprintf("%s/akeneo/backend/config/settings.%s.json", cwd, env)

	// Try multiple paths for config file
	if env == "pipeline" {
		_, compilationPath, _, _ := runtime.Caller(0) //nolint:errcheck // compilation path is always available
		projectPath := filepath.Join(filepath.Dir(compilationPath), "../../../..")
		configPath = fmt.Sprintf("%s/config/%s/settings.%s.json", projectPath, servicePath, env)
	}

	// If running from backend directory, adjust path
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = fmt.Sprintf("%s/../config/settings.%s.json", cwd, env)
	}

	viperContext.SetConfigName(filepath.Base(configPath))
	viperContext.SetConfigType("json")
	viperContext.AddConfigPath(filepath.Dir(configPath))

	// Enable VIPER to read Environment Variables
	viperContext.AutomaticEnv()
	viperContext.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viperContext.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("config file not found: %s (searched: %s)", filepath.Base(configPath), filepath.Dir(configPath))
		}
		return fmt.Errorf("fatal error config file: %w", err)
	}

	fmt.Printf("Loaded config from: %s\n", viperContext.ConfigFileUsed())
	viper.Set(context, viperContext)

	return nil
}
