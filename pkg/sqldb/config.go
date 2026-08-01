package sqldb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     string `default:"localhost" envconfig:"DB_HOST"`
	Port     string `default:"5432" envconfig:"DB_PORT"`
	DBName   string `default:"bookmark_db" envconfig:"DB_NAME"`
	User     string `default:"postgres" envconfig:"DB_USER"`
	Password string `default:"postgres" envconfig:"DB_PASSWORD"`
	SSLMode  string `default:"disable" envconfig:"DB_SSL_MODE"`
	Timezone string `default:"Asia/Ho_Chi_Minh" envconfig:"DB_TIMEZONE"`
}

func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode, c.Timezone,
	)
}