package api

import "github.com/kelseyhightower/envconfig"

// Config represents the configuration of the application
type Config struct {
	AppPort        string `default:"8080" envconfig:"APP_PORT"`
	ServiceName    string `default:"bookmark_service" envconfig:"SERVICE_NAME"`
	InstanceID     string `default:"" envconfig:"INSTANCE_ID"`
	BasePath       string `default:"/" envconfig:"BASE_PATH"`
	LogLevel       string `default:"info" envconfig:"LOG_LEVEL"`
	PrivateKeyPath string `default:"./private.key" envconfig:"PRIVATE_KEY_PATH"`
	PublicKeyPath  string `default:"./public.key" envconfig:"PUBLIC_KEY_PATH"`
}

// NewConfig creates a new config struct based on the environment variables
func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
